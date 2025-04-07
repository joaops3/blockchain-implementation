package blockchain

import (
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

	newBlock := NewBlock(data, lastHash)

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

func NewBlockchain() *Blockchain {
	var tip []byte
	

	Db, err := bolt.Open(DbFile, 0600, nil)

	if err != nil {
		log.Fatalf("Error opening database: %s", err.Error())
	}

	err = Db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			genesisBlock := NewGenesisBlock()
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