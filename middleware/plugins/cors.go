package plugins

import (
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/valyala/fasthttp"
)

type CorsMiddleware struct {
	nextMiddleware middleware.Middleware
	enabled        bool
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{enabled: true}
}

func (m *CorsMiddleware) Name() string {
	return "CorsMiddleware"
}

func (m *CorsMiddleware) Enabled() bool {
	return m.enabled
}

func (m *CorsMiddleware) Next() middleware.Middleware {
	return m.nextMiddleware
}

func (m *CorsMiddleware) SetNext(middleware middleware.Middleware) {
	m.nextMiddleware = middleware
}

func (m *CorsMiddleware) OnRequest(session *rpc.Session) error {
	return nil
}

func (m *CorsMiddleware) OnProcess(session *rpc.Session) error {
	return nil
}

func (m *CorsMiddleware) OnResponse(session *rpc.Session) error {
	if ctx, ok := session.RequestCtx.(*fasthttp.RequestCtx); ok {
		if session.Method == "OPTIONS" {
			ctx.Response.Reset()
			ctx.Response.Header.Set("Access-Control-Max-Age", "86400")
		}
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token, Authorization")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		ctx.Response.Header.Set("X-Relay-Node", session.NodeName)

		ctx.SetStatusCode(fasthttp.StatusOK)
	}
	return nil
}
