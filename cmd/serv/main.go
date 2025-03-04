package main

import (
	"github.com/KevinZonda/GinTemplate/controller"
	"github.com/KevinZonda/GinTemplate/shared"
	"github.com/KevinZonda/GoX/pkg/iox"
	"github.com/KevinZonda/GoX/pkg/panicx"
)

func initCfg() {
	bs, err := iox.ReadAllByte("config.json")
	panicx.NotNilErr(err)
	panicx.NotNilErr(shared.LoadConfig(bs))
}

func main() {
	initCfg()
	shared.Init()

	controller.Init(shared.Engine)

	shared.RunGin()
}
