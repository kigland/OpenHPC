package common

import (
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/chzyer/readline"
)

var Rl *readline.Instance

func InitRL() {
	var err error
	Rl, err = readline.New("> ")
	panicx.NotNilErr(err)
}

func rlStr() string {
	line, err := Rl.Readline()
	panicx.NotNilErr(err)
	return strings.TrimSpace(line)
}
