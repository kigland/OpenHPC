package handler

import (
	"fmt"

	"github.com/kigland/OpenHPC/tools/common"
)

func Upd() {
	cidToFunc(upd)
}

func upd(cid string) {
	info, err := common.Upgrade(cid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info.String())
}
