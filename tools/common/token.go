package common

import (
	"context"
	"log"
	"strings"

	"github.com/kigland/HPC-Scheduler/lib/utils"
)

func Token(cid string) []string {
	env := Env(cid)
	var tokens []string
	for _, e := range env {
		if strings.Contains(utils.TrimLower(e), "token") {
			tokens = append(tokens, e)
		}
	}
	return tokens
}

func Env(cid string) []string {
	summary, ok := DockerHelper.TryGetContainer(cid)
	if !ok {
		log.Fatalf("Container not found or not managed by KHS")
		return nil
	}
	inspect := utils.RdrErr(DockerHelper.Cli().ContainerInspect(context.Background(), summary.ID))
	return inspect.Config.Env
}
