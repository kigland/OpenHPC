package utils

import (
	"crypto/rand"
	"math/big"
)

const DEFAULT_ID_LENGTH = 8
const ALPHABET_NUMERIC = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RndId(length int) string {
	if length <= 0 {
		length = DEFAULT_ID_LENGTH
	}
	return RndIdCharset(length, ALPHABET_NUMERIC)
}

func RndIdCharset(length int, charset string) string {
	if length <= 0 {
		return ""
	}
	charsetLength := big.NewInt(int64(len(charset)))

	b := make([]byte, length)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			b[i] = ALPHABET_NUMERIC[0]
			continue
		}
		b[i] = charset[randomIndex.Int64()]
	}

	return string(b)
}
