package common

import (
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/docker/docker/client"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
)

var (
	DockerClient *client.Client
	DockerHelper *dockerProv.DockerHelper
)

func InitDocker() {
	d, err := client.NewClientWithOpts(client.WithHost(consts.PODMAN_UNIX_SOCKET), client.WithAPIVersionNegotiation())
	panicx.NotNilErr(err)
	h := dockerProv.NewDockerHelper(d)
	DockerClient = d
	DockerHelper = h
}
