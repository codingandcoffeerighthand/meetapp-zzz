package config

import (
	"os"
	ws_config "proxy-srv/config/ws"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerConfig ServerConfig `yaml:"server"`
	ProxyConfig  ProxyConfig  `yaml:"proxy"`
}
type ServerConfig struct {
	Port     int `yaml:"port"`
	GrpcPort int `yaml:"grpc_port"`
}
type ProxyConfig struct {
	Url string `yaml:"url"`
}

type ConfigPath string

func NewConfig(filePath ConfigPath) (Config, error) {
	var (
		configBytes = ws_config.DefaultConfig
		config      = Config{}
		err         error
	)
	if filePath != "" {
		configBytes, err = os.ReadFile(string(filePath))
		if err != nil {
			return config, err
		}
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
