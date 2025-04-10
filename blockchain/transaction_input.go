package blockchain

import "bytes"

type TXInput struct {
	Txid      []byte
	Vout      int // inputs must reference previous transactions outputs
	Signature []byte // bitcoin uses ScriptSig for Signature and PubKey
	PubKey    []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}