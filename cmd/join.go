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
	rootCmd.AddCommand(joinCmd)
}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		layerToRun, err := layer.ParseString(args[0])
		if err != nil {
			return err
		}

		cfg, _ := config.LoadConfig()

		p := pid.New(cfg.GetL0PidFilename())
		err = p.Load()
		if err != nil {
			return err
		}

		fmt.Printf("Joining layer: %s...\n", *layerToRun)

		err = node.Join(cfg, *layerToRun)
		if err != nil {
			return err
		}

		fmt.Printf("Layer %s joined\n", *layerToRun)
		return nil
	},
}
