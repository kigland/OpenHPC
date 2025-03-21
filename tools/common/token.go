package common

import (
	"context"
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
	inspect := utils.RdrErr(DockerHelper.Cli().ContainerInspect(context.Background(), cid))
	return inspect.Config.Env
}
