package handler

import (
	"fmt"
	"os"

	"github.com/kigland/OpenHPC/tools/common"
)

func IDs() {
	if len(os.Args) == 3 {
		ids(os.Args[2])
		return
	}
	ids(common.InputWithPrompt("Container ID or Service Tag or Short Code:"))
}

func ids(cid string) {
	v := common.IDs(cid)
	fmt.Println(v.String())
}
