package github

import (
	"dag-cli/domain/github"
	"dag-cli/infrastructure/config"
	"fmt"
)

type gh struct {
	organization string
	repository   string
}

func (gh *gh) GetFetchURL(version string, filename string) string {
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", gh.organization, gh.repository, version, filename)
}

func GetClient(organization string, repository string) github.GitHub {
	return &gh{
		organization: organization,
		repository:   repository,
	}
}

func GetDefaultClient(cfg config.GithubConfig) github.GitHub {
	return GetClient(cfg.Organization, cfg.Repository)
}
