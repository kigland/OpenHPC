package handler

import (
	"os"

	"github.com/kigland/OpenHPC/tools/common"
)

func PopFst() string {
	if len(os.Args) == 0 {
		return ""
	}
	arg := os.Args[0]
	os.Args = os.Args[1:]
	return arg
}

func cidToFunc(f func(cid string)) {
	cid := PopFst()
	if cid != "" {
		f(cid)
		return
	}
	f(common.InputWithPrompt("Container ID:"))
}
