package common

import (
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/docker/docker/client"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/dockerHelper"
)

var (
	DockerClient *client.Client
	DockerHelper *dockerHelper.DockerHelper
)

func InitDocker() {
	d, err := client.NewClientWithOpts(client.WithHost(consts.PODMAN_UNIX_SOCKET), client.WithAPIVersionNegotiation())
	panicx.NotNilErr(err)
	h := dockerHelper.NewDockerHelper(d)
	DockerClient = d
	DockerHelper = h
}
