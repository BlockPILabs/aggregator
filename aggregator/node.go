package aggregator

import "net/url"

type Node struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Weight   int64  `json:"weight"`
	ReadOnly bool   `json:"read_only"`
	Disabled bool   `json:"disabled"`
}

func (node *Node) Host() string {
	_url, _ := url.Parse(node.Endpoint)
	return _url.Host + ":443"
}
