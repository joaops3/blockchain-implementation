package blockchain

import (
	"blockchain/db"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)


const DbFile = "blockchain_%s.db"

// bucket key = blockHash, value = block
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	Db db.Storage
}


func NewBlockchain(address string) *Blockchain {
	var tip []byte
	

	Db, err := db.NewBoltStorage(DbFile)

	if err != nil {
		log.Fatalf("Error opening database: %s", err.Error())
	}

	err = Db.Update(func(tx db.WriteBucket) error {
		 
		tx.CreateBucketIfNotExists(blocksBucket)
		

		lastHash, err := tx.Get(blocksBucket, []byte("l"))
		if err == nil && len(lastHash) > 0 {
			tip = lastHash
			return nil
		}

		cbtx := NewCoinbaseTX(address, "Genesis block transaction data")
		genesis := NewGenesisBlock(cbtx)

		err = tx.Put(blocksBucket, genesis.Hash, genesis.Serialize())
		if err != nil {
			return err
		}
		err = tx.Put(blocksBucket, []byte("l"), genesis.Hash)
		if err != nil {
			return err
		}

		tip = genesis.Hash
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to create blockchain: %v", err)
	}
	bc := &Blockchain{tip: tip, Db: Db}

	return bc
}

// add block with empty Transaction to blockchain
func (b *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := b.Db.View(func(tx db.ReadBucket) error {
	
		lastHashDb, err := tx.Get(blocksBucket, []byte("l"))

		if err != nil {
			return err
		}
		lastHash = lastHashDb
		return nil
	})

	if err != nil {
		log.Fatalf("Error adding block: %s", err.Error())
	}

	newBlock := NewBlock([]*Transaction{}, lastHash)

	err = b.Db.Update(func(tx db.WriteBucket) error {
		
		err := tx.Put(blocksBucket, newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = tx.Put(blocksBucket, []byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}
		b.tip = newBlock.Hash
		return nil
	})
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block{
	var lastHash []byte

	for _, tx := range transactions {
		if !bc.VerifyTransactions(tx) {
			log.Panic("Invalid transaction")
		}
	}

	err := bc.Db.View(func(tx db.ReadBucket) error {
	
		lastHashDb, err := tx.Get(blocksBucket, []byte("l"))
		if err != nil {
			return err
		}
		lastHash = lastHashDb
		return nil
	})

	if err != nil {
		log.Fatalf("Error getting last hash: %s", err.Error())
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.Db.Update(func(tx db.WriteBucket) error {
		err := tx.Put(blocksBucket, newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = tx.Put(blocksBucket, []byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}
		bc.tip = newBlock.Hash
		return nil
	})
	return newBlock
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
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
			fmt.Printf("tx.ID: %x\n", tx.ID)
			fmt.Printf("ID: %x\n", ID)
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction not found")
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	// to obtain block from top to bottom
	bci := &BlockchainIterator{CurrentHash: bc.tip, Db: bc.Db}

	return bci
}





func (b *Blockchain) FindUnspentTransactionsByAddress(address string) []Transaction {
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

				if otx.IsLockedWithKey(GetPubKeyFromAddress(address)) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, input := range tx.Vin {
					if input.UsesKey([]byte(address)) {
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


func (b *Blockchain) FindAllUnspentUTXO() map[string][]TXOutput {
	UTXO := make(map[string][]TXOutput)
	spentTXOs := make(map[string][]int)
	bci := b.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs = append(outs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}


func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactionsByAddress(address)
	accumulated := 0

	Work:
		for _, tx := range unspentTXs {
			txID := hex.EncodeToString(tx.ID)

			for outIdx, out := range tx.Vout {
				if out.IsLockedWithKey(GetPubKeyFromAddress(address)) && accumulated < amount {
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