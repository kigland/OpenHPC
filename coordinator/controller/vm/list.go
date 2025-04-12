package vm

import (
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/coordinator/shared"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func list(c *gin.Context) {
	vmlists := []apimod.VmListProvider{}

	for provider, docker := range shared.Containers {
		uidToContainers, err := docker.UserContainerRelations()
		lists := []apimod.VmListItem{}
		if err != nil {
			log.Println(err)
			continue
		}
		for _, containers := range uidToContainers {
			for svcTagStr, summary := range containers {
				item, ok := containInfo(svcTagStr, summary)
				if !ok {
					continue
				}
				lists = append(lists, item)
			}
		}
		vmlists = append(vmlists, apimod.VmListProvider{
			Provider:   string(provider),
			Containers: lists,
		})
	}

	c.JSON(200, vmlists)
}

func containInfo(svcTagStr string, summary container.Summary) (apimod.VmListItem, bool) {
	svcTag, err := svcTag.Parse(svcTagStr)
	if err != nil {
		log.Println(err)
		return apimod.VmListItem{}, false
	}
	item := apimod.VmListItem{
		Cid:     summary.ID,
		SvcTag:  svcTag.String(),
		Sc:      svcTag.ShortCode(),
		Status:  summary.Status,
		Owner:   svcTag.Owner,
		Project: svcTag.Project,
	}

	ports := []apimod.VmListItemMount{}
	for _, port := range summary.Ports {
		ports = append(ports, apimod.VmListItemMount{
			Host:      fmt.Sprintf("%s:%d", port.IP, port.PublicPort),
			Container: fmt.Sprintf(":%d", port.PrivatePort),
		})
	}
	item.Port = ports

	mounts := []apimod.VmListItemMount{}
	for _, mount := range summary.Mounts {
		mounts = append(mounts, apimod.VmListItemMount{
			Host:      mount.Source,
			Container: mount.Destination,
			Readonly:  !mount.RW,
		})
	}
	item.Mount = mounts
	return item, true
}
