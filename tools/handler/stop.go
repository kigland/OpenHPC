package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/kigland/HPC-Scheduler/lib/utils"
	"github.com/kigland/HPC-Scheduler/tools/common"
)

func Stop() {
	cid := common.InputWithPrompt("Container ID or Service Tag:")
	summary, ok := tryGetContainer(cid)
	if !ok {
		panic("Container not found or not managed by KHS")
	}
	fmt.Println(summary.Names, summary.ID, summary.Status, summary.Image, summary.ImageID, summary.Created, summary.State)
	confirm := common.InputWithPrompt("Are you sure to stop the container? (y/n)")
	if confirm != "y" {
		fmt.Println("Stopping container cancelled")
		return
	}

	common.DockerHelper.Cli().ContainerStop(context.Background(), summary.ID, container.StopOptions{})
}

func tryGetContainer(cid string) (container.Summary, bool) {
	cs := utils.RdrErr(common.DockerHelper.AllKHSContainers())
	for n, c := range cs {
		if n == cid || strings.HasPrefix(c.ID, cid) {
			return c, true
		}
	}
	return container.Summary{}, false
}
