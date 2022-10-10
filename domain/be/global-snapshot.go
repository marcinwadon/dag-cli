package be

type GlobalSnapshot struct {
	Hash             string   `json:"hash"`
	Ordinal          int64    `json:"ordinal"`
	Height           int64    `json:"height"`
	SubHeight        int64    `json:"subHeight"`
	LastSnapshotHash string   `json:"lastSnapshotHash"`
	Blocks           []string `json:"blocks"`
	Timestamp        string   `json:"timestamp"`
}
