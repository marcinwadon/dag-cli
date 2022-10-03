package node

import (
	"bytes"
	"dag-cli/domain/node"
	"dag-cli/pkg/api"
	"encoding/json"
	"fmt"
	"net/http"
)

type nd struct {
	host string
	port int
}

func (nd *nd) getEndpoint(endpoint string) string {
	return fmt.Sprintf("http://%s:%d/%s", nd.host, nd.port, endpoint)
}

func (nd *nd) GetNodeInfo() (*node.NodeInfo, error) {
	resp, err := http.Get(nd.getEndpoint("node/info"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var nodeInfo node.NodeInfo
	if err := json.NewDecoder(resp.Body).Decode(&nodeInfo); err != nil {
		return nil, err
	}

	return &nodeInfo, nil
}

func (nd nd) Join(id string, host string, p2pPort int) error {
	body, _ := json.Marshal(map[string]any{
		"id":      id,
		"ip":      host,
		"p2pPort": p2pPort,
	})
	requestBody := bytes.NewBuffer(body)

	resp, err := http.Post(nd.getEndpoint("cluster/join"), api.ApplicationJson, requestBody)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("an error occured when joining %s:%d, status code: %d", host, p2pPort, resp.StatusCode)
	}

	return nil
}

func GetClient(host string, port int) node.NodeClient {
	return &nd{host: host, port: port}
}
