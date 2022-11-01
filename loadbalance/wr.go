package loadbalance

import (
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/notify"
	"math/rand"
	"sync"
)

// WrSelector weighted-random selector
type WrSelector struct {
	nodes     []aggregator.Node
	sumWeight int64

	mutex sync.Mutex
}

func (s *WrSelector) SetNodes(nodes []aggregator.Node) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var nodesSelected []aggregator.Node
	var sumWeight int64 = 0
	for _, node := range nodes {
		if !node.Disabled {
			if node.Weight > 0 && len(node.Endpoint) > 0 {
				sumWeight += node.Weight
				nodesSelected = append(nodesSelected, node)
			} else {
				notify.SendError("load balance: node is not selected", node.Name, node.Endpoint)
			}
		} else {
			logger.Warn("Node is disabled", "node", node.Name, "endpoint", node.Endpoint)
		}
	}
	s.nodes = nodesSelected
	s.sumWeight = sumWeight
}

func (s *WrSelector) NextNode() *aggregator.Node {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sumWeight > 0 {
		w := rand.Int63n(s.sumWeight)
		var weight int64 = 0
		for _, node := range s.nodes {
			//if !node.Disabled {
			weight += node.Weight
			if weight >= w {
				return &node
			}
			//}
		}
	}
	return nil
}
