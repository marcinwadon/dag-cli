package config

import (
	"dag-cli/pkg/color"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"path/filepath"
)

type KeyConfig struct {
	Path     string
	Alias    string
	Password string
}

func (cfg *KeyConfig) Show() {
	fmt.Printf("--- Key configuration:\n")
	fmt.Printf("Key: %s, alias=***, password=***\n", color.Ize(color.Green, cfg.Path))
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

func (cfg *GithubConfig) Show() {
	fmt.Printf("--- GitHub configuration:\n")
	fmt.Printf("Github: %s/%s\n", color.Ize(color.Green, cfg.Organization), color.Ize(color.Green, cfg.Repository))
	fmt.Printf("L0: %s, L1: %s, Seedlist: %s\n", color.Ize(color.Green, cfg.L0Filename), color.Ize(color.Green, cfg.L1Filename), color.Ize(color.Green, cfg.SeedlistFilename))
}

type JavaConfig struct {
	Xmx string
	Xms string
	Xss string
}

func (cfg *JavaConfig) Show() {
	fmt.Printf("Java xmx=%s xms=%s xss=%s\n", color.Ize(color.Green, cfg.Xmx), color.Ize(color.Green, cfg.Xms), color.Ize(color.Green, cfg.Xss))
}

type PortConfig struct {
	Public int
	P2P    int
	CLI    int
}

func (cfg *PortConfig) Show() {
	fmt.Printf("Public=%s P2P=%s CLI=%s\n", color.Ize(color.Green, fmt.Sprintf("%d", cfg.Public)), color.Ize(color.Green, fmt.Sprintf("%d", cfg.P2P)), color.Ize(color.Green, fmt.Sprintf("%d", cfg.CLI)))
}

type L0Peer struct {
	Discovery bool
	Id        string
	Host      string
	Port      int
}

func (cfg *L0Peer) Show() {
	if cfg.Discovery {
		fmt.Printf("L0 joining peer: auto-discovery\n")
	} else {
		fmt.Printf("L0 joining peer: %s:%s id=%s\n", color.Ize(color.Green, cfg.Host), color.Ize(color.Green, fmt.Sprintf("%d", cfg.Port)), color.Ize(color.Green, cfg.Id))
	}
}

type L0Config struct {
	SnapshotsPath string `mapstructure:"snapshots-path"`
	Java          JavaConfig
	LoadBalancer  string `mapstructure:"load-balancer"`
	Port          PortConfig
}

func (cfg *L0Config) Show() {
	fmt.Printf("--- L0 configuration:\n")
	fmt.Printf("Snapshots path: %s\n", color.Ize(color.Green, cfg.SnapshotsPath))
	fmt.Printf("Load-balancer: %s\n", color.Ize(color.Green, cfg.LoadBalancer))
	cfg.Java.Show()
	cfg.Port.Show()
}

type L1Config struct {
	L0Peer       L0Peer `mapstructure:"l0-peer"`
	Java         JavaConfig
	LoadBalancer string `mapstructure:"load-balancer"`
	Port         PortConfig
}

func (cfg *L1Config) Show() {
	fmt.Printf("--- L1 configuration:\n")
	fmt.Printf("Load-balancer: %s\n", color.Ize(color.Green, cfg.LoadBalancer))
	cfg.Java.Show()
	cfg.Port.Show()
	cfg.L0Peer.Show()
}

type Config struct {
	Env           string
	BlockExplorer string `mapstructure:"block-explorer"`
	Verbose       bool
	WorkingDir    string `mapstructure:"working-directory"`
	Key           KeyConfig
	Tessellation  TessellationConfig
	ExternalIP    string `mapstructure:"external-ip"`
	Github        GithubConfig
	L0            L0Config
	L1            L1Config
}

func (cfg *Config) GetL0JarFilename() string {
	return filepath.Join(cfg.WorkingDir, "l0.jar")
}
func (cfg *Config) GetL1JarFilename() string {
	return filepath.Join(cfg.WorkingDir, "l1.jar")
}
func (cfg *Config) GetL0PidFilename() string {
	return filepath.Join(cfg.WorkingDir, "l0.pid")
}
func (cfg *Config) GetL1PidFilename() string {
	return filepath.Join(cfg.WorkingDir, "l1.pid")
}
func (cfg *Config) GetL0SeedlistFilename() string {
	return filepath.Join(cfg.WorkingDir, "seedlist.csv")
}
func (cfg *Config) Show() {
	fmt.Printf("--------------------------------------\n")
	fmt.Printf("--- Configuration:\n")
	fmt.Printf("Environment: %s\n", color.Ize(color.Green, cfg.Env))
	fmt.Printf("Working dir: %s\n", color.Ize(color.Green, cfg.WorkingDir))
	fmt.Printf("External IP: %s\n", color.Ize(color.Green, cfg.ExternalIP))
	fmt.Printf("Block explorer: %s\n", color.Ize(color.Green, cfg.BlockExplorer))
	cfg.Key.Show()
	cfg.Github.Show()
	cfg.L0.Show()
	cfg.L1.Show()
	fmt.Printf("--------------------------------------\n")
}

func LoadConfig() (Config, error) {
	var C Config
	err := viper.Unmarshal(&C, func(config *mapstructure.DecoderConfig) {
		config.ErrorUnset = true
	})
	if err != nil {
		return C, err
	}
	return C, nil
}

func (cfg *Config) Persist() error {
	viper.Set("env", cfg.Env)
	viper.Set("working-directory", cfg.WorkingDir)
	viper.Set("external-ip", cfg.ExternalIP)
	viper.Set("block-explorer", cfg.BlockExplorer)
	viper.Set("tessellation.version", cfg.Tessellation.Version)
	viper.Set("key.path", cfg.Key.Path)
	viper.Set("key.alias", cfg.Key.Alias)
	viper.Set("key.password", cfg.Key.Password)
	viper.Set("github.organization", cfg.Github.Organization)
	viper.Set("github.repository", cfg.Github.Repository)
	viper.Set("github.l0-filename", cfg.Github.L0Filename)
	viper.Set("github.l1-filename", cfg.Github.L1Filename)
	viper.Set("github.seedlist-filename", cfg.Github.SeedlistFilename)
	viper.Set("l0.snapshots-path", cfg.L0.SnapshotsPath)
	viper.Set("l0.load-balancer", cfg.L0.LoadBalancer)
	viper.Set("l0.java.xmx", cfg.L0.Java.Xmx)
	viper.Set("l0.java.xms", cfg.L0.Java.Xms)
	viper.Set("l0.java.xss", cfg.L0.Java.Xss)
	viper.Set("l0.port.public", cfg.L0.Port.Public)
	viper.Set("l0.port.p2p", cfg.L0.Port.P2P)
	viper.Set("l0.port.cli", cfg.L0.Port.CLI)
	viper.Set("l1.load-balancer", cfg.L1.LoadBalancer)
	viper.Set("l1.java.xmx", cfg.L1.Java.Xmx)
	viper.Set("l1.java.xms", cfg.L1.Java.Xms)
	viper.Set("l1.java.xss", cfg.L1.Java.Xss)
	viper.Set("l1.port.public", cfg.L1.Port.Public)
	viper.Set("l1.port.p2p", cfg.L1.Port.P2P)
	viper.Set("l1.port.cli", cfg.L1.Port.CLI)
	viper.Set("l1.l0-peer.discovery", cfg.L1.L0Peer.Discovery)
	viper.Set("l1.l0-peer.id", cfg.L1.L0Peer.Id)
	viper.Set("l1.l0-peer.host", cfg.L1.L0Peer.Host)
	viper.Set("l1.l0-peer.port", cfg.L1.L0Peer.Port)

	return viper.WriteConfig()
}

func DefaultMainnet() Config {
	return Config{
		Env:          "mainnet",
		Verbose:      false,
		WorkingDir:   "./",
		BlockExplorer: "be-mainnet.constellationnetwork.io",
		Key:          KeyConfig{Path: "./key.p12", Alias: "alias", Password: ""},
		Tessellation: TessellationConfig{Version: ""},
		ExternalIP:   "",
		Github: GithubConfig{
			Organization:     "Constellation-Labs",
			Repository:       "tessellation",
			L0Filename:       "cl-node.jar",
			L1Filename:       "cl-dag-l1.jar",
			SeedlistFilename: "mainnet-seedlist",
		},
		L0: L0Config{
			SnapshotsPath: "./snapshots",
			Java: JavaConfig{
				Xmx: "10G",
				Xms: "1024M",
				Xss: "394K",
			},
			LoadBalancer: "https://l0-lb-mainnet.constellationnetwork.io",
			Port: PortConfig{
				Public: 9000,
				P2P:    9001,
				CLI:    9002,
			},
		},
		L1: L1Config{
			L0Peer: L0Peer{
				Discovery: true,
				Port:      9000,
			},
			Java: JavaConfig{
				Xmx: "4G",
				Xms: "1024M",
				Xss: "394K",
			},
			LoadBalancer: "https://l1-lb-mainnet.constellationnetwork.io",
			Port: PortConfig{
				Public: 9010,
				P2P:    9011,
				CLI:    9012,
			},
		},
	}
}

func DefaultTestnet() Config {
	return Config{
		Env:          "testnet",
		Verbose:      false,
		WorkingDir:   "./",
		BlockExplorer: "be-testnet.constellationnetwork.io",
		Key:          KeyConfig{Path: "./key.p12", Alias: "alias", Password: ""},
		Tessellation: TessellationConfig{Version: ""},
		ExternalIP:   "",
		Github: GithubConfig{
			Organization:     "Constellation-Labs",
			Repository:       "tessellation",
			L0Filename:       "cl-node.jar",
			L1Filename:       "cl-dag-l1.jar",
			SeedlistFilename: "mainnet-seedlist",
		},
		L0: L0Config{
			SnapshotsPath: "./snapshots",
			Java: JavaConfig{
				Xmx: "10G",
				Xms: "1024M",
				Xss: "394K",
			},
			LoadBalancer: "https://l0-lb-testnet.constellationnetwork.io",
			Port: PortConfig{
				Public: 9000,
				P2P:    9001,
				CLI:    9002,
			},
		},
		L1: L1Config{
			L0Peer: L0Peer{
				Discovery: true,
				Port:      9000,
			},
			Java: JavaConfig{
				Xmx: "4G",
				Xms: "1024M",
				Xss: "394K",
			},
			LoadBalancer: "https://l1-lb-testnet.constellationnetwork.io",
			Port: PortConfig{
				Public: 9010,
				P2P:    9011,
				CLI:    9012,
			},
		},
	}
}
