package be

type Response struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

type BlockExplorer interface {
	GetLatestGlobalSnapshot() (*GlobalSnapshot, error)
	GetGlobalSnapshotByHash(hash string) (*GlobalSnapshot, error)
	GetGlobalSnapshotByOrdinal(ordinal int64) (*GlobalSnapshot, error)
}