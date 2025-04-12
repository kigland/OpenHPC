package user

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.POST("/user/register", register)
	r.POST("/user/login", login)
	r.GET("/user/quota", quota)
}
