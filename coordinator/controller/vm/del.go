package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/coordinator/utils"
)

func del(c *gin.Context) {
	req := utils.BodyAsF[apimod.VmDelReq](c)

	docker := MustGetProvider(c, req.Provider)
	if docker == nil {
		utils.ErrorMsg(c, 400, "provider not found")
		return
	}

	summary, ok := docker.TryGetContainer(req.Id)
	if !ok {
		utils.ErrorMsg(c, 404, "Container not found")
		return
	}

	err := docker.StopContainer(summary.ID)
	if err != nil {
		utils.ErrorMsg(c, 500, "Failed to stop container")
		return
	}

	err = docker.RemoveContainer(summary.ID)
	if err != nil {
		utils.ErrorMsg(c, 500, "Failed to remove container")
		return
	}
	c.Status(200)
}
