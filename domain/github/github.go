package github

type GitHub interface {
	GetFetchURL(version string, filename string) string
}