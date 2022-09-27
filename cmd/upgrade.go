package cmd

import (
	"dag-cli/pkg/node"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use: "upgrade",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		node.Start()
		//node.Upgrade()
	},
}
