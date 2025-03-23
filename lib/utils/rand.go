package utils

import (
	"github.com/KevinZonda/GoX/pkg/randx"
)

const DEFAULT_ID_LENGTH = 8

func RndId(length int) string {
	if length <= 0 {
		length = DEFAULT_ID_LENGTH
	}
	return randx.RndIdCharset(length, randx.ALPHABET_NUMERIC)
}
