package dockerHelper

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

func (d *DockerHelper) StopContainer(containerID string) error {
	return d.cli.ContainerStop(context.Background(), containerID, container.StopOptions{
		Signal:  "SIGTERM",
		Timeout: nil,
	})
}
