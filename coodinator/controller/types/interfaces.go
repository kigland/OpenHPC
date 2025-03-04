package types

import "github.com/gin-gonic/gin"

type IController interface {
	Init(r gin.IRouter)
}
