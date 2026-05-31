package dockerProv

import (
	"context"
	"log"
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
	svcTag.SvcTag
}

func (d *DockerHelper) TryGetContainer(cid string) (ContainerSummaryWithSvcTag, bool) {
	cs, err := d.AllKHSContainers()
	if err != nil {
		return ContainerSummaryWithSvcTag{}, false
	}
	var tag svcTag.SvcTag
	if strings.Contains(cid, "@") {
		tag, err = svcTag.Parse(cid)
		if err != nil {
			return ContainerSummaryWithSvcTag{}, false
		}
		cid = tag.String()
	}
	for n, c := range cs {
		if n == cid || n == "/"+cid || strings.HasPrefix(c.ID, cid) {
			log.Println("Container found for id:", cid, "with svcTag:", c.Labels["svcTag"])
			return ContainerSummaryWithSvcTag{
				Summary: c,
				SvcTag:  tag,
			}, true
		}
	}
	return ContainerSummaryWithSvcTag{}, false
}
