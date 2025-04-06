package common

import (
	"log"
	"sync"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/KevinZonda/GoX/pkg/stringx"
	"github.com/docker/docker/client"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
)

var (
	DockerClient *client.Client
	DockerHelper *dockerProv.DockerHelper
)

func InitDocker() {
	provider := InputProvider()
	d, err := client.NewClientWithOpts(client.WithHost(providerToSocket(provider)), client.WithAPIVersionNegotiation())
	panicx.NotNilErr(err)
	h := dockerProv.NewDockerHelper(d)
	DockerClient = d
	DockerHelper = h
}

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

func InputProvider() dockerProv.Provider {
	provider := rlStrWithPrompt("Please input the provider (docker/podman) (default: docker): ")
	provider = stringx.TrimLower(provider)
	switch provider {
	case "docker", "":
		return dockerProv.ProviderDocker
	case "podman":
		return dockerProv.ProviderPodman
	}
	log.Fatalf("Invalid provider: %s", provider)
	return dockerProv.ProviderDocker
}

var dockerInit sync.Once

func MustInitDocker() {
	dockerInit.Do(InitDocker)
}
