package common

import (
	"fmt"

	"github.com/kigland/OpenHPC/lib/config"
)

var cfg config.Config

func InitConfig() {
	var err error
	cfg, err = config.LoadConfigFromDefaultPaths()
	if err != nil {
		fmt.Println("Failed to load config, fallback to default config")
	}
}
