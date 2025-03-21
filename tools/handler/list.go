package handler

import (
	"fmt"
	"strings"

	"github.com/disiqueira/gotree"
	"github.com/docker/docker/api/types/container"
	"github.com/kigland/HPC-Scheduler/lib/svcTag"
	"github.com/kigland/HPC-Scheduler/lib/utils"
	"github.com/kigland/HPC-Scheduler/tools/common"
)

func List() {
	uidToContainers := utils.RdrErr(common.DockerHelper.UserContainerRelations())
	tree := SummaryToTree(uidToContainers, false)
	fmt.Println(tree.Print())
}

func containerToStr(c container.Summary, tag string, showCID bool) string {
	var ports []string
	for _, p := range c.Ports {
		ports = append(ports, fmt.Sprintf(":%d->%s:%d", p.PrivatePort, p.IP, p.PublicPort))
	}
	mount := ""
	for _, m := range c.Mounts {
		mount += fmt.Sprintf(" %s:%s", m.Source, m.Destination)
		if m.RW {
			mount += "(RW)"
		} else {
			mount += "(RO)"
		}
	}
	tagName := tag
	svcTag, err := svcTag.Parse(tag)
	if err == nil {
		tagName = svcTag.ShortName()
	}
	mount = strings.TrimSpace(mount)
	if showCID {
		return fmt.Sprintf("[%s] %s %s %s CID: %s", tagName, c.Status, strings.Join(ports, ", "), mount, c.ID)
	}
	return fmt.Sprintf("[%s] %s %s %s", tagName, c.Status, strings.Join(ports, ", "), mount)

}

func SummaryToTree(uidToContainers map[string]map[string]container.Summary, showCID bool) gotree.Tree {
	tree := gotree.New("Users")
	for uid, containers := range uidToContainers {
		user := tree.Add(uid)
		for svcTag, c := range containers {
			user.Add(containerToStr(c, svcTag, showCID))
		}
	}

	return tree
}
func ListUser(u string) {
	uidToContainers := utils.RdrErr(common.DockerHelper.UserContainerRelations())
	u = strings.TrimSpace(u)
	if u == "" {
		panic("user is empty")
	}
	u = strings.ToLower(u)
	if _, ok := uidToContainers[u]; !ok {
		panic("user not found")
	}
	tree := gotree.New(u)
	for svcTag, c := range uidToContainers[u] {
		tree.Add(containerToStr(c, svcTag, false))
	}
	fmt.Println(tree.Print())
}
