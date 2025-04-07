package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(), 
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, 
		Nonce: 0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}



func (b *Block) Serialize() []byte {
	result := &bytes.Buffer{}
	enconder := gob.NewEncoder(result)

	err := enconder.Encode(b)
	if err != nil {
		log.Fatalf("error serializing block: %s", err.Error())
		panic(err)
	}
	return result.Bytes()
}

func Deserialize(data []byte) *Block {
	block := &Block{}
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(block)
	if err != nil {
		log.Fatalf("error deserializing block: %s", err.Error())
		panic(err)
	}
	return block
}