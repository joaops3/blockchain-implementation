package blockchain

type TXInput struct {
	Txid      []byte
	Vout      int // inputs must reference previous transactions outputs
	ScriptSig string // script which provides data to be used in an output’s ScriptPubKey
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}