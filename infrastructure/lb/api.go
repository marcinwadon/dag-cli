package lb

import (
	"dag-cli/domain/lb"
	node2 "dag-cli/domain/node"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type loadbalancer struct {
	url string
}

func (lb *loadbalancer) GetClusterInfo() ([]node2.PeerInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/cluster/info", lb.url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var clusterInfo []node2.PeerInfo
	if err := json.NewDecoder(resp.Body).Decode(&clusterInfo); err != nil {
		return nil, err
	}

	return clusterInfo, nil
}

func (lb *loadbalancer) GetRandomReadyPeer() (*node2.PeerInfo, error) {
	peers, err := lb.GetClusterInfo()
	if err != nil {
		return nil, err
	}

	var readyPeers []node2.PeerInfo

	for _, peer := range peers {
		if node2.ParseString(peer.State) == node2.Ready {
			readyPeers = append(readyPeers, peer)
		}
	}

	if len(readyPeers) == 0 {
		return nil, errors.New("no ready peers found")
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(readyPeers)

	return &peers[n], nil
}

func GetClient(url string) lb.LoadBalancer {
	return &loadbalancer{url}
}
