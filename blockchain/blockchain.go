package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"log"

	"github.com/boltDb/bolt"
)

const DbFile = "blockchain_%s.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	Db *bolt.DB
}
type BlockchainIterator struct {
	CurrentHash []byte
	Db          *bolt.DB
}




func NewBlockchain(address string) *Blockchain {
	var tip []byte
	

	Db, err := bolt.Open(DbFile, 0600, nil)

	if err != nil {
		log.Fatalf("Error opening database: %s", err.Error())
	}

	err = Db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(blocksBucket))
		

		if bucket == nil {
		
			cbtx := NewCoinbaseTX(address, "Genesis block transaction data")
			genesisBlock := NewGenesisBlock(cbtx)
			bucket, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Fatalf("Error creating bucket: %s", err.Error())
			}
			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Fatalf("Error put bucket: %s", err.Error())
			}
			err = bucket.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Fatalf("Error put 1 bucket: %s", err.Error())
			}

		}else {
			tip = bucket.Get([]byte("l"))
		}

		return nil
	})

	return &Blockchain{tip: tip, Db: Db}
}

func (b *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := b.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			return errors.New("Bucket not found")
		}
		lastHash = bucket.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Fatalf("Error adding block: %s", err.Error())
	}

	newBlock := NewBlock(data, lastHash, []*Transaction{})

	err = b.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			log.Fatalf("Bucket not found")
		}
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}
		b.tip = newBlock.Hash
		return nil
	})
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block{
	return nil
}

func (bc *Blockchain) VerifyTransactions(transaction *Transaction) bool {

	if transaction.IsCoinbase() {
		return true
	}

	
	prevTXs := make(map[string]Transaction)

	for _, in := range transaction.Vin {
		
		prevTx, err := bc.FindTransaction(in.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return transaction.Verify(prevTXs)

}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	// to obtain block from top to bottom
	bci := &BlockchainIterator{CurrentHash: bc.tip, Db: bc.Db}

	return bci
}


func (bci *BlockchainIterator) Next() *Block {
	block := &Block{}
	err := bci.Db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(bci.CurrentHash)
		block = Deserialize(encodedBlock)
		return nil
	})

	if err != nil {
		log.Fatalf("error getting next: %s", err.Error())
	}

	bci.CurrentHash = block.PrevBlockHash
	return block
}


func (b *Blockchain) FindUnspentTransactions(address string) []Transaction {
	unspentTxs := []Transaction{}
	spentTxo := make(map[string][]int)
	bci := b.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			
		Outputs:
			for outputIndex, otx := range tx.Vout {

				if spentTxo[txID] != nil {
					for _, spentIndex := range spentTxo[txID] {
						if outputIndex == spentIndex {
							continue Outputs
						}

					}

				}

				if otx.CanBeUnlockedWith(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, input := range tx.Vin {
					if input.CanUnlockOutputWith(address) {
						inTxId := hex.EncodeToString(input.Txid)
						spentTxo[inTxId] = append(spentTxo[inTxId], input.Vout)
					}
				}

			}

		}


		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTxs
}


func (b *Blockchain) FindUTXO(address string) []TXOutput{
	txOutputs := []TXOutput{}
	unspentTx := b.FindUnspentTransactions(address)

	for _, tx := range unspentTx {
		for _, output := range tx.Vout {

			if output.CanBeUnlockedWith(address) {
				txOutputs = append(txOutputs, output)
			}

		}

	}

	return txOutputs
}


func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

	Work:
		for _, tx := range unspentTXs {
			txID := hex.EncodeToString(tx.ID)

			for outIdx, out := range tx.Vout {
				if out.CanBeUnlockedWith(address) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

					if accumulated >= amount {
						break Work
					}
				}
			}
		}

	return accumulated, unspentOutputs
}