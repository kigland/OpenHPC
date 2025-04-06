package dockerProv

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/kigland/OpenHPC/lib/svcTag"
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

func (d *DockerHelper) TryGetContainer(cid string) (container.Summary, bool) {
	cs, err := d.AllKHSContainers()
	if err != nil {
		return container.Summary{}, false
	}
	if strings.Contains(cid, "@") {
		svcTag, err := svcTag.Parse(cid)
		if err != nil {
			return container.Summary{}, false
		}
		cid = svcTag.String()
	}
	for n, c := range cs {
		if n == cid || n == "/"+cid || strings.HasPrefix(c.ID, cid) {
			return c, true
		}
	}
	return container.Summary{}, false
}
