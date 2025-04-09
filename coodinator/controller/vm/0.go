package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/controller/mid"
	"github.com/kigland/OpenHPC/coodinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.POST("/vm/request", mid.FakeAuth, request)
	r.POST("/vm/token", mid.FakeAuth, token)
	r.GET("/vm/list", mid.FakeAuth, list)
	r.POST("/vm/del", mid.FakeAuth, del)

	r.POST("/vm/extend", mid.FakeAuth, extend)
}
