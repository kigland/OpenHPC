package utils

import "github.com/KevinZonda/GoX/pkg/panicx"

func RdrErr[T any](t T, e error) T {
	panicx.NotNilErr(e)
	return t
}
