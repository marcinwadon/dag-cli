package config

import (
	"github.com/spf13/viper"
)

type KeyConfig struct {
	Path     string
	Alias    string
	Password string
}

type TessellationConfig struct {
	Version string
}

type GithubConfig struct {
	Organization     string
	Repository       string
	L0Filename       string `mapstructure:"l0-filename"`
	L1Filename       string `mapstructure:"l1-filename"`
	SeedlistFilename string `mapstructure:"seedlist-filename"`
}

type JavaConfig struct {
	Xmx string
	Xms string
	Xss string
}

type PortConfig struct {
	Public int
	P2P    int
	CLI    int
}

type L0Peer struct {
	Id   string
	Host string
	Port int
}

type L0Config struct {
	Env           string
	Path          string
	PidPath       string `mapstructure:"pid-path"`
	SeedlistPath  string `mapstructure:"seedlist-path"`
	SnapshotsPath string `mapstructure:"snapshots-path"`
	Java          JavaConfig
	LoadBalancer  string `mapstructure:"load-balancer"`
	Port          PortConfig
}

type L1Config struct {
	Env          string
	Path         string
	PidPath      string `mapstructure:"pid-path"`
	L0Peer       L0Peer `mapstructure:"l0-peer"`
	Java         JavaConfig
	LoadBalancer string `mapstructure:"load-balancer"`
	Port         PortConfig
}

type Config struct {
	Key          KeyConfig
	Tessellation TessellationConfig
	ExternalIP   string `mapstructure:"external-ip"`
	Github       GithubConfig
	L0           L0Config
	L1           L1Config
}

func LoadConfig() (Config, error) {
	var C Config
	err := viper.Unmarshal(&C)
	if err != nil {
		return C, err
	}
	return C, nil
}
