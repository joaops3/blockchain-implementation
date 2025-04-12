package blockchain

import (
	"log"

	"github.com/boltDb/bolt"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Db          *bolt.DB
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