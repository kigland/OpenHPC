package container

import (
	"crypto/rand"
	"math/big"
)

const PREFIX = "KHS"

const DEFAULT_ID_LENGTH = 8

func RndId(length int) string {
	if length <= 0 {
		length = DEFAULT_ID_LENGTH
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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

func NewContainerName(userID string) string {
	return PREFIX + "-" + userID + "-" + RndId(DEFAULT_ID_LENGTH)
}
