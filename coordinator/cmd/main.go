package main

import (
	"github.com/KevinZonda/GoX/pkg/iox"
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/kigland/OpenHPC/coordinator/controller"
	"github.com/kigland/OpenHPC/coordinator/shared"
)

func initCfg() {
	bs, err := iox.ReadAllByte("/etc/openhpc/host.json")
	panicx.NotNilErr(err)
	panicx.NotNilErr(shared.LoadConfig(bs))
}

func main() {
	initCfg()
	shared.Init()

	controller.Init(shared.Engine)

	shared.RunGin()
}
