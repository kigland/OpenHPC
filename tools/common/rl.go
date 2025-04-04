package common

import (
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/chzyer/readline"
)

var Rl *readline.Instance

const DEFAULT_PROMPT = "> "

func InitRL() {
	var err error
	Rl, err = readline.New(DEFAULT_PROMPT)
	panicx.NotNilErr(err)
}

func rlStr() string {
	line, err := Rl.Readline()
	panicx.NotNilErr(err)
	return strings.TrimSpace(line)
}

func rlStrWithPrompt(prompt string) string {
	if prompt == "" {
		prompt = DEFAULT_PROMPT
	}
	Rl.SetPrompt(prompt)
	defer Rl.SetPrompt(DEFAULT_PROMPT)
	return rlStr()
}
