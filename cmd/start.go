package cmd

import (
	"dag-cli/domain/layer"
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/lb"
	"dag-cli/infrastructure/node"
	"errors"
	"fmt"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/spf13/cobra"
)

var skipVersionCheck bool

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVar(&skipVersionCheck, "skip-version-check", false, "skip checking cluster tessellation version")
}

var startCmd = &cobra.Command{
	Use:   "start [layer]",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		layerToRun, err := layer.ParseString(args[0])
		if err != nil {
			return err
		}

		cfg, _ := config.LoadConfig()

		if !semver.IsValid(cfg.Tessellation.Version) {
			return errors.New("invalid tessellation version; run upgrade")
		}

		if !skipVersionCheck {
			err = checkClusterVersion(cfg, layerToRun)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Starting layer: %s...\n", *layerToRun)

		err = node.Start(cfg, *layerToRun)
		if err != nil {
			return err
		}

		fmt.Printf("Layer %s started\n", *layerToRun)
		return nil
	},
}

func checkClusterVersion(cfg config.Config, layerToRun *layer.Layer) error {
	lbClient := lb.GetClient(cfg.L0.LoadBalancer)
	if *layerToRun == layer.L1 {
		lbClient = lb.GetClient(cfg.L1.LoadBalancer)
	}

	peer, err := lbClient.GetRandomReadyPeer()
	if err != nil {
		return err
	}
	nodeClient := node.GetClient(peer.Ip, peer.PublicPort)

	nodeInfo, err := nodeClient.GetNodeInfo()
	if err != nil {
		return err
	}

	if fmt.Sprintf("v%s", nodeInfo.Version) != cfg.Tessellation.Version {
		return fmt.Errorf("cluster version: v%s does not match the local version: %s", nodeInfo.Version, cfg.Tessellation.Version)
	}

	return nil
}