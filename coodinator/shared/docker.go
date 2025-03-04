package shared

import (
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/docker/docker/client"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
)

var Docker *client.Client
var DockerHelper *dockerHelper.DockerHelper

func initDocker() {
	var err error
	Docker, err = client.NewClientWithOpts(client.WithHost(GetConfig().DockerHost), client.WithAPIVersionNegotiation())
	panicx.NotNilErr(err)
	DockerHelper = dockerHelper.NewDockerHelper(Docker)
}
