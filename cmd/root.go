package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
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
	return rootCmd.Execute()
}

func init() {

}
