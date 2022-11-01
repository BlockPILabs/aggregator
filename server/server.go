package server

import (
	"github.com/BlockPILabs/aggregator/config"
	"github.com/BlockPILabs/aggregator/log"
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/valyala/fasthttp"
)

var (
	logger = log.Module("server")
)

var requestHandler = func(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("error", "msg", err)
		}
	}()

	var err error

	session := &rpc.Session{RequestCtx: ctx}
	err = session.Init()
	if err != nil {
		ctx.Error(string(session.NewJsonRpcError(err).Marshal()), fasthttp.StatusOK)
		return
	}
	for {
		session.Tries++
		err = middleware.OnRequest(session)
		if err != nil {
			if session.IsMaxRetriesExceeded() {
				ctx.Error(string(session.NewJsonRpcError(err).Marshal()), fasthttp.StatusOK)
				return
			}
			continue
		}

		err = middleware.OnProcess(session)
		if err != nil {
			if session.IsMaxRetriesExceeded() {
				ctx.Error(string(session.NewJsonRpcError(err).Marshal()), fasthttp.StatusOK)
				return
			}
			continue
		}

		err = middleware.OnResponse(session)
		if err != nil {
			if session.IsMaxRetriesExceeded() {
				ctx.Error(string(session.NewJsonRpcError(err).Marshal()), fasthttp.StatusOK)
				return
			}
			continue
		}
		return
	}
}

func NewServer() error {
	var err error
	addr := ":8011"
	logger.Info("Starting proxy server", "addr", addr)

	for _, chain := range config.Chains() {
		logger.Info("Registered RPC", "endpoint", "http://localhost:8011/"+chain)
	}

	s := &fasthttp.Server{
		Handler:            fasthttp.CompressHandlerLevel(requestHandler, 6),
		MaxRequestBodySize: fasthttp.DefaultMaxRequestBodySize * 10,
	}

	err = s.ListenAndServe(addr)
	if err != nil {
		notify.SendError("Error start aggregator server.", err.Error())
		return err
	}
	return nil
}
