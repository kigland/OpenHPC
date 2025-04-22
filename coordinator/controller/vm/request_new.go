package vm

import (
	"fmt"
	"slices"

	"github.com/KevinZonda/GoX/pkg/randx"
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/coordinator/shared"
	"github.com/kigland/OpenHPC/coordinator/utils"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func requestNew(c *gin.Context) {
	req := utils.BodyAsF[apimod.VmReq](c)

	if len(image.ALLOWED_IMAGES) == 0 {
		utils.ErrorMsg(c, 500, "no image supported")
		return
	}

	if req.Image == "" {
		req.Image = string(image.ALLOWED_IMAGES[0])
	}

	provider, docker := MustGetProviderWithProvId(c, req.Provider)
	if docker == nil {
		utils.ErrorMsg(c, 400, "provider not found")
		return
	}

	imgName := image.AllowedImages(req.Image)
	if !imgName.IsAllowed() {
		c.JSON(400, gin.H{
			"message": "image not supported",
		})
		return
	}

	rndPort := randx.RndRange(0, shared.GetConfig().MaxPortShift)

	creq := CreateRequest{
		Provider: provider,
		Dk:       docker,
		Image:    imgName,
		Tag:      svcTag.New(req.Owner).WithProject(req.Project),
		Passwd:   utils.RndId(32),

		BindHost:    shared.GetConfig().BindHTTPHost,
		BindPort:    shared.GetConfig().BindHTTPPort + rndPort,
		BindSSHHost: shared.GetConfig().BindSSHHost,
		BindSSHPort: shared.GetConfig().BindSSHPort + rndPort,

		MaxMemByte: int64(req.MaxMem) * 1024 * 1024, // B -> KB -> MB
	}

	if req.EnableGpu {
		if req.Gpu.All {
			creq.AllGPU = true
		} else {
			if len(req.Gpu.Ids) > 0 {
				creq.GPUIds = make([]int, len(req.Gpu.Ids))
				for i, id := range req.Gpu.Ids {
					creq.GPUIds[i] = int(id)
				}
			}
		}
	}

	if err := validateGPUIds(creq.GPUIds); err != nil {
		utils.ErrorMsg(c, 400, err.Error())
		return
	}

	if req.Shm >= 64 {
		creq.ShmSize = int(req.Shm)
	}

	// RDS Support
	if req.EnableRds {
		var err error
		rdsDir, err := shared.Rds.GetRDSPath(req.Owner, req.RdsFolder, true)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "failed to get rds path",
			})
			return
		}
		creq.RdsDir = rdsDir
		creq.RdsMountAt = imgName.RdsDir()
	}

	info, err := CreateContainerCustomRDS(creq)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to create container",
		})
		return
	}

	c.JSON(200, createdInfoToVMInfo(info))
}

func validateGPUIds(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	ids = slices.Compact(ids)
	slices.Sort(ids)
	for _, id := range ids {
		if id < 0 {
			return fmt.Errorf("invalid gpu id: %d", id)
		}
	}
	// TODO: Verify GPU Ids
	return nil
}
