package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/controller/mid"
	"github.com/kigland/OpenHPC/coordinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.POST("/vm/request", mid.ACLAuth, requestNew)
	r.POST("/vm/token", mid.ACLAuth, token)
	r.GET("/vm/list", mid.ACLAuth, list)
	r.POST("/vm/del", mid.ACLAuth, del)
	r.POST("/vm/extend", mid.ACLAuth, extend)
	r.POST("/vm/upgrade", mid.ACLAuth, upgrade)
}
