package node

type NodeClient interface {
	GetNodeInfo() (*NodeInfo, error)
	Join(id string, host string, p2pPort int) error
	Leave() error
}
