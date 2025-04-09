package common

import (
	"fmt"
	"log"

	"github.com/KevinZonda/GoX/pkg/stringx"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
)

func providerToSocket(provider dockerProv.Provider) string {
	switch provider {
	case dockerProv.ProviderDocker:
		return consts.DOCKER_UNIX_SOCKET
	case dockerProv.ProviderPodman:
		return consts.PODMAN_UNIX_SOCKET
	default:
		log.Fatalf("Invalid provider: %s", provider)
		return ""
	}
}

var provider dockerProv.Provider

func SetProvider(p dockerProv.Provider) {
	provider = p
}

func LoadProvider() dockerProv.Provider {
	if dockerProv.ValidateProvider(provider) {
		return provider
	}

	// load from config
	prvStr := cfg.Provider
	if prv, ok := dockerProv.ParseProvider(prvStr); ok {
		provider = prv
		goto success
	}

	provider = inputProvider()
success:
	fmt.Println("Using provider:", provider)
	return provider
}

func inputProvider() dockerProv.Provider {
	prv := rlStrWithPrompt("Please input the provider (docker/podman) (default: docker): ")
	prv = stringx.TrimLower(prv)
	if prv == "" {
		return dockerProv.ProviderDocker
	}
	prvParsed, ok := dockerProv.ParseProvider(prv)
	if !ok {
		log.Fatalf("Invalid provider: %s", prv)
	}
	return prvParsed
}
