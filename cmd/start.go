package cmd

import (
	"dag-cli/pkg/config"
	"dag-cli/pkg/layer"
	"dag-cli/pkg/node"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use: "start [layer]",
	Short: "",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		layerToRun, err := layer.ParseString(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Starting layer: %s...\n", layerToRun)

		cfg, _ := config.LoadConfig()
		err = node.Start(cfg, layerToRun)
		if err != nil {
			return err
		}

		fmt.Printf("Layer %s started\n", layerToRun)
		return nil
	},
}
