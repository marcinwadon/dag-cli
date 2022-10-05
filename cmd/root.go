package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "dag",
		Short:   "DAG Command Line Utility",
		Version: "0.3.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("status")

			return nil
		},
	}
	Verbose bool
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("dag")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dag")
	viper.AddConfigPath("$HOME/.dag")
	viper.AddConfigPath(".dag")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
