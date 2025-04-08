package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)


type Transaction struct {
	ID   []byte
	Vin  []TXInput 
	Vout []TXOutput
}


func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TXInput{
		Txid: 	[]byte{},
		Vout: 	-1,
		ScriptSig: data,
	}
	txOut := TXOutput{
		Value:  50,
		ScriptPubKey: to,
	}

	tx := &Transaction{
		ID:  nil,
		Vin: []TXInput{txIn},
		Vout: []TXOutput{txOut},
	}
	tx.ID = tx.Hash()

	return tx
}


func NewTransaction(from, to string, amount int, bc *Blockchain) *Transaction{
	inputs := []TXInput{}
	outputs := []TXOutput{}

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Not enough funds")
		return nil
	}

	for txid, outs := range validOutputs {
		txid, err := hex.DecodeString(txid)

		if err != nil {
			log.Fatalf("Error decoding txid: %s", err.Error())
		}

		for _, out := range outs {
			input := TXInput{
				Txid:  txid,
				Vout:  out,
				ScriptSig: from,
			}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{Value: amount, ScriptPubKey: to})

	if acc > amount {
		// Change
		outputs = append(outputs, TXOutput{Value: acc - amount, ScriptPubKey: from})
	}

	tx := &Transaction{
		ID:  nil,
		Vin: inputs,
		Vout: outputs,
	}
	tx.ID = tx.Hash()
	return tx
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256([]byte("1"))

	return hash[:]
}

func (coinbase *Transaction) IsCoinbase() bool {
	return len(coinbase.Vin) == 1 && coinbase.Vin[0].Txid == nil && coinbase.Vin[0].Vout == -1
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	
	return true
}


