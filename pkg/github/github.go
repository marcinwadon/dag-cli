package github

import (
	"fmt"
	"github.com/spf13/viper"
)

type GitHub interface {
	GetFetchURL(version string, filename string) string
}

type github struct {
	organization string
	repository   string
}

func (gh *github) GetFetchURL(version string, filename string) string {
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", gh.organization, gh.repository, version, filename)
}

func GetClient(organization string, repository string) GitHub {
	return &github{
		organization: organization,
		repository:   repository,
	}
}

func GetDefaultClient() GitHub {
	organization := viper.GetString("github.organization")
	repository := viper.GetString("github.repository")

	return GetClient(organization, repository)
}