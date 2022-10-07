package cmd

import (
	"dag-cli/infrastructure/config"
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
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()

		version := args[0]

		if semver.IsValid(version) {
			return node.Upgrade(cfg, version)
		} else {
			return errors.New("invalid version")
		}
	},
}
