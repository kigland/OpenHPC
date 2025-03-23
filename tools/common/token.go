package common

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/KevinZonda/GoX/pkg/ruby"
	"github.com/KevinZonda/GoX/pkg/stringx"
	"github.com/kigland/HPC-Scheduler/lib/image"
)

func Token(cid string) []string {
	env := Env(cid)
	return filterToken(env)
}

func filterToken(env []string) []string {
	var tokens []string
	for _, e := range env {
		if strings.Contains(stringx.TrimLower(e), "token") {
			tokens = append(tokens, e)
		}
	}
	return tokens
}

func tokenMap(env []string) map[string]string {
	tokenMap := make(map[string]string)
	for _, token := range env {
		parts := strings.Split(token, "=")
		if len(parts) == 2 {
			tokenMap[parts[0]] = parts[1]
		}
	}
	return tokenMap
}

func Env(cid string) []string {
	summary, ok := DockerHelper.TryGetContainer(cid)
	if !ok {
		log.Fatalf("Container not found or not managed by KHS")
		return nil
	}
	inspect := ruby.RdrErr(DockerHelper.Cli().ContainerInspect(context.Background(), summary.ID))
	return inspect.Config.Env
}

func Upgrade(cid string) (ContainerInfo, error) {
	summary, ok := DockerHelper.TryGetContainer(cid)
	if !ok {
		return ContainerInfo{}, fmt.Errorf("container not found or not managed by KHS")
	}
	inspect := ruby.RdrErr(DockerHelper.Cli().ContainerInspect(context.Background(), summary.ID))
	ids := IDs(summary.ID)
	imgStr := inspect.Config.Image
	img := image.AllowedImages(imgStr)
	if !slices.Contains(image.ALLOWED_IMAGES, img) {
		return ContainerInfo{}, fmt.Errorf("image not supported")
	}

	tokens := filterToken(inspect.Config.Env)
	tokenMap := tokenMap(tokens)
	token := tokenMap[image.JUPYTER_TOKEN]
	if token == "" {
		return ContainerInfo{}, fmt.Errorf("token not found")
	}

	port := -1
	for _, p := range summary.Ports {
		if p.PrivatePort == 8888 {
			port = int(p.PublicPort)
			break
		}
	}
	if port == -1 {
		return ContainerInfo{}, fmt.Errorf("port not found")
	}

	rdsFrom, rdsTo := "", ""

	for _, m := range inspect.Mounts {
		if strings.Contains(m.Destination, "/rds") {
			rdsFrom = m.Source
			rdsTo = m.Destination
			break
		}
	}

	return CreateContainerCustomRDS(DockerHelper, img, ids.SvcTag.Owner, token, port, ids.SvcTag.Project, rdsFrom, rdsTo)
}
