package dockerHelper

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func (d *DockerHelper) ListAllContainers(runningOnly bool) ([]container.Summary, error) {
	if !runningOnly {
		return d.cli.ContainerList(context.Background(), container.ListOptions{
			All: true,
		})
	}

	return d.cli.ContainerList(context.Background(), container.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("status", "running"),
		),
	})
}
