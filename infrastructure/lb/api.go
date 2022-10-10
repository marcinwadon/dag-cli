package lb

import (
	"dag-cli/domain/lb"
	"dag-cli/domain/node"
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

func (lb *loadbalancer) GetClusterInfo() ([]node.PeerInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/cluster/info", lb.url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var clusterInfo []node.PeerInfo
	if err := json.NewDecoder(resp.Body).Decode(&clusterInfo); err != nil {
		return nil, err
	}

	return clusterInfo, nil
}

func (lb *loadbalancer) GetRandomReadyPeer() (*node.PeerInfo, error) {
	peers, err := lb.GetClusterInfo()
	if err != nil {
		return nil, err
	}

	var readyPeers []node.PeerInfo

	for _, peer := range peers {
		if node.ParseString(peer.State) == node.Ready {
			readyPeers = append(readyPeers, peer)
		}
	}

	if len(readyPeers) == 0 {
		return nil, errors.New("no ready peers found")
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(readyPeers)

	return &readyPeers[n], nil
}

func GetClient(url string) lb.LoadBalancer {
	return &loadbalancer{url}
}
