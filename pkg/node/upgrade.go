package node

import (
	"dag-cli/pkg/fetch"
	"dag-cli/pkg/github"
	"github.com/spf13/viper"
)

func UpgradeL0() (err error) {
	version := viper.GetString("tessellation.version")
	filename := viper.GetString("github.l0-filename")
	path := viper.GetString("l0.path")

	githubClient := github.GetDefaultClient()
	url := githubClient.GetFetchURL(version, filename)

	err = fetch.DownloadFile(path, url)
	if err != nil {
		return err
	}

	return nil
}

func UpgradeL1() (err error) {
	version := viper.GetString("tessellation.version")
	filename := viper.GetString("github.l1-filename")
	path := viper.GetString("l1.path")

	githubClient := github.GetDefaultClient()
	url := githubClient.GetFetchURL(version, filename)

	err = fetch.DownloadFile(path, url)
	if err != nil {
		return err
	}

	return nil
}

func Upgrade() (err error) {
	err = UpgradeL0()
	if err != nil {
		return err
	}

	err = UpgradeL1()
	if err != nil {
		return err
	}

	return nil
}