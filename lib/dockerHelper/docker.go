package dockerHelper

import (
	"github.com/docker/docker/client"
	"github.com/kigland/HPC-Scheduler/lib/consts"
)

type DockerHelper struct {
	cli    *client.Client
	Prefix string
}

func NewDockerHelper(cli *client.Client) *DockerHelper {
	return &DockerHelper{
		cli:    cli,
		Prefix: consts.IDENTIFIER,
	}
}

func (d *DockerHelper) Cli() *client.Client {
	return d.cli
}
