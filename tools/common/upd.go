package common

import (
	"fmt"
	"slices"
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/KevinZonda/GoX/pkg/ruby"
	"github.com/kigland/OpenHPC/lib/image"
)

func Upgrade(cid string) (ContainerInfo, error) {
	summary, ok := DockerHelper.TryGetContainer(cid)
	if !ok {
		return ContainerInfo{}, fmt.Errorf("container not found or not managed by OHPC")
	}
	inspect := ruby.RdrErr(DockerHelper.ContainerInspect(summary.ID))
	ids := IDs(summary.ID)
	imgStr := inspect.Config.Image
	imgStr = strings.TrimPrefix(imgStr, "docker.io/")
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
	needSSH := false
	for _, p := range summary.Ports {
		if p.PrivatePort == 8888 {
			port = int(p.PublicPort)
			break
		}
		if p.PrivatePort == 22 {
			needSSH = true
			continue
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

	fmt.Println("Stopping container...")
	panicx.NotNilErr(DockerHelper.StopContainer(summary.ID))
	fmt.Println("Removing container...")
	panicx.NotNilErr(DockerHelper.RemoveContainer(summary.ID))

	fmt.Println("Creating new container...")
	return CreateContainerCustomRDS(DockerHelper, img, ids.SvcTag, token, port, rdsFrom, rdsTo, needSSH)
}
