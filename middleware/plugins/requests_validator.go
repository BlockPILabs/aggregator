package plugins

import (
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/rpc"
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
	return nil
}

func (m *RequestValidatorMiddleware) OnResponse(session *rpc.Session) error {
	return nil
}
