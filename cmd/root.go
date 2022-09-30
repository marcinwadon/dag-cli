package cmd

import (
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/lb"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dag",
		Short: "DAG Command Line Utility",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("status")
			cfg, _ := config.LoadConfig()

			fmt.Printf("\nEnter Layer to configure (yes/no) [yes]: ")
			fmt.Printf("\nEnter L0 load balancer url [https://l0-lb-testnet.constellationnetwork.io]: ")
			var lay string
			fmt.Scanln(&lay)

			fmt.Println(lay)

			return nil

			lbClient := lb.GetClient(cfg.L0.LoadBalancer)
			peer, err := lbClient.GetRandomReadyPeer()
			if err != nil {
				return err
			}
			fmt.Println(peer.Id)

			return nil
		},
	}
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("dag")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dag")
	viper.AddConfigPath("$HOME/.dag")
	viper.AddConfigPath(".dag")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
