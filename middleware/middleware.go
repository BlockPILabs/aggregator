package middleware

import (
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/log"
	"github.com/BlockPILabs/aggregator/rpc"
)

var (
	middlewareChain []Middleware
	logger          = log.Module("middleware")
)

type Middleware interface {
	Name() string

	Next() Middleware
	SetNext(next Middleware)

	Enabled() bool

	OnRequest(session *rpc.Session) error
	OnProcess(session *rpc.Session) error
	OnResponse(session *rpc.Session) error
}

func Append(middlewares ...Middleware) {
	for _, mw := range middlewares {
		if middlewareChain != nil && len(middlewareChain) > 0 {
			middlewareChain[len(middlewareChain)-1].SetNext(mw)
		}

		middlewareChain = append(middlewareChain, mw)
	}
}

func OnRequest(session *rpc.Session) error {
	for _, mw := range middlewareChain {
		err := mw.OnRequest(session)
		if err != nil {
			if err == aggregator.ErrMustReturn {
				return nil
			}
			logger.Error("an error occurred", "sid", session.SId(), "middleware", mw.Name(), "error", err)
			return err
		}
	}
	return nil
}

func OnProcess(session *rpc.Session) error {
	for _, mw := range middlewareChain {
		err := mw.OnProcess(session)
		if err != nil {
			if err == aggregator.ErrMustReturn {
				return nil
			}
			logger.Error("an error occurred", "sid", session.SId(), "middleware", mw.Name(), "error", err)
			return err
		}
	}
	return nil
}

func OnResponse(session *rpc.Session) error {
	for _, mw := range middlewareChain {
		err := mw.OnResponse(session)
		if err != nil {
			if err == aggregator.ErrMustReturn {
				return nil
			}
			logger.Error("an error occurred", "sid", session.SId(), "middleware", mw.Name(), "error", err)
			return err
		}
	}
	return nil
}
