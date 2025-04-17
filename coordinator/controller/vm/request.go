package vm

import (
	"log"
	"strconv"
	"strings"

	"github.com/KevinZonda/GoX/pkg/randx"
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/coordinator/shared"
	"github.com/kigland/OpenHPC/coordinator/utils"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func request(c *gin.Context) {
	req := utils.BodyAsF[apimod.VmReq](c)

	if req.Image == "" {
		req.Image = string(image.ImageJupyterHub)
	}

	provider, docker := MustGetProviderWithProvId(c, req.Provider)
	if docker == nil {
		return
	}

	passwd := utils.RndId(32)

	rndPort := randx.RndRange(0, shared.GetConfig().MaxPortShift)

	imgName := image.AllowedImages(req.Image)
	if !imgName.IsAllowed() {
		c.JSON(400, gin.H{
			"message": "image not supported",
		})
		return
	}

	img := image.Factory{
		Password:    passwd,
		BindHost:    shared.GetConfig().BindHTTPHost,
		BindPort:    shared.GetConfig().BindHTTPPort + rndPort,
		BindSSHHost: shared.GetConfig().BindSSHHost,
		BindSSHPort: shared.GetConfig().BindSSHPort + rndPort,
		Provider:    provider,
	}.Image(imgName).WithGPU(1).
		WithAutoRestart().
		WithBaseURL(imgName.BaseURLEnvVar(), consts.BASE_URL(shared.GetConfig().BindHTTPPort+rndPort))

	if req.Shm >= 64 {
		img = img.WithShmSize(int64(req.Shm))
	}

	// RDS Support
	rdsMountAt := ""
	if req.EnableRds {
		var err error
		rdsDir, err := shared.Rds.GetRDSPath(req.Owner, req.RdsFolder)
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
		defer func() {
			err := docker.RemoveContainer(id)
			if err != nil {
				log.Fatalf("Failed to remove container: %v", err)
			}
		}()
		return
	}

	cinfo := apimod.VmCreatedInfo{
		Cid:   id,
		Image: string(imgName),
		RdsAt: rdsMountAt,
		Token: passwd,
		Http:  parseHTTPVisitURL(shared.GetConfig().BindHTTPPort + rndPort),

		SvcTag: svgT.String(),
		Sc:     svgT.ShortCode(),
	}

	if imgName.SupportSSH() {
		cinfo.Ssh = shared.GetConfig().VisitSSHHost + ":" + strconv.Itoa(shared.GetConfig().BindSSHPort+rndPort)
	}

	c.JSON(200, cinfo)
}

func parseHTTPVisitURL(newPort int) string {
	return strings.ReplaceAll(shared.GetConfig().VisitHTTPHost, "$PORT", strconv.Itoa(newPort))
}
