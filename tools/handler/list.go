package handler

import (
	"fmt"
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/disiqueira/gotree"
	"github.com/docker/docker/api/types/container"
	"github.com/kigland/HPC-Scheduler/tools/common"
)

func List() {
	uidToContainers, err := common.UserContainerRelations()
	panicx.NotNilErr(err)
	tree := SummaryToTree(uidToContainers)
	fmt.Println(tree.Print())
}

func contrainerToStr(c container.Summary, svcTag string) string {
	var ports []string
	for _, p := range c.Ports {
		ports = append(ports, fmt.Sprintf("%d->%s:%d", p.PrivatePort, p.IP, p.PublicPort))
	}
	mount := ""
	for _, m := range c.Mounts {
		mount += fmt.Sprintf(" %s:%s", m.Source, m.Destination)
		if m.RW {
			mount += " (RW)"
		} else {
			mount += " (RO)"
		}
	}
	mount = strings.TrimSpace(mount)

	return fmt.Sprintf("[%s] %s CID: %s\n%s %s", svcTag, c.Status, c.ID, strings.Join(ports, ", "), mount)
}

func SummaryToTree(uidToContainers map[string]map[string]container.Summary) gotree.Tree {
	tree := gotree.New("Users")
	for uid, containers := range uidToContainers {
		user := tree.Add(uid)
		for svcTag, c := range containers {
			user.Add(contrainerToStr(c, svcTag))
		}
	}

	return tree
}
func ListUser(u string) {
	uidToContainers, err := common.UserContainerRelations()
	panicx.NotNilErr(err)
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
		tree.Add(contrainerToStr(c, svcTag))
	}
	fmt.Println(tree.Print())
}
