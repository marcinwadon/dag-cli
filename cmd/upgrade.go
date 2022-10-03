package cmd

import (
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/node"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		return node.Upgrade(cfg)
	},
}
