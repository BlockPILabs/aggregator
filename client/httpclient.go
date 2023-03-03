package client

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"strings"
	"time"
)

type Client struct {
	client     fasthttp.Client
	timeout    int64
	maxRetries int64
	proxy      string
}

func DefaultClient() *Client {
	return NewClient(30, "")
}

func NewClient(timeout int64, proxy string) *Client {
	cli := &Client{
		client: fasthttp.Client{
			MaxConnsPerHost: 65000,
			//Dial: func(addr string) (net.Conn, error) {
			//	return nil, nil
			//},
		},
		timeout: timeout,
		proxy:   proxy,
	}
	if proxy != "" {
		if strings.HasPrefix(proxy, "socks5://") {
			cli.client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
		} else {
			cli.client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
		}
	}
	return cli
}

func (cli *Client) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	return cli.client.DoTimeout(req, resp, time.Second*time.Duration(cli.timeout))
}
