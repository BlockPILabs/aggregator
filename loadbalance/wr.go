package loadbalance

import (
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/notify"
	"math/rand"
	"sync"
	"time"
)

// WrSelector weighted-random selector
type WrSelector struct {
	nodes     []aggregator.Node
	sumWeight int64
	timeouts  map[string]time.Time

	mutex sync.Mutex
}

func (s *WrSelector) SetNodes(nodes []aggregator.Node) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.timeouts == nil {
		s.timeouts = map[string]time.Time{}
	}

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

func (s *WrSelector) TimeoutNode(name string, d time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.timeouts == nil {
		s.timeouts = map[string]time.Time{}
	}
	s.timeouts[name] = time.Now().Add(d)
}

func (s *WrSelector) NextNode() *aggregator.Node {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sumWeight > 0 {
		w := rand.Int63n(s.sumWeight)
		var weight int64 = 0
		now := time.Now()
		for i := 0; i < len(s.nodes); i++ {
			name := s.nodes[i].Name
			if until, ok := s.timeouts[name]; ok {
				if now.Before(until) {
					continue
				}
				delete(s.timeouts, name)
			}
			weight += s.nodes[i].Weight
			if weight >= w {
				return &s.nodes[i]
			}
		}
	}
	return nil
}
