package shared

import (
	"encoding/json"
	"log"

	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
)

type ProviderConfig struct {
	Provider dockerProv.Provider `json:"provider"`
	Socket   string              `json:"socket"`
}

type Config struct {
	Addr  string `json:"addr"`
	Debug bool   `json:"debug"`

	AvailableProviders []ProviderConfig    `json:"avail_providers"`
	DefaultProvider    dockerProv.Provider `json:"default_provider"`

	BindSSHHost string `json:"bind_ssh_host"`
	BindSSHPort int    `json:"bind_ssh_port"`

	BindHTTPHost string `json:"bind_http_host"`
	BindHTTPPort int    `json:"bind_http_port"`

	MaxPortShift int `json:"max_port_shift"`

	VisitHTTPHost string `json:"visit_http_host"`
	VisitSSHHost  string `json:"visit_ssh_host"`

	MySQL   string `json:"mysql"`
	Storage string `json:"storage"`
}

func (c *Config) normaliseProvider() {
	providers := []ProviderConfig{}
	for _, p := range c.AvailableProviders {
		if dockerProv.ValidateProvider(p.Provider) {
			switch p.Provider {
			case dockerProv.ProviderDocker:
				p.Socket = consts.DOCKER_UNIX_SOCKET
			case dockerProv.ProviderPodman:
				p.Socket = consts.PODMAN_UNIX_SOCKET
			default:
				log.Printf("Unknown provider: %s", p.Provider)
				continue
			}
			providers = append(providers, p)
		}
	}
	c.AvailableProviders = providers

	defaultInProviders := false
	for _, p := range providers {
		if p.Provider == c.DefaultProvider {
			defaultInProviders = true
			break
		}
	}
	if !defaultInProviders {
		if len(providers) > 0 {
			c.DefaultProvider = providers[0].Provider
		} else {
			panic("No provider found")
		}
	}
}

func (c *Config) Normalise() {
	c.normaliseProvider()
}

var cfg *Config

func GetConfig() *Config {
	return cfg
}

func LoadConfig(bs []byte) error {
	err := json.Unmarshal(bs, &cfg)
	if err == nil {
		cfg.Normalise()
	}
	return err
}
