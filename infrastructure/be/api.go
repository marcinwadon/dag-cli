package be

import (
	"dag-cli/domain/be"
	"encoding/json"
	"fmt"
	"net/http"
)

type b struct {
	url string
}

func (b *b) getEndpoint(endpoint string) string {
	return fmt.Sprintf("https://%s/%s", b.url, endpoint)
}

func (b b) GetLatestGlobalSnapshot() (*be.GlobalSnapshot, error) {
	return b.GetGlobalSnapshotByHash("latest")
}

func (b b) GetGlobalSnapshotByHash(hash string) (*be.GlobalSnapshot, error) {
	resp, err := http.Get(b.getEndpoint(fmt.Sprintf("global-snapshots/%s", hash)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var globalSnapshot be.GlobalSnapshot
	var response = be.Response{Data: &globalSnapshot}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &globalSnapshot, nil
}

func (b b) GetGlobalSnapshotByOrdinal(ordinal int64) (*be.GlobalSnapshot, error) {
	return b.GetGlobalSnapshotByHash(fmt.Sprintf("%d", ordinal))
}

func GetClient(url string) be.BlockExplorer {
	return &b{url: url}
}