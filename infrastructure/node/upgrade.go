package node

import (
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/github"
	"dag-cli/pkg/fetch"
)

func UpgradeL0(cfg config.Config) (err error) {
	version := cfg.Tessellation.Version

	githubClient := github.GetDefaultClient(cfg.Github)
	url := githubClient.GetFetchURL(version, cfg.Github.L0Filename)

	err = fetch.DownloadFile(cfg.L0.Path, url)
	if err != nil {
		return err
	}

	seedlistUrl := githubClient.GetFetchURL(version, cfg.Github.SeedlistFilename)

	err = fetch.DownloadFile(cfg.L0.SeedlistPath, seedlistUrl)
	if err != nil {
		return err
	}

	return nil
}

func UpgradeL1(cfg config.Config) (err error) {
	githubClient := github.GetDefaultClient(cfg.Github)
	url := githubClient.GetFetchURL(cfg.Tessellation.Version, cfg.Github.L1Filename)

	err = fetch.DownloadFile(cfg.L1.Path, url)
	if err != nil {
		return err
	}

	return nil
}

func Upgrade(cfg config.Config) (err error) {
	err = UpgradeL0(cfg)
	if err != nil {
		return err
	}

	err = UpgradeL1(cfg)
	if err != nil {
		return err
	}

	return nil
}