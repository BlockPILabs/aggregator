package server

import (
	"fmt"
	"github.com/BlockPILabs/aggregator/config"
	"github.com/BlockPILabs/aggregator/jsonrpc"
	"github.com/BlockPILabs/aggregator/loadbalance"
	"github.com/BlockPILabs/aggregator/log"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/valyala/fasthttp"
	//proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
	"strings"
	"time"
)

var (
	logger = log.Module("server")
)

//func factory(hostAddr string) (*proxy.ReverseProxy, error) {
//	p := proxy.NewReverseProxy(hostAddr)
//	return p, nil
//}

func parseChainFromPath(path string) string {
	ss := strings.Split(path, "/")
	if len(ss) == 2 {
		return strings.Trim(ss[1], " ")
	}
	return ""
}

var requestHandler = func(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("error", "msg", err)
		}
	}()

	path := string(ctx.Path())
	chain := parseChainFromPath(path)
	if !config.HasChain(chain) {
		ctx.Error("Unsupported chain "+path, fasthttp.StatusNotFound)
		return
	}
	tries := 0
	for {
		node := loadbalance.NextNode(chain)
		if node == nil {
			ctx.Error("load balancer: Node not found", fasthttp.StatusNotFound)
			return
		}
		logger.Debug("proxy", "src", path, "dst", node.Endpoint, "tries", tries+1)

		ctx.Request.SetRequestURI(node.Endpoint)

		//proxyServer, err := pool.Get(node.Host())
		//if err != nil {
		//	ctx.Error("load balancer: Node not found", fasthttp.StatusInternalServerError)
		//	return
		//}
		//defer pool.Put(proxyServer)
		//proxyServer.ServeHTTP(ctx)
		client := fasthttp.Client{}
		err := client.DoTimeout(&ctx.Request, &ctx.Response, time.Second*time.Duration(config.Config.RequestTimeout))
		if err != nil {
			if config.Config.MaxRetries > 0 {
				tries++
				if tries >= config.Config.MaxRetries {
					ctx.Write(jsonrpc.Error(-3200, "Max retries exceeded"))
					return
				} else {
					continue
				}
			} else {
				ctx.Write(jsonrpc.Error(-3200, fmt.Sprintf("Error connect to node [%s]", node.Name)))
				return
			}
		}
		return
	}
}

func NewServer() error {
	var err error
	//pool, err = proxy.NewChanPool(100, 10000, factory)
	//if err != nil {
	//	notify.SendError("Error start aggregator server.", err.Error())
	//	return err
	//}
	addr := ":8011"
	logger.Info("Starting aggregator proxy server", "addr", addr)
	err = fasthttp.ListenAndServe(addr, requestHandler)
	if err != nil {
		notify.SendError("Error start aggregator server.", err.Error())
		return err
	}
	return nil
}
