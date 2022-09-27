package main

import (
	"dag-cli/cmd"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigName("dag")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dag")
	viper.AddConfigPath("$HOME/.dag")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln(err)
	}

	cmd.Execute()
}
