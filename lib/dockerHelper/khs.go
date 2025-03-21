package dockerHelper

import (
	"strings"

	"github.com/docker/docker/api/types/container"
)

func (d *DockerHelper) AllKHSContainers() (map[string]container.Summary, error) {
	docker := d
	containers, err := docker.ListAllContainers(false)
	if err != nil {
		return nil, err
	}
	prefix := d.Prefix + "-"
	cs := map[string]container.Summary{}
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.HasPrefix(n, prefix) || strings.HasPrefix(n, "/"+prefix) {
				cs[n] = c
				break
			}
		}
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
		names := strings.Split(n, "-")
		if len(names) != 3 {
			continue
		}
		userID := names[1]
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
	prefix := d.Prefix + "-"
	userID = strings.ToLower(userID)
	userCs := map[string]container.Summary{}
	for n, c := range cs {
		if strings.HasPrefix(n, prefix+userID+"-") || strings.HasPrefix(n, "/"+prefix+userID+"-") {
			userCs[n] = c
		}
	}
	return userCs, nil
}
