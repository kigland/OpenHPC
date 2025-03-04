package dockerHelper

import (
	"github.com/docker/docker/client"
)

type DockerHelper struct {
	cli *client.Client
}

func NewDockerHelper(cli *client.Client) *DockerHelper {
	return &DockerHelper{
		cli: cli,
	}
}
