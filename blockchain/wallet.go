package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey []byte
    PublicKey  []byte
}


func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}
func newKeyPair() ([]byte, []byte) {
    curve := elliptic.P256()
    private, err := ecdsa.GenerateKey(curve, rand.Reader)
    if err != nil {
        panic(err)
    }
    privateKeyBytes := private.D.Bytes()
    pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

    return privateKeyBytes, pubKey
}

func DeserializePrivateKey(data []byte) ecdsa.PrivateKey {
    curve := elliptic.P256()
    privateKey := ecdsa.PrivateKey{
        PublicKey: ecdsa.PublicKey{Curve: curve},
        D:         new(big.Int).SetBytes(data),
    }
    privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(data)
    return privateKey
}

func (w *Wallet) SerializePrivateKey() []byte {
    return w.PrivateKey
}
func (w *Wallet) GetAddress() []byte {
    pubkeyHash := HashPubKey(w.PublicKey)
    versionedPayload := append([]byte{version}, pubkeyHash...)
    checkSum := checkSum(versionedPayload)

    fullPayload := append(versionedPayload, checkSum...)
    address := EncodeBase58(fullPayload)
    return []byte(address)
}

func HashPubKey(pubKey []byte) []byte {
	first := sha256.Sum256(pubKey)
	second := sha256.Sum256(first[:])
	return second[:]
}

func checkSum(payload []byte) []byte {
	firstSha := sha256.Sum256(payload)
	secondSha := sha256.Sum256(firstSha[:])
	return secondSha[:addressChecksumLen]

}