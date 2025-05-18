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

type Web3Config struct {
	RpcUrl          string `yaml:"rpc_url"`
	ChanId          int64  `yaml:"chan_id"`
	PrivateKey      string `yaml:"private_key"`
	ContractAddress string `yaml:"contract_address"`
}

type CloudflareConfig struct {
	BaseURL   string `yaml:"base_url"`
	AppId     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
	TurnId    string `yaml:"turn_id"`
	TurnToken string `yaml:"turn_token"`
}

type WsService struct {
	Url string `yaml:"url"`
}

type Config struct {
	ServerConfig     ServerConfig     `yaml:"server"`
	CloudflareConfig CloudflareConfig `yaml:"cloudflare"`
	Web3Config       Web3Config       `yaml:"web3"`
	WsService        WsService        `yaml:"ws"`
	PrivateKeyRSA    string           `yaml:"private_key_rsa"`
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
