package cmd

import (
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/lb"
	"dag-cli/infrastructure/node"
	"errors"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()

		version := args[0]

		if version == "latest" {
			l0Lb := lb.GetClient(cfg.L0.LoadBalancer)
			randomPeer, err := l0Lb.GetRandomReadyPeer()
			if err != nil {
				return err
			}
			randomNodeClient := node.GetClient(randomPeer.Ip, randomPeer.PublicPort)
			peerInfo, err := randomNodeClient.GetNodeInfo()
			if err != nil {
				return err
			}
			version = "v" + peerInfo.Version
		}

		if semver.IsValid(version) {
			return node.Upgrade(cfg, version)
		} else {
			return errors.New("invalid version")
		}
	},
}
