package utils

import (
	"strconv"

	"github.com/mr-tron/base58"
)

func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

func EncodeBase58(data []byte) string {
	return base58.Encode(data)
}


func DecodeBase58(encoded string) ([]byte, error) {
	return base58.Decode(encoded)
}