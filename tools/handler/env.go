package handler

import (
	"fmt"
	"strings"

	"github.com/kigland/OpenHPC/tools/common"
)

func Env() {
	cidToFunc(env)
}

func env(cid string) {
	env := common.Env(cid)
	envs := strings.Join(env, "\n")
	fmt.Println(envs)
}

func Token() {
	cidToFunc(token)
}

func token(cid string) {
	tokens := common.Token(cid)
	tokensStr := strings.Join(tokens, "\n")
	fmt.Println(tokensStr)
}

func IDs() {
	cidToFunc(ids)
}

func ids(cid string) {
	v := common.IDs(cid)
	fmt.Println(v.String())
}
