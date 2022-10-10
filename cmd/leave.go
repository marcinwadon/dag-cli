package cmd

import (
	"dag-cli/domain/layer"
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/node"
	"dag-cli/pkg/pid"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(leaveCmd)
}

var leaveCmd = &cobra.Command{
	Use:   "leave [layer]",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		layerToStop, err := layer.ParseString(args[0])
		if err != nil {
			return err
		}

		cfg, _ := config.LoadConfig()

		p := pid.New(cfg.GetL0PidFilename())
		err = p.Load()
		if err != nil {
			return err
		}

		fmt.Printf("Leaving layer: %s...\n", *layerToStop)

		nodeClient := node.GetClient("127.0.0.1", cfg.L0.Port.CLI)
		if *layerToStop == "l1" {
			nodeClient = node.GetClient("127.0.0.1", cfg.L1.Port.CLI)
		}

		err = nodeClient.Leave()
		if err != nil {
			return err
		}

		fmt.Printf("Layer %s left\n", *layerToStop)
		return nil
	},
}
