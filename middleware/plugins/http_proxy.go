package plugins

import (
	"github.com/BlockPILabs/aggregator/client"
	"github.com/BlockPILabs/aggregator/log"
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/valyala/fasthttp"
	"sync"
	"time"
)

type HttpProxyMiddleware struct {
	nextMiddleware  middleware.Middleware
	enabled         bool
	client          *client.Client
	clientCreatedAt time.Time
	clientRenew     time.Duration
	mu              sync.Mutex
}

func NewHttpProxyMiddleware() *HttpProxyMiddleware {
	return &HttpProxyMiddleware{
		enabled:     true,
		clientRenew: time.Second * 60,
		mu:          sync.Mutex{},
	}
}

func (m *HttpProxyMiddleware) Name() string {
	return "HttpProxyMiddleware"
}

func (m *HttpProxyMiddleware) Enabled() bool {
	return m.enabled
}

func (m *HttpProxyMiddleware) Next() middleware.Middleware {
	return m.nextMiddleware
}

func (m *HttpProxyMiddleware) SetNext(middleware middleware.Middleware) {
	m.nextMiddleware = middleware
}

func (m *HttpProxyMiddleware) OnRequest(session *rpc.Session) error {
	return nil
}

func (m *HttpProxyMiddleware) OnProcess(session *rpc.Session) error {
	if ctx, ok := session.RequestCtx.(*fasthttp.RequestCtx); ok {
		logger.Debug("relay rpc -> "+session.RpcMethod(), "sid", session.SId(), "node", session.NodeName, "isTx", session.IsWriteRpcMethod, "tries", session.Tries)
		err := m.GetClient(session).Do(&ctx.Request, &ctx.Response)
		//if ctx, ok := session.RequestCtx.(*fasthttp.RequestCtx); ok {
		//	ctx.Response.Header.Set("Access-Control-Max-Age", "86400")
		//	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		//	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		//	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		//	ctx.Response.Header.Set("X-Do-Node", session.NodeName)
		//}

		shouldDisableEndpoint := false
		if err != nil {
			log.Error(err.Error(), "node", session.NodeName)
			shouldDisableEndpoint = true
		}

		statusCode := ctx.Response.StatusCode()
		if statusCode/100 != 2 {
			log.Error("error status code", "code", statusCode, "node", session.NodeName)
			shouldDisableEndpoint = true
		}

		if shouldDisableEndpoint {
			//todo disable endpoint
		}

		return err
	}

	return nil
}

func (m *HttpProxyMiddleware) OnResponse(session *rpc.Session) error {
	return nil
}

func (m *HttpProxyMiddleware) GetClient(session *rpc.Session) *client.Client {
	m.mu.Lock()
	defer m.mu.Unlock()

	if time.Since(m.clientCreatedAt) <= m.clientRenew {
		if m.client != nil {
			return m.client
		}
	}

	//log.Debug("renew proxy http client")
	m.client = client.NewClient(session.Cfg.RequestTimeout, session.Cfg.Proxy)
	m.clientCreatedAt = time.Now()

	return m.client
}
