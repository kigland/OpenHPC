package vm

import (
	"log"
	"strconv"

	"github.com/KevinZonda/GoX/pkg/randx"
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/models/openapi"
	"github.com/kigland/OpenHPC/coodinator/shared"
	"github.com/kigland/OpenHPC/coodinator/utils"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func request(c *gin.Context) {
	req := utils.BodyAsF[openapi.VmReq](c)

	provider, docker := MustGetProviderWithProvId(c, req.Provider)
	if docker == nil {
		return
	}

	passwd := utils.RndId(32)

	rndPort := randx.RndRange(0, shared.GetConfig().MaxPortShift)

	imgName := image.ImageJupyterHub

	img := image.Factory{
		Password:    passwd,
		BindHost:    shared.GetConfig().BindHTTPHost,
		BindPort:    shared.GetConfig().BindHTTPPort + rndPort,
		BindSSHHost: shared.GetConfig().BindSSHHost,
		BindSSHPort: shared.GetConfig().BindSSHPort + rndPort,
		Provider:    provider,
	}.Image(imgName).WithGPU(1).WithAutoRestart()

	// RDS Support
	rdsMountAt := ""
	if req.EnableRDS {
		var err error
		rdsDir, err := shared.Rds.GetRDSPath(req.Owner, req.RDSFolder)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "failed to get rds path",
			})
			return
		}
		rdsMountAt = imgName.RdsDir()
		img = img.WithMountRW(rdsDir, rdsMountAt)
	}

	svgT := svcTag.New(req.Owner).WithProject(req.Project)
	img.ContainerName = svgT.String()

	id, err := docker.StartContainer(img, false)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
		return
	}

	cinfo := openapi.VmCreatedInfo{
		CID:   id,
		RDSAt: rdsMountAt,
		Token: passwd,
		SSH:   shared.GetConfig().VisitSSHHost + ":" + strconv.Itoa(shared.GetConfig().BindSSHPort+rndPort),
		HTTP:  shared.GetConfig().VisitHTTPHost + ":" + strconv.Itoa(shared.GetConfig().BindHTTPPort+rndPort),

		SvcTag:    svgT.String(),
		ShortCode: svgT.ShortCode(),
	}

	c.JSON(200, cinfo)
}
