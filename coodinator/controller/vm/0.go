package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/controller/mid"
	"github.com/kigland/OpenHPC/coodinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.POST("/request", mid.FakeAuth, request)
	r.POST("/token", mid.FakeAuth, token)
	r.POST("/extend", mid.FakeAuth, extend)
}
