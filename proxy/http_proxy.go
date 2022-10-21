package proxy

import "github.com/valyala/fasthttp"

type HttpProxy struct {
	clients []*fasthttp.Client
}

func (p *HttpProxy) getClient() *fasthttp.Client {
	if p.clients == nil {
		return nil

	}
	return p.clients[0]
}
