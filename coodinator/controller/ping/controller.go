package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
