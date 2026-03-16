package loadbalance

import (
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/config"
	"github.com/BlockPILabs/aggregator/log"
	"sync"
	"time"
)

var (
	_selectors = map[string]*WrSelector{}
	_mutex     sync.Mutex
	logger     = log.Module("load-balance")
)

func SetNodes(chain string, nodes []aggregator.Node) {
	_mutex.Lock()
	defer _mutex.Unlock()

	selector := &WrSelector{}
	selector.SetNodes(nodes)
	_selectors[chain] = selector
}

func NextNode(chain string) *aggregator.Node {
	_mutex.Lock()
	defer _mutex.Unlock()

	selector := _selectors[chain]
	if selector != nil {
		return selector.NextNode()
	}

	return nil
}

func TimeoutNode(chain string, nodeName string, d time.Duration) {
	_mutex.Lock()
	defer _mutex.Unlock()
	selector := _selectors[chain]
	if selector != nil {
		selector.TimeoutNode(nodeName, d)
	}
}

func LoadFromConfig() {
	for chain, nodes := range config.Default().Nodes {
		logger.Info("New load balancer", "chain", chain, "nodes", len(nodes))
		SetNodes(chain, nodes)
	}
}
