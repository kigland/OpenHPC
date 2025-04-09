package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/shared"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
)

func GetProvider(prov string) *dockerProv.DockerHelper {
	_, docker := GetProviderWithProvId(prov)
	return docker
}

func MustGetProvider(c *gin.Context, prov string) *dockerProv.DockerHelper {
	_, docker := MustGetProviderWithProvId(c, prov)
	return docker
}

func MustGetProviderWithProvId(c *gin.Context, prov string) (dockerProv.Provider, *dockerProv.DockerHelper) {
	id, docker := GetProviderWithProvId(prov)
	if docker == nil {
		c.JSON(400, gin.H{
			"message": "provider not found",
		})
		return "", nil
	}
	return id, docker
}

func GetProviderWithProvId(prov string) (dockerProv.Provider, *dockerProv.DockerHelper) {
	provider, _ := dockerProv.ParseProvider(prov)
	if provider == "" {
		provider = shared.GetDefaultProvider()
	}

	return provider, shared.Containers[provider]
}
