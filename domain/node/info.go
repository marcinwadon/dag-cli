package node

type NodeInfo struct {
	State          string `json:"state"`
	Session        int64  `json:"session"`
	ClusterSession int64  `json:"clusterSession"`
	Version        string `json:"version"`
	Host           string `json:"host"`
	PublicPort     int    `json:"publicPort"`
	P2PPort        int    `json:"p2pPort"`
	Id             string `json:"id"`
}
