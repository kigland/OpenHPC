package shared

import "encoding/json"

type Config struct {
	Addr  string `json:"addr"`
	Debug bool   `json:"debug"`

	Storage string `json:"storage"`

	DockerHost string `json:"docker_host"`

	MySQL string `json:"mysql"`
}

var cfg *Config

func GetConfig() *Config {
	return cfg
}

func LoadConfig(bs []byte) error {
	return json.Unmarshal(bs, &cfg)
}
