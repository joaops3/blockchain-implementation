package blockchain

import (
	"blockchain/db"
	"log"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Db          db.Storage
}

func (bci *BlockchainIterator) Next() *Block {
	block := &Block{}
	err := bci.Db.View(func(tx db.ReadBucket) error {

	
		encodedBlock, err := tx.Get(blocksBucket, bci.CurrentHash)
		if err != nil {
			log.Fatalf("error getting block: %s", err.Error())
		}
		block = Deserialize(encodedBlock)
		return nil
	})

	if err != nil {
		log.Fatalf("error getting next: %s", err.Error())
	}

	bci.CurrentHash = block.PrevBlockHash
	return block
}