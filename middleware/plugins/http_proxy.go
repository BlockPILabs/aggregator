package plugins

import (
	"github.com/BlockPILabs/aggregator/client"
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/valyala/fasthttp"
)

type HttpProxyMiddleware struct {
	nextMiddleware middleware.Middleware
	enabled        bool
}

func NewHttpProxyMiddleware() *HttpProxyMiddleware {
	return &HttpProxyMiddleware{enabled: true}
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
		logger.Debug("relay rpc -> "+session.RpcMethod(), "sid", session.SId(), "tries", session.Tries)
		err := client.NewClient(session.Cfg.RequestTimeout, session.Cfg.Proxy).Relay(&ctx.Request, &ctx.Response)
		//if ctx, ok := session.RequestCtx.(*fasthttp.RequestCtx); ok {
		//	ctx.Response.Header.Set("Access-Control-Max-Age", "86400")
		//	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		//	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		//	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		//	ctx.Response.Header.Set("X-Relay-Node", session.NodeName)
		//}
		return err
	}

	return nil
}

func (m *HttpProxyMiddleware) OnResponse(session *rpc.Session) error {
	return nil
}
