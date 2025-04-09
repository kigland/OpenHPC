package common

import (
	"log"
	"strings"

	"github.com/KevinZonda/GoX/pkg/ruby"
	"github.com/KevinZonda/GoX/pkg/stringx"
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
	inspect := ruby.RdrErr(DockerHelper.ContainerInspect(summary.ID))
	return inspect.Config.Env
}
