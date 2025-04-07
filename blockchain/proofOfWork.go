package blockchain

import (
	"blockchain/utils"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)


const maxNonce = math.MaxInt64
const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}


func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	
	target.Lsh(target, uint(256-targetBits)) 
	
	pow := &ProofOfWork{block: block, target: target}
	return pow
}

func (p *ProofOfWork) prepareData(nonce int64) []byte{
	data := append([]byte{}, p.block.PrevBlockHash...)
	data = append(data, p.block.Data...)
	data = append(data, utils.IntToHex(p.block.Timestamp)...)
	data = append(data, utils.IntToHex(int64(targetBits))...)
	data = append(data, utils.IntToHex(nonce)...)
	return data
}

func (p *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Started mining new block")
	for nonce < maxNonce {
		data := p.prepareData(int64(nonce))

		hash = sha256.Sum256(data)
		// only to print hash
		if math.Remainder(float64(nonce), 100000) == 0 {
			fmt.Printf("\r%x", hash)
		}
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 {
			break
		}else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(int64(pow.block.Nonce))
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
	
}