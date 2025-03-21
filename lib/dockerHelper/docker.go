package dockerHelper

import (
	"github.com/docker/docker/client"
)

type DockerHelper struct {
	cli    *client.Client
	Prefix string
}

func NewDockerHelper(cli *client.Client) *DockerHelper {
	return &DockerHelper{
		cli:    cli,
		Prefix: "KHS",
	}
}

func (d *DockerHelper) Cli() *client.Client {
	return d.cli
}
