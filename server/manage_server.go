package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/BlockPILabs/aggregator/config"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

var basicAuthPrefix = []byte("Basic ")

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
	logger.Info("Starting management server", "addr", addr)
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

		auth := ctx.Request.Header.Peek("Authorization")
		if bytes.HasPrefix(auth, basicAuthPrefix) {
			payload, err := base64.StdEncoding.DecodeString(string(auth[len(basicAuthPrefix):]))
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && bytes.Equal(pair[0], []byte("blockpi")) && bytes.Equal(pair[1], []byte(config.Config.Password)) {
					r.Handler(ctx)
					return
				}
			}
		}
		ctx.Error("Unauthorized", fasthttp.StatusUnauthorized)
	})
	if err != nil {
		notify.SendError("Error start manage server.", err.Error())
		return err
	}
	return nil
}
