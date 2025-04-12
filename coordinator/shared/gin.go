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

	c := cors.DefaultConfig()
	c.AllowAllOrigins = true
	c.AddAllowHeaders("Authorization")

	Engine.Use(cors.New(c))
}

func RunGin() {
	err := Engine.Run(GetConfig().Addr)
	panicx.NotNilErr(err)
}
