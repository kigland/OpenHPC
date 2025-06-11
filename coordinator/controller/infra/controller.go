package infra

import (
	"log"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/controller/mid"
	"github.com/kigland/OpenHPC/coordinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.GET("/infra/px/restart", mid.ACLAuth, func(c *gin.Context) {
		o, err := exec.Command("systemctl", "restart", "frpc-px").Output()
		if err != nil {
			c.JSON(500, gin.H{"error": "running error"})
			return
		}
		log.Println(o)
		c.JSON(200, gin.H{"message": "ok"})
	})
}
