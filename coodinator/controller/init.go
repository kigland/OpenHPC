package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/HPC-Scheduler/coodinator/controller/ping"
	"github.com/kigland/HPC-Scheduler/coodinator/controller/types"
)

func Init(r gin.IRouter) {
	register(r, &ping.Controller{})
}

func register(r gin.IRouter, cs ...types.IController) {
	for _, c := range cs {
		c.Init(r)
	}
}
