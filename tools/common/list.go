package common

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
)

func AllKHSContainers() (map[string]container.Summary, error) {
	docker := DockerHelper
	containers, err := docker.ListAllContainers(false)
	if err != nil {
		return nil, err
	}
	cs := map[string]container.Summary{}
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.HasPrefix(n, "KHS-") {
				cs[n] = c
				break
			}
		}
	}
	return cs, nil
}

func UserContainerRelations() (map[string]map[string]container.Summary, error) {
	cs, err := AllKHSContainers()
	fmt.Println(cs, err)

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

func UserContainers(userID string) (map[string]container.Summary, error) {
	cs, err := AllKHSContainers()
	if err != nil {
		return nil, err
	}
	userCs := map[string]container.Summary{}
	for n, c := range cs {
		if strings.HasPrefix(n, "KHS-"+userID+"-") {
			userCs[n] = c
		}
	}
	return userCs, nil
}
