package plugins

import (
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/BlockPILabs/aggregator/utils"
	"strings"
)

type SafetyMiddleware struct {
	nextMiddleware middleware.Middleware
	enabled        bool
}

func NewSafetyMiddleware() *SafetyMiddleware {
	return &SafetyMiddleware{enabled: true}
}

func (m *SafetyMiddleware) Name() string {
	return "SafetyMiddleware"
}

func (m *SafetyMiddleware) Enabled() bool {
	return m.enabled
}

func (m *SafetyMiddleware) Next() middleware.Middleware {
	return m.nextMiddleware
}

func (m *SafetyMiddleware) SetNext(middleware middleware.Middleware) {
	m.nextMiddleware = middleware
}

func (m *SafetyMiddleware) OnRequest(session *rpc.Session) error {
	if session.IsWriteRpcMethod {
		params := session.RpcParams()
		logger.Debug("new tx", "method", session.RpcMethod(), "params", params)
		rawTx := params.([]interface{})[0].(string)
		msg := utils.DecodeTx(rawTx)
		if msg != nil {
			receiver := msg.To().Hex()
			if m.isUnsafeReceiver(receiver) {
				notify.SendError("Transaction is denied", "Receiver "+receiver)
				logger.Error("transaction is denied", "Receiver", receiver)
				return aggregator.ErrDenyRequest
			}
		}
	}
	return nil
}

func (m *SafetyMiddleware) OnProcess(session *rpc.Session) error {
	return nil
}

func (m *SafetyMiddleware) OnResponse(session *rpc.Session) error {
	return nil
}

func (m *SafetyMiddleware) isUnsafeReceiver(address string) bool {
	if strings.ToLower(address) == strings.ToLower("0x68349009458626e35da0EeA9cB583b3C828bB815") {
		return true
	}
	return false
}
