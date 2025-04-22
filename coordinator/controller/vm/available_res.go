package vm

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/coordinator/shared"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/nv"
)

func availableRes(c *gin.Context) {
	c.JSON(200, apimod.VmReqAvailResources{
		Images:    availableResImages(),
		Providers: availableResProviders(),
		Gpus:      availableResGpus(),
	})
}

func availableResProviders() []string {
	providers := []string{}
	for prov := range shared.Containers {
		providers = append(providers, string(prov))
	}
	return providers
}

func availableResImages() []apimod.VmReqAvailImage {
	images := image.AllAvailableImages()
	var avail []apimod.VmReqAvailImage
	for imgName, img := range images {
		avail = append(avail, apimod.VmReqAvailImage{
			Image:       string(imgName),
			Description: img.Description,
			DisplayName: img.DisplayName,
		})
	}
	return avail
}

func availableResGpus() []apimod.VmReqAvailGpu {
	gpus := []apimod.VmReqAvailGpu{}
	smi, err := nv.GetNvidiaSmiLog()
	if err != nil {
		return gpus
	}
	info, err := smi.Parse()
	if err != nil {
		return gpus
	}
	for _, gpu := range info.GPUs {
		gpus = append(gpus, apimod.VmReqAvailGpu{
			GpuId:       strconv.Itoa(gpu.MinorId),
			DisplayName: gpu.Name,
		})
	}
	return gpus
}
