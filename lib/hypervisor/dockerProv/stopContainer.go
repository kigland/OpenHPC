package dockerProv

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

func (d *DockerHelper) RemoveContainer(containerID string) error {
	return d.cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	})
}
