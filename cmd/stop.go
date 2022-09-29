package cmd

import (
	"dag-cli/pkg/config"
	"dag-cli/pkg/layer"
	"dag-cli/pkg/node"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use: "stop [layer]",
	Short: "",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		layerToStop, err := layer.ParseString(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Stopping layer: %s...\n", layerToStop)

		cfg, _ := config.LoadConfig()
		err = node.Stop(cfg, layerToStop)
		if err != nil {
			return err
		}

		fmt.Printf("Layer %s stopped\n", layerToStop)
		return nil
	},
}
