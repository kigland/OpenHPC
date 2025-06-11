package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/controller/infra"
	"github.com/kigland/OpenHPC/coordinator/controller/ping"
	"github.com/kigland/OpenHPC/coordinator/controller/stat"
	"github.com/kigland/OpenHPC/coordinator/controller/types"
	"github.com/kigland/OpenHPC/coordinator/controller/vm"
)

func Init(r gin.IRouter) {
	register(r, &ping.Controller{}, &vm.Controller{},
		&stat.Controller{}, &infra.Controller{}) // &user.Controller{},
}

func register(r gin.IRouter, cs ...types.IController) {
	for _, c := range cs {
		c.Init(r)
	}
}
