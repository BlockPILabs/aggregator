package server

import (
	"encoding/json"
	"github.com/BlockPILabs/aggregator/config"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

func routeIndex(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

func routeConfig(ctx *fasthttp.RequestCtx) {
	data, _ := json.Marshal(config.Config)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(data)
}

func routeUpdateConfig(ctx *fasthttp.RequestCtx) {
	data, _ := json.Marshal(config.Config)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(data)
}

func NewManageServer() error {
	r := router.New()
	r.GET("/", routeIndex)
	r.GET("/config", routeConfig)
	r.POST("/config", routeUpdateConfig)

	addr := ":8012"
	logger.Info("Starting aggregator manage server", "addr", addr)
	err := fasthttp.ListenAndServe(addr, func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		if string(ctx.Method()) == "OPTIONS" {
			ctx.SetStatusCode(http.StatusNoContent)
			ctx.SetBodyString("ok")
			return
		}
		r.Handler(ctx)
	})
	if err != nil {
		notify.SendError("Error start manage server.", err.Error())
		return err
	}
	return nil
}
