package utils

import (
	"crypto/rand"
	"math/big"
)

func RndId(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	charsetLength := big.NewInt(int64(len(charset)))

	b := make([]byte, length)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			b[i] = charset[0]
			continue
		}
		b[i] = charset[randomIndex.Int64()]
	}

	return string(b)
}
