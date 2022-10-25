package server

import (
	"github.com/BlockPILabs/aggregator/config"
	"github.com/BlockPILabs/aggregator/loadbalance"
	"github.com/BlockPILabs/aggregator/log"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"

	//proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
	"strings"
	"time"
)

var (
	logger = log.Module("server")
)

func parseChainFromPath(path string) string {
	ss := strings.Split(path, "/")
	if len(ss) == 2 {
		return strings.Trim(ss[1], " ")
	}
	return ""
}

func doRpcRelay(request *fasthttp.Request, response *fasthttp.Response) *rpc.JsonRpcResponse {
	rpcReq := rpc.MustUnmarshalJsonRpcRequest(request.Body())
	if rpcReq == nil {
		return rpc.ErrorInvalidRequest(0, "Invalid request")
	}

	path := string(request.URI().Path())
	chain := parseChainFromPath(path)
	if !config.HasChain(chain) {
		return rpc.ErrorInvalidRequest(rpcReq.Id, "Unsupported chain path "+path)
	}

	tries := 0
	for {
		node := loadbalance.NextNode(chain)
		if node == nil {
			return rpc.ErrorServerError(rpcReq.Id, "Node not found")
		}
		logger.Debug("proxy "+rpcReq.Method, "src", path, "dst", node.Endpoint, "tries", tries+1)

		request.SetRequestURI(node.Endpoint)

		client := fasthttp.Client{}
		if config.Config.Proxy != "" {
			if strings.HasPrefix(config.Config.Proxy, "socks5://") {
				client.Dial = fasthttpproxy.FasthttpSocksDialer(config.Config.Proxy)
			} else {
				client.Dial = fasthttpproxy.FasthttpHTTPDialer(config.Config.Proxy)
			}
		}

		err := client.DoTimeout(request, response, time.Second*time.Duration(config.Config.RequestTimeout))
		if err != nil {
			if config.Config.MaxRetries > 0 {
				tries++
				if tries >= config.Config.MaxRetries {
					return rpc.ErrorServerError(rpcReq.Id, "Max retries exceeded")
				} else {
					continue
				}
			} else {
				return rpc.ErrorServerError(rpcReq.Id, "Connect to node failed")
			}
		}

		return nil
	}
}

var requestHandler = func(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("error", "msg", err)
		}
	}()

	err := doRpcRelay(&ctx.Request, &ctx.Response)
	if err != nil {
		ctx.Error(string(err.Marshal()), fasthttp.StatusOK)
	}
}

func NewServer() error {
	var err error
	addr := ":8011"
	logger.Info("Starting proxy server", "addr", addr)

	for _, chain := range config.Chains() {
		logger.Info("Registered RPC", "endpoint", "http://localhost:8011/"+chain)
	}

	err = fasthttp.ListenAndServe(addr, requestHandler)
	if err != nil {
		notify.SendError("Error start aggregator server.", err.Error())
		return err
	}
	return nil
}
