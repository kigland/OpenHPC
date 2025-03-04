package shared

import "encoding/json"

type Config struct {
	Addr  string `json:"addr"`
	Debug bool   `json:"debug"`
}

var cfg *Config

func GetConfig() *Config {
	return cfg
}

func LoadConfig(bs []byte) error {
	return json.Unmarshal(bs, &cfg)
}
