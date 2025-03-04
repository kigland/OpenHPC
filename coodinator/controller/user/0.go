package user

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/HPC-Scheduler/coodinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/quota", quota)
}
