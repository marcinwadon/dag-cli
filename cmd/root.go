package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use: "dag",
		Short: "DAG Command Line Utility",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("status")
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
