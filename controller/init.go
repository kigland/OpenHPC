package controller

import (
	"github.com/KevinZonda/GinTemplate/controller/ping"
	"github.com/KevinZonda/GinTemplate/controller/types"
	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	register(r, &ping.Controller{})
}

func register(r gin.IRouter, cs ...types.IController) {
	for _, c := range cs {
		c.Init(r)
	}
}
