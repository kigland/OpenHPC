package main

import (
	"github.com/kigland/OpenHPC/tools/common"
	"github.com/kigland/OpenHPC/tools/handler"
)

func main() {
	common.InitRL()
	defer common.Rl.Close()
	common.InitDocker()

	handler.Request()
}
