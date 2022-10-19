package lb

import (
	"dag-cli/domain/node"
)

type LoadBalancer interface {
	GetClusterInfo() ([]node.PeerInfo, error)
	GetRandomReadyPeer() (*node.PeerInfo, error)
	FindPeerByHostPort(host string, port int)  (*node.PeerInfo, error)
	FindPeerById(id string) (*node.PeerInfo, error)
}
