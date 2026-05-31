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

type ContainerSummaryWithSvcTag struct {
	container.Summary
	*svcTag.SvcTag
}

func (d *DockerHelper) TryGetContainer(cid string) (ContainerSummaryWithSvcTag, bool) {
	cs, err := d.AllKHSContainers()
	if err != nil {
		return ContainerSummaryWithSvcTag{}, false
	}
	var tag *svcTag.SvcTag
	if strings.Contains(cid, "@") {
		_tag, err := svcTag.Parse(cid)
		if err != nil {
			return ContainerSummaryWithSvcTag{}, false
		}
		tag = &_tag
		cid = tag.String()
	}
	for n, c := range cs {
		if n == cid || n == "/"+cid || strings.HasPrefix(c.ID, cid) {
			return ContainerSummaryWithSvcTag{
				Summary: c,
				SvcTag:  tag,
			}, true
		}
	}
	return ContainerSummaryWithSvcTag{}, false
}
