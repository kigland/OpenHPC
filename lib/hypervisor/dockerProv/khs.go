package dockerProv

import (
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func (d *DockerHelper) getRelatedName(c container.Summary) (string, svcTag.SvcTag) {
	for _, n := range c.Names {
		tag, err := svcTag.Parse(n)
		if err != nil || tag.Identifier != d.Identifier {
			continue
		}
		return n, tag
	}
	return "", svcTag.SvcTag{}
}
func (d *DockerHelper) AllKHSContainers() (map[string]container.Summary, error) {
	docker := d
	containers, err := docker.ListAllContainers(false)
	if err != nil {
		return nil, err
	}
	cs := map[string]container.Summary{}
	for _, c := range containers {
		n, _ := d.getRelatedName(c)
		cs[n] = c
	}
	return cs, nil
}

func (d *DockerHelper) UserContainerRelations() (map[string]map[string]container.Summary, error) {
	cs, err := d.AllKHSContainers()

	if err != nil {
		return nil, err
	}
	rsh := map[string]map[string]container.Summary{}
	for n, c := range cs {
		tag, err := svcTag.Parse(n)
		if err != nil || tag.Identifier != d.Identifier {
			continue
		}
		userID := tag.Owner
		if _, ok := rsh[userID]; !ok {
			rsh[userID] = map[string]container.Summary{}
		}
		rsh[userID][n] = c
	}
	return rsh, nil
}

func (d *DockerHelper) UserContainers(userID string) (map[string]container.Summary, error) {
	cs, err := d.AllKHSContainers()
	if err != nil {
		return nil, err
	}
	userID = strings.ToLower(userID)
	userCs := map[string]container.Summary{}
	for n, c := range cs {
		tag, err := svcTag.Parse(n)
		if err != nil || tag.Identifier != d.Identifier {
			continue
		}
		if tag.Owner == userID {
			userCs[n] = c
		}
	}
	return userCs, nil
}
