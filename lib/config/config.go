package config

import (
	"encoding/json"
	"fmt"

	"github.com/KevinZonda/GoX/pkg/iox"
)

type Config struct {
	Provider string `json:"provider"`
}

var DEFAULT_CONFIG_PATHS = []string{
	"/etc/openhpc/config.json",
	"/etc/khs/config.json",
}

func LoadConfigFromDefaultPaths() (Config, error) {
	return LoadConfig(DEFAULT_CONFIG_PATHS...)
}

func DefaultConfig() Config {
	return Config{
		Provider: "",
	}
}

var ErrNoConfig = fmt.Errorf("no config found")

// LoadConfig loads config from config file
// If load failed, return default and error
// If load success, return config and nil
func LoadConfig(cfg_path ...string) (Config, error) {
	cfg := DefaultConfig()
	for _, path := range cfg_path {
		bs, err := iox.ReadAllBytes(path)
		if err != nil {
			continue
		}
		if err = json.Unmarshal(bs, &cfg); err != nil {
			continue
		}
		return cfg, nil
	}
	return cfg, ErrNoConfig
}
