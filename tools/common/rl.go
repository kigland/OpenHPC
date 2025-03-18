package common

import (
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/chzyer/readline"
)

var Rl *readline.Instance

func InitRL() {
	var err error
	Rl, err = readline.New("> ")
	panicx.NotNilErr(err)
}
