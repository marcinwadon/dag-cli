package node

import "fmt"

type NodeState string

const (
	Initial         NodeState = "Initial"
	ReadyToJoin     NodeState = "ReadyToJoin"
	LoadingGenesis  NodeState = "LoadingGenesis"
	GenesisReady    NodeState = "GenesisReady"
	StartingSession NodeState = "StartingSession"
	SessionStarted  NodeState = "SessionStarted"
	Ready           NodeState = "Ready"
	Leaving         NodeState = "Leaving"
	Offline         NodeState = "Offline"
	Unknown         NodeState = "Unknown"
	Observing       NodeState = "Observing"
	Undefined       NodeState = "Undefined" // Internal status when we could not obtain status for the node
)

var ValidStatuses = [...]NodeState{Initial, ReadyToJoin, LoadingGenesis, GenesisReady, StartingSession, SessionStarted, Ready, Leaving, Offline, Unknown, Observing, Undefined}

func StateFromString(in string) NodeState {
	for _, v := range ValidStatuses {
		if in == fmt.Sprint(v) {
			return v
		}
	}

	return Unknown
}

type Addr struct {
	Ip   string `json:"ip"`
	Port int    `json:"publicPort"`
}

type PeerInfo struct {
	Id         string `json:"id"`
	Ip         string `json:"ip"`
	PublicPort int    `json:"publicPort"`
	P2PPort    int    `json:"p2pPort"`
	Session    string `json:"session"`
	State      string `json:"state"`

	cardinalState NodeState
}
type ClusterInfo struct {
	Id    string
	Peers *Peers
}

type Peers []PeerInfo

func (p *PeerInfo) Addr() Addr {
	return Addr{
		p.Ip,
		p.PublicPort,
	}
}
