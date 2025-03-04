package ping

import (
	"github.com/KevinZonda/GinTemplate/controller/types"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
