package vm

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/models/apimod"
	"github.com/kigland/OpenHPC/coodinator/shared"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func list(c *gin.Context) {
	vmlists := map[dockerProv.Provider][]apimod.VmListItem{}

	for provider, docker := range shared.Containers {
		uidToContainers, err := docker.UserContainerRelations()
		lists := []apimod.VmListItem{}
		if err != nil {
			log.Println(err)
			continue
		}
		for _, containers := range uidToContainers {
			for svcTagStr, summary := range containers {
				svcTag, err := svcTag.Parse(svcTagStr)
				if err != nil {
					log.Println(err)
					continue
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

				lists = append(lists, item)
			}
		}
		vmlists[provider] = lists
	}

	c.JSON(200, vmlists)
}
