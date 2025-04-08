package blockchain

type TXOutput struct {
	Value        int
	ScriptPubKey string // unlocking script, which determines the logic of unlocking the output.
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}