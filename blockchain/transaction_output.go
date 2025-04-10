package blockchain

import (
	"bytes"
	"log"
)

type TXOutput struct {
	Value        int
	PubKeyHash []byte // unlocking script, which determines the logic of unlocking the output.
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash, err := DecodeBase58(string(address))
	if err != nil {
		log.Fatalf("Failed to decode address: %v", err)
	}
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}