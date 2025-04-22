package vm

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/coordinator/shared"
	"github.com/kigland/OpenHPC/coordinator/utils"
	"github.com/kigland/OpenHPC/lib/image"
)

func upgrade(c *gin.Context) {
	req := utils.BodyAsF[apimod.VmUpgradeReq](c)
	provider, docker := MustGetProviderWithProvId(c, req.Provider)
	if docker == nil {
		utils.ErrorMsg(c, 400, "provider not found")
		return
	}

	summary, ok := docker.TryGetContainer(req.Id)
	if !ok {
		utils.ErrorMsg(c, 404, "container not found")
		return
	}
	inspect, err := docker.ContainerInspect(summary.ID)
	if err != nil {
		utils.ErrorMsg(c, 500, "failed to inspect container")
		return
	}
	ids, err := IDs(docker, summary.ID)
	if err != nil {
		utils.ErrorMsg(c, 500, "failed to get container ids")
		return
	}

	imgStr := image.PruneImageStr(inspect.Config.Image)
	img := image.AllowedImages(imgStr)
	if !img.IsAllowed() {
		utils.ErrorMsg(c, 400, "image not supported")
		log.Println("image not supported", imgStr)
		return
	}

	imgCfg := img.Cfg()

	tokens := filterToken(inspect.Config.Env)
	tokenMap := tokenMap(tokens)
	token := tokenMap[imgCfg.Env.Token]
	if token == "" {
		utils.ErrorMsg(c, 400, "token not found")
		return
	}

	port := -1
	needSSH := false
	for _, p := range summary.Ports {
		if p.PrivatePort == uint16(imgCfg.HTTP) {
			port = int(p.PublicPort)
			break
		}
		if p.PrivatePort == uint16(imgCfg.SSH) {
			needSSH = true
			continue
		}
	}
	if port == -1 {
		utils.ErrorMsg(c, 400, "port not found")
		return
	}

	rdsFrom, rdsTo := "", ""

	for _, m := range inspect.Mounts {
		if strings.Contains(m.Destination, "/rds") {
			rdsFrom = m.Source
			rdsTo = m.Destination
			break
		}
	}

	err = docker.StopContainer(summary.ID)
	if err != nil {
		utils.ErrorMsg(c, 500, "failed to stop container")
		return
	}

	err = docker.RemoveContainer(summary.ID)
	if err != nil {
		utils.ErrorMsg(c, 500, "failed to remove container")
		return
	}

	// fmt.Println("Creating new container...")
	creq := CreateRequest{
		Provider: provider,
		Dk:       docker,
		Image:    img,
		Tag:      ids.SvcTag,
		Passwd:   token,

		BindHost: shared.GetConfig().BindHTTPHost,
		BindPort: port,

		RdsDir:     rdsFrom,
		RdsMountAt: rdsTo,
		ShmSize:    int(req.Shm),

		AllGPU: req.Gpu, // FIXME: ALL GPU

		MaxMemByte: inspect.HostConfig.Resources.Memory,
	}

	if needSSH {
		creq.BindSSHPort = port
		creq.BindSSHHost = shared.GetConfig().BindSSHHost
	}

	info, err := CreateContainerCustomRDS(creq)
	if err != nil {
		utils.ErrorMsg(c, 500, "failed to create container")
		return
	}
	c.JSON(200, createdInfoToVMInfo(info))
}
