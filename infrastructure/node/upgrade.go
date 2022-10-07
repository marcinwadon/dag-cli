package node

import (
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/github"
	"dag-cli/pkg/fetch"
)

func UpgradeL0(cfg config.Config, version string) (err error) {
	githubClient := github.GetDefaultClient(cfg.Github)
	url := githubClient.GetFetchURL(version, cfg.Github.L0Filename)

	err = fetch.DownloadFile(cfg.GetL0JarFilename(), url)
	if err != nil {
		return err
	}

	seedlistUrl := githubClient.GetFetchURL(version, cfg.Github.SeedlistFilename)

	err = fetch.DownloadFile(cfg.GetL0SeedlistFilename(), seedlistUrl)
	if err != nil {
		return err
	}

	_ = cfg.Persist()

	return nil
}

func UpgradeL1(cfg config.Config, version string) (err error) {
	githubClient := github.GetDefaultClient(cfg.Github)
	url := githubClient.GetFetchURL(version, cfg.Github.L1Filename)

	err = fetch.DownloadFile(cfg.GetL1JarFilename(), url)
	if err != nil {
		return err
	}

	return nil
}

func Upgrade(cfg config.Config, version string) (err error) {
	err = UpgradeL0(cfg, version)
	if err != nil {
		return err
	}

	err = UpgradeL1(cfg, version)
	if err != nil {
		return err
	}

	cfg.Tessellation.Version = version
	_ = cfg.Persist()

	return nil
}
