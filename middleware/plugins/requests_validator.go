package plugins

import (
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/valyala/fasthttp"
)

type RequestValidatorMiddleware struct {
	nextMiddleware middleware.Middleware
	enabled        bool
}

func NewRequestValidatorMiddleware() *RequestValidatorMiddleware {
	return &RequestValidatorMiddleware{enabled: true}
}

func (m *RequestValidatorMiddleware) Name() string {
	return "RequestValidatorMiddleware"
}

func (m *RequestValidatorMiddleware) Enabled() bool {
	return m.enabled
}

func (m *RequestValidatorMiddleware) Next() middleware.Middleware {
	return m.nextMiddleware
}

func (m *RequestValidatorMiddleware) SetNext(middleware middleware.Middleware) {
	m.nextMiddleware = middleware
}

func (m *RequestValidatorMiddleware) OnRequest(session *rpc.Session) error {
	//if !session.InitOnce {
	//	err := session.Init()
	//	logger.Debug("recv new request", "sid", session.SId(), "method", session.RpcMethod())
	//	return err
	//}
	if session.Method == "OPTIONS" {
		return aggregator.ErrMustReturn
	}

	if session.Method != "POST" {
		return aggregator.ErrInvalidMethod
	}

	session.RpcMethod()

	return nil
}

func (m *RequestValidatorMiddleware) OnProcess(session *rpc.Session) error {
	if session.Method == "OPTIONS" {
		if ctx, ok := session.RequestCtx.(*fasthttp.RequestCtx); ok {
			ctx.Error("ok", fasthttp.StatusOK)
		}

		return aggregator.ErrMustReturn
	}
	return nil
}

func (m *RequestValidatorMiddleware) OnResponse(session *rpc.Session) error {
	return nil
}
