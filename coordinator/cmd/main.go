package main

import (
	"os"

	"github.com/KevinZonda/GoX/pkg/iox"
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/kigland/OpenHPC/coordinator/controller"
	"github.com/kigland/OpenHPC/coordinator/shared"
)

func initCfg() {
	path := "/etc/openhpc/host.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	bs, err := iox.ReadAllByte(path)
	panicx.NotNilErr(err)
	panicx.NotNilErr(shared.LoadConfig(bs))
}

func main() {
	initCfg()
	shared.Init()

	controller.Init(shared.Engine)

	shared.RunGin()
}
