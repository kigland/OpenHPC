package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/controller/mid"
	"github.com/kigland/OpenHPC/coordinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.POST("/vm/request", mid.MustAuth, request)
	r.POST("/vm/token", mid.MustAuth, token)
	r.GET("/vm/list", mid.MustAuth, list)
	r.POST("/vm/del", mid.MustAuth, del)
	r.POST("/vm/extend", mid.MustAuth, extend)
}
