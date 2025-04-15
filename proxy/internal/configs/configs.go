package configs

import (
	"os"
	"proxy-srv/config"

	"gopkg.in/yaml.v3"
)

type ConfigPath string

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

type CloudflareConfig struct {
	BaseURL   string `yaml:"base_url"`
	AppId     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
	TurnId    string `yaml:"turn_id"`
	TurnToken string `yaml:"turn_token"`
}

type Config struct {
	ServerConfig     ServerConfig     `yaml:"server"`
	CloudflareConfig CloudflareConfig `yaml:"cloudflare"`
}

func NewConfig(filePath ConfigPath) (Config, error) {
	var (
		configBytes = config.DefaultConfig
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
