package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"log"

	"github.com/boltDb/bolt"
)

// bucket key = txid, value = txo
const utxoBucket = "utxo"

type UTXOSet struct {
	Blockchain *Blockchain
}

func NewUTXOSet(bc *Blockchain) *UTXOSet {
	return &UTXOSet{bc}
}

// gets all unspent outputs from blockchain, and finally it saves the outputs to the bucket.
func (u *UTXOSet) Reindex() {
	db := u.Blockchain.Db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
	
		if b != nil {
			err := tx.DeleteBucket(bucketName)
			if err != nil {
				return err
			}
		}
	
		_, err := tx.CreateBucket(bucketName)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error reindexing UTXO set: %s", err.Error())
	}

	// remover o address futuramente
	UTXO := u.Blockchain.FindAllUnspentUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("UTXO bucket not found")
		}
		
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(string(txID))
			if err != nil {
				return err
			}
			err = bucket.Put(key, SerializeOutputs(outs))
			if err != nil {
				return err
			}
		}
		return nil
	})

}

// get the unspent outputs by address and amount
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int)(int, map[string][]int){
	// txid -> [outIdx]
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	

	err := u.Blockchain.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utxoBucket))
		if bucket == nil {
			return errors.New("UTXO bucket not found")
		}
		cursor := bucket.Cursor()
		
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeTXOutput(v)
			for outIdx, out := range outs {
			
				
			
                if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
                    accumulated += out.Value
                    unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
                }
            }

		}
		return nil
	})


		if err != nil {
			log.Fatalf("Error finding spendable outputs: %s", err.Error())
		}

	return accumulated, unspentOutputs
}


// check balance by address
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
    var UTXOs []TXOutput
    

    err := u.Blockchain.Db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            outs := DeserializeTXOutput(v)
			
            for _, out := range outs {
                if out.IsLockedWithKey(pubKeyHash) {
                    UTXOs = append(UTXOs, out)
                }
            }
        }

        return nil
    })
	if err != nil {
		log.Fatalf("Error finding UTXOs: %s", err.Error())
	}
    return UTXOs
}


func (u *UTXOSet) Update(block *Block) {
	err := u.Blockchain.Db.Update(func(txn *bolt.Tx) error {
		bucket := txn.Bucket([]byte(utxoBucket))
		if bucket == nil {
			return errors.New("UTXO bucket not found")
		}
	
		for _, tx := range block.Transactions {

			if tx.IsCoinbase() == false {
				for _, vin := range tx.Vin {
					var updatedOutputs = []TXOutput{}
					outBytes := bucket.Get(tx.ID)

					if outBytes == nil {
						continue
					}
		
					outs := DeserializeTXOutput(outBytes)
					for outIdx, out := range outs {
						if outIdx != vin.Vout {
							updatedOutputs = append(updatedOutputs, out)
						}  
					}

					if len(updatedOutputs) == 0 {
						err  := bucket.Delete(vin.Txid)
						if err != nil {
							return err
						}
					}else {
						var buff bytes.Buffer

						encoder := gob.NewEncoder(&buff)
						err := encoder.Encode(updatedOutputs)
						if err != nil {
							return err
						}

						err = bucket.Put(vin.Txid, buff.Bytes())
						if err != nil {
							return err
						}
					}
				}
			}
		newOutputs := []TXOutput{}
			for _, out := range tx.Vout {
				newOutputs = append(newOutputs, out)
			}

			var buff bytes.Buffer

						encoder := gob.NewEncoder(&buff)
						err := encoder.Encode(newOutputs)
						if err != nil {
							return err
						}

			err = bucket.Put(tx.ID, buff.Bytes())
			if err != nil {
				log.Panic(err)
			}
		}
		
		return nil
	})

	if err != nil {
		log.Fatalf("Error updating UTXO set: %s", err.Error())
	}

}

func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.Db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}