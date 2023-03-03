package plugins

import (
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/loadbalance"
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/valyala/fasthttp"
)

type LoadBalanceMiddleware struct {
	nextMiddleware middleware.Middleware
	enabled        bool
}

func NewLoadBalanceMiddleware() *LoadBalanceMiddleware {
	return &LoadBalanceMiddleware{enabled: true}
}

func (m *LoadBalanceMiddleware) Name() string {
	return "LoadBalanceMiddleware"
}

func (m *LoadBalanceMiddleware) Enabled() bool {
	return m.enabled
}

func (m *LoadBalanceMiddleware) Next() middleware.Middleware {
	return m.nextMiddleware
}

func (m *LoadBalanceMiddleware) SetNext(middleware middleware.Middleware) {
	m.nextMiddleware = middleware
}

func (m *LoadBalanceMiddleware) OnRequest(session *rpc.Session) error {
	node := loadbalance.NextNode(session.Chain)
	if node == nil {
		return aggregator.ErrServerError
	}
	session.NodeName = node.Name
	//logger.Debug("load balance", "sid", session.SId(), "node", node.Name)
	if ctx, ok := session.RequestCtx.(*fasthttp.RequestCtx); ok {
		ctx.Request.SetRequestURI(node.Endpoint)
	}
	return nil
}

func (m *LoadBalanceMiddleware) OnProcess(session *rpc.Session) error {
	return nil
}

func (m *LoadBalanceMiddleware) OnResponse(session *rpc.Session) error {
	return nil
}
