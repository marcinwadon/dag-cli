package cmd

import (
	"dag-cli/infrastructure/config"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().Bool("advanced", false, "Advanced configuration")
	configureCmd.Flags().Bool("reset", false, "Reset existing configuration")
}

var configureCmd = &cobra.Command{
	Use: "configure",
	Short: "Configure your node",
	RunE: func(cmd *cobra.Command, args []string) error {
		advanced, _ := cmd.Flags().GetBool("advanced")
		resetExistingConfig, _ := cmd.Flags().GetBool("reset")

		ac := config.New(os.Stdin)

		var tmpConfig config.Config

		currentConfig, err := config.LoadConfig()

		fmt.Printf("Default values are in square brackets. Press enter to keep the default value.\n")

		env := ac.Ask("Environment", "mainnet")

		if err != nil || resetExistingConfig {
			if env == "mainnet" {
				tmpConfig = config.DefaultMainnet()
			} else if env == "testnet" {
				tmpConfig = config.DefaultTestnet()
			}
		} else {
			tmpConfig = currentConfig
		}

		tmpConfig.WorkingDir = ac.AskWorkingDir(tmpConfig.WorkingDir)
		tmpConfig.ExternalIP = ac.AskExternalIP(tmpConfig.ExternalIP)

		fmt.Printf("-- Key configuration -- \n")
		tmpConfig.Key = ac.AskKey(tmpConfig.Key)

		fmt.Printf("-- L0 Peer configuration -- \n")
		tmpConfig.L0 = ac.AskL0(tmpConfig.L0)
		fmt.Printf("-- L1 Peer configuration -- \n")
		tmpConfig.L1 = ac.AskL1(tmpConfig.L1)

		if advanced {
			fmt.Printf("-- Advanced configuration -- \n")
			tmpConfig.Github = ac.AskGithub(tmpConfig.Github)
		}

		tmpConfig.Show()

		save := ac.Ask("Would you like to save this configuration? (y|n)", "no")
		if strings.ToLower(save) == "y" || strings.ToLower(save) == "yes" {
			err = tmpConfig.Persist()
		}

		return nil
	},
}
