package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/controller/ping"
	"github.com/kigland/OpenHPC/coodinator/controller/types"
	"github.com/kigland/OpenHPC/coodinator/controller/user"
	"github.com/kigland/OpenHPC/coodinator/controller/vm"
)

func Init(r gin.IRouter) {
	register(r, &ping.Controller{}, &user.Controller{}, &vm.Controller{})
}

func register(r gin.IRouter, cs ...types.IController) {
	for _, c := range cs {
		c.Init(r)
	}
}
