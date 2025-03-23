package common

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/KevinZonda/GoX/pkg/ruby"
	"github.com/docker/docker/api/types/container"
	"github.com/kigland/HPC-Scheduler/lib/image"
)

func Upgrade(cid string) (ContainerInfo, error) {
	summary, ok := DockerHelper.TryGetContainer(cid)
	if !ok {
		return ContainerInfo{}, fmt.Errorf("container not found or not managed by KHS")
	}
	inspect := ruby.RdrErr(DockerHelper.Cli().ContainerInspect(context.Background(), summary.ID))
	ids := IDs(summary.ID)
	imgStr := inspect.Config.Image
	img := image.AllowedImages(imgStr)
	if !slices.Contains(image.ALLOWED_IMAGES, img) {
		return ContainerInfo{}, fmt.Errorf("image not supported")
	}

	tokens := filterToken(inspect.Config.Env)
	tokenMap := tokenMap(tokens)
	token := tokenMap[image.JUPYTER_TOKEN]
	if token == "" {
		return ContainerInfo{}, fmt.Errorf("token not found")
	}

	port := -1
	for _, p := range summary.Ports {
		if p.PrivatePort == 8888 {
			port = int(p.PublicPort)
			break
		}
	}
	if port == -1 {
		return ContainerInfo{}, fmt.Errorf("port not found")
	}

	rdsFrom, rdsTo := "", ""

	for _, m := range inspect.Mounts {
		if strings.Contains(m.Destination, "/rds") {
			rdsFrom = m.Source
			rdsTo = m.Destination
			break
		}
	}

	panicx.NotNilErr(DockerHelper.Cli().ContainerStop(context.Background(), summary.ID, container.StopOptions{}))

	fmt.Println("Waiting for container to stop...")
	time.Sleep(time.Second * 4)

	fmt.Println("Creating new container...")
	return CreateContainerCustomRDS(DockerHelper, img, ids.SvcTag, token, port, rdsFrom, rdsTo)
}
