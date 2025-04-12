package stat

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/controller/types"
)

type Controller struct{}

var _ types.IController = (*Controller)(nil)

func (c *Controller) Init(r gin.IRouter) {
	r.GET("/stat/nvidia-smi", NvidiaSMIHandler)
	r.GET("/stat/cpu", CPUHandler)
	r.GET("/stat/mem", MemHandler)
}
