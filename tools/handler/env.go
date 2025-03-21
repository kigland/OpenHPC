package handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/kigland/HPC-Scheduler/tools/common"
)

func Env() {
	if len(os.Args) == 3 {
		env(os.Args[2])
		return
	}
	env(common.InputWithPrompt("Container ID:"))
}

func env(cid string) {
	env := common.Env(cid)
	envs := strings.Join(env, "\n")
	fmt.Println(envs)
}

func Token() {
	if len(os.Args) == 3 {
		token(os.Args[2])
		return
	}
	token(common.InputWithPrompt("Container ID:"))
}

func token(cid string) {
	tokens := common.Token(cid)
	tokensStr := strings.Join(tokens, "\n")
	fmt.Println(tokensStr)
}
