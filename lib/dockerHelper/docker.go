package dockerHelper

import (
	"github.com/docker/docker/client"
	"github.com/kigland/HPC-Scheduler/lib/consts"
)

type DockerHelper struct {
	cli        *client.Client
	Identifier string
}

func NewDockerHelper(cli *client.Client) *DockerHelper {
	return &DockerHelper{
		cli:        cli,
		Identifier: consts.IDENTIFIER,
	}
}

func (d *DockerHelper) Cli() *client.Client {
	return d.cli
}

func ShortId(cid string) string {
	if len(cid) > 12 {
		return cid[:12]
	}
	return cid
}
