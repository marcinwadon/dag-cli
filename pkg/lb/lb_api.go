package lb

type Loadbalancer interface {
	GetClusterInfo() (interface{}, error)
}

type loadbalancer struct {
	url string
}

func (lb *loadbalancer) GetClusterInfo() (interface{}, error) {
	panic("Not implemented")
}

func GetClient(url string) Loadbalancer {
	return &loadbalancer{url}
}
