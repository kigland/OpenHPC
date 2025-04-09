package handler

import (
	"fmt"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/kigland/OpenHPC/tools/common"
)

func Stop() {
	cidToFunc(stop)
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

	panicx.NotNilErr(common.DockerHelper.StopContainer(summary.ID))
	panicx.NotNilErr(common.DockerHelper.RemoveContainer(summary.ID))
}
