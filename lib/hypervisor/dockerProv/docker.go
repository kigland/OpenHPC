package dockerProv

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/kigland/OpenHPC/lib/consts"
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

func (d *DockerHelper) ContainerInspect(cid string) (container.InspectResponse, error) {
	return d.cli.ContainerInspect(context.Background(), cid)
}
