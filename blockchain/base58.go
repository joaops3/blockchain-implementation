package blockchain

import "github.com/mr-tron/base58"


func EncodeBase58(data []byte) string {
	return base58.Encode(data)
}


func DecodeBase58(encoded string) ([]byte, error) {
	return base58.Decode(encoded)
}