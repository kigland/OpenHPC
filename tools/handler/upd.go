package handler

import (
	"fmt"
	"os"

	"github.com/kigland/OpenHPC/tools/common"
)

func Upd() {
	if len(os.Args) == 3 {
		upd(os.Args[2])
		return
	}
	upd(common.InputWithPrompt("Container ID or Service Tag or Short Code:"))
}

func upd(cid string) {
	info, err := common.Upgrade(cid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info.String())
}
