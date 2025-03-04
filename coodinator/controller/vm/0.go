package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/HPC-Scheduler/coodinator/controller/mid"
	"github.com/kigland/HPC-Scheduler/coodinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.POST("/request", mid.MustAuth, request)
	r.POST("/extend", mid.MustAuth, extend)
}
