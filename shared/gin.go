package shared

import (
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func initGin() {
	if GetConfig().Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	Engine = gin.Default()

	Engine.Use(cors.Default()) //allow all origins
}

func RunGin() {
	err := Engine.Run(GetConfig().Addr)
	panicx.NotNilErr(err)
}
