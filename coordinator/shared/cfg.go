package shared

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/KevinZonda/GoX/pkg/stringx"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
	"github.com/kigland/OpenHPC/lib/image"
)

type ProviderConfig struct {
	Provider dockerProv.Provider `json:"provider"`
	Socket   string              `json:"socket"`
}

type ACL struct {
	AllowAll bool     `json:"allow_all"`
	APIKeys  []string `json:"api_keys"`
}

type ACLParsed struct {
	Raw     ACL
	apikeys map[string][]string
}

func (a ACLParsed) GetACLAllowList(k string) []string {
	if a.Raw.AllowAll {
		return []string{""}
	}
	return a.apikeys[k]
}

type Config struct {
	Addr  string `json:"addr"`
	Debug bool   `json:"debug"`

	AvailableProviders []ProviderConfig    `json:"available_providers"`
	DefaultProvider    dockerProv.Provider `json:"default_provider"`

	ACL       ACL       `json:"acl"`
	ACLParsed ACLParsed `json:"-"`

	BindSSHHost string `json:"bind_ssh_host"`
	BindSSHPort int    `json:"bind_ssh_port"`

	BindHTTPHost string `json:"bind_http_host"`
	BindHTTPPort int    `json:"bind_http_port"`

	MaxPortShift int `json:"max_port_shift"`

	VisitHTTPHost string `json:"visit_http_host"`
	VisitSSHHost  string `json:"visit_ssh_host"`

	MySQL   string `json:"mysql"`
	Storage string `json:"storage"`

	Images []image.HPCImage `json:"images"`
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
			log.Printf("Recognised provider: %s", p.Provider)
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

func (c *Config) parseACL() {
	acl := c.ACL
	mapList := make(map[string][]string)
	for _, apikey := range acl.APIKeys {
		parts := strings.SplitN(apikey, ":", 2)
		left := ""
		var right []string
		if len(parts) <= 0 {
			continue
		}

		left = parts[0]
		if len(parts) == 1 {
			right = []string{""}
		} else {
			right = stringx.TrimAll(strings.Split(parts[1], ","))
		}
		mapList[left] = right
	}
	c.ACLParsed = ACLParsed{
		Raw:     acl,
		apikeys: mapList,
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
		cfg.parseACL()
		cfg.Normalise()
	}
	image.InitAllowedImages(cfg.Images)
	return err
}
