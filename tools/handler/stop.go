package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/kigland/HPC-Scheduler/tools/common"
)

func Stop() {
	if len(os.Args) == 3 {
		stop(os.Args[2])
		return
	}
	stop(common.InputWithPrompt("Container ID or Service Tag or Short Code:"))
}

func stop(cid string) {
	summary, ok := common.DockerHelper.TryGetContainer(cid)
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
