package node

type PeerInfo struct {
	Id         string `json:"id"`
	Ip         string `json:"ip"`
	PublicPort int    `json:"publicPort"`
	P2PPort    int    `json:"p2pPort"`
	Session    string `json:"session"`
	State      string `json:"state"`
}

type Addr struct {
	Ip   string `json:"ip"`
	Port int    `json:"publicPort"`
}

func (p *PeerInfo) Addr() Addr {
	return Addr{
		p.Ip,
		p.PublicPort,
	}
}

type L0Peer struct {
	Id   string
	Host string
	Port int
}

func L0PeerFromPeerInfo(peer PeerInfo) L0Peer {
	return L0Peer{
		Id:   peer.Id,
		Host: peer.Ip,
		Port: peer.PublicPort,
	}
}