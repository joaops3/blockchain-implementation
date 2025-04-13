package blockchain

import (
	"blockchain/utils"
	"bytes"
	"encoding/gob"
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
	pubKeyHash, err := utils.DecodeBase58(string(address))
	if err != nil {
		log.Fatalf("Failed to decode address: %v", err)
	}
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (out *TXOutput) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(out)
	if err != nil {
		log.Fatalf("Failed to serialize TXOutput: %v", err)
	}
	return result.Bytes()
}

func SerializeOutputs(outputs []TXOutput) []byte {
    var result bytes.Buffer
    encoder := gob.NewEncoder(&result)

    err := encoder.Encode(outputs)
    if err != nil {
        log.Fatalf("Failed to serialize TXOutputs: %v", err)
    }
    return result.Bytes()
}

func DeserializeTXOutput(data []byte) []TXOutput {
    var txo []TXOutput

    decoder := gob.NewDecoder(bytes.NewReader(data))
    err := decoder.Decode(&txo)
    if err != nil {
        log.Fatalf("Failed to deserialize TXOutput: %v", err)
    }
    return txo
}