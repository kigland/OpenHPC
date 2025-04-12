package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/controller/mid"
	"github.com/kigland/OpenHPC/coodinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	x := r.Use(mid.ACLAuth)
	{
		x.POST("/vm/request", request)
		x.POST("/vm/token", token)
		x.GET("/vm/list", list)
		x.POST("/vm/del", del)
		x.POST("/vm/extend", extend)
	}
}
