package blockchain

import (
	"blockchain/db"
	"encoding/hex"
	"log"
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
	err := u.Blockchain.Db.Update(func(bucket db.WriteBucket) error {
		_ = bucket.CreateBucketIfNotExists(utxoBucket)

		// limpa bucket
		_ = bucket.ForEach(utxoBucket, func(k, _ []byte) error {
			return bucket.Delete(utxoBucket, k)
		})

		utxos := u.Blockchain.FindAllUnspentUTXO()

		for txID, outs := range utxos {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			if err := bucket.Put(utxoBucket, key, SerializeOutputs(outs)); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error reindexing UTXO set: %s", err)
	}
}

// get the unspent outputs by address and amount
func (u *UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	acc := 0
	// map[txid][]outidx
	unspent := make(map[string][]int)

	_ = u.Blockchain.Db.View(func(bucket db.ReadBucket) error {
		_ = bucket.ForEach(utxoBucket, func(k, v []byte) error {
			txID := hex.EncodeToString(k)
			outs := DeserializeTXOutput(v)

			for idx, out := range outs {
				if out.IsLockedWithKey(pubKeyHash) && acc < amount {
					acc += out.Value
					unspent[txID] = append(unspent[txID], idx)
				}
			}
			return nil
		})
		return nil
	})

	return acc, unspent
}

// check balance by address
func (u *UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput

	_ = u.Blockchain.Db.View(func(bucket db.ReadBucket) error {
		_ = bucket.ForEach(utxoBucket, func(_, v []byte) error {
			outs := DeserializeTXOutput(v)
			for _, out := range outs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
			return nil
		})
		return nil
	})

	return UTXOs
}


func (u *UTXOSet) Update(block *Block) {
	err := u.Blockchain.Db.Update(func(bucket db.WriteBucket) error {
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, vin := range tx.Vin {
					outsRaw, err := bucket.Get(utxoBucket, vin.Txid)
					if err != nil || outsRaw == nil {
						continue
					}
					outs := DeserializeTXOutput(outsRaw)

					var newOuts []TXOutput
					for idx, out := range outs {
						if idx != vin.Vout {
							newOuts = append(newOuts, out)
						}
					}

					if len(newOuts) == 0 {
						_ = bucket.Delete(utxoBucket, vin.Txid)
					} else {
						_ = bucket.Put(utxoBucket, vin.Txid, SerializeOutputs(newOuts))
					}
				}
			}

			// adiciona os novos outputs
			_ = bucket.Put(utxoBucket, tx.ID, SerializeOutputs(tx.Vout))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error updating UTXO set: %s", err)
	}
}

func (u *UTXOSet) CountTransactions() int {
	count := 0
	_ = u.Blockchain.Db.View(func(bucket db.ReadBucket) error {
		_ = bucket.ForEach(utxoBucket, func(_, _ []byte) error {
			count++
			return nil
		})
		return nil
	})
	return count
}