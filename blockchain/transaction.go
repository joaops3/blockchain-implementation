package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/jinzhu/copier"
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
		Signature: nil,
		PubKey: []byte(to),
	}
	txOut :=  *NewTXOutput(10, to)

	tx := &Transaction{
		ID:  nil,
		Vin: []TXInput{txIn},
		Vout: []TXOutput{txOut},
	}
	tx.ID = tx.Hash()

	return tx
}


func NewTransaction(wallet *Wallet, to string, amount int, bc *Blockchain) *Transaction{
	inputs := []TXInput{}
	outputs := []TXOutput{}

	acc, validOutputs := bc.FindSpendableOutputs(string(wallet.PublicKey), amount)

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
				PubKey: wallet.PublicKey,
				Signature: nil,
			}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, to))

	if acc > amount {
		// Change
		outputs = append(outputs,  *NewTXOutput(acc - amount, string(wallet.GetAddress())))
	}

	tx := &Transaction{
		ID:  nil,
		Vin: inputs,
		Vout: outputs,
	}
	tx.ID = tx.Hash()
	bc.SignTransaction(tx, DeserializePrivateKey(wallet.PrivateKey))
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


func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	txCopy := &Transaction{}
	err := copier.Copy(txCopy, tx)
	if err != nil {
		log.Fatalf("Error copying transaction: %s", err.Error())
	}

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Fatalf("Error signing transaction: %s", err.Error())
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}
}


