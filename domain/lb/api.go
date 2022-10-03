package lb

import (
	"dag-cli/domain/node"
)

type LoadBalancer interface {
	GetClusterInfo() ([]node.PeerInfo, error)
	GetRandomReadyPeer() (*node.PeerInfo, error)
}
