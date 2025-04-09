package shared

import (
	"log"

	"github.com/docker/docker/client"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
	"github.com/kigland/OpenHPC/lib/rds"
)

var Containers map[dockerProv.Provider]*dockerProv.DockerHelper = map[dockerProv.Provider]*dockerProv.DockerHelper{}

var Rds *rds.RDS
var DefaultProvider dockerProv.Provider

func GetDefaultProvider() dockerProv.Provider {
	if DefaultProvider == "" {
		for k, _ := range Containers {
			DefaultProvider = k
			return DefaultProvider
		}
		return ""
	}
	return DefaultProvider
}

func _initContainer(p ProviderConfig) error {
	dk, err := client.NewClientWithOpts(client.WithHost(p.Socket), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	Containers[p.Provider] = dockerProv.NewDockerHelper(dk)
	return nil
}

func initDocker() {
	for _, p := range GetConfig().AvailableProviders {
		err := _initContainer(p)
		if err != nil {
			log.Println("init container controller failed:", p.Provider, err)
		}
	}

	if GetConfig().DefaultProvider != "" {
		DefaultProvider = GetConfig().DefaultProvider
	}

	Rds = &rds.RDS{
		BasePath: consts.RDS_PATH,
	}
}
