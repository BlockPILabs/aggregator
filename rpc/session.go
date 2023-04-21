package rpc

import (
	"fmt"
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/config"
	"github.com/valyala/fasthttp"
	"strings"
	"sync"
	"sync/atomic"
)

var _id int64 = 0

type Session struct {
	once       sync.Once
	sId        any
	RequestCtx any
	Method     string
	Path       string
	Chain      string
	Request    *JsonRpcRequest
	RawRequest []byte
	Cfg        config.Config

	Tries            int
	NodeName         string
	IsWriteRpcMethod bool

	//Tx *types.Transaction
}

func (s *Session) Init() error {
	var err error

	s.once.Do(func() {
		s.sId = atomic.AddInt64(&_id, 1)
		s.Cfg = config.Clone()
		if ctx, ok := s.RequestCtx.(*fasthttp.RequestCtx); ok {
			s.Method = string(ctx.Method())
			s.Path = string(ctx.URI().Path())
			s.RawRequest = ctx.Request.Body()

			ss := strings.Split(s.Path, "/")
			if len(ss) != 2 {
				err = aggregator.ErrInvalidRequest
				return
			}
			s.Chain = strings.Trim(ss[1], " ")
			s.Request = MustUnmarshalJsonRpcRequest(ctx.Request.Body())
		}

		if !s.Cfg.HasChain(s.Chain) {
			err = aggregator.ErrInvalidChain
			return
		}
	})

	return err
}

func (s *Session) SId() string {
	return fmt.Sprintf("s-%016d", s.sId)
}

func (s *Session) Id() any {
	var id any = 1
	if s.Request != nil {
		id = s.Request.Id
	}
	return id
}

func (s *Session) IsMaxRetriesExceeded() bool {
	return s.Tries >= s.Cfg.MaxRetries
}

func (s *Session) RpcMethod() string {
	if s.Request != nil {
		return s.Request.Method
	}
	return ""
}

func (s *Session) RpcParams() interface{} {
	if s.Request != nil {
		return s.Request.Params
	}
	return nil
}

func (s *Session) NewJsonRpcError(err error) *JsonRpcResponse {
	id := s.Id()
	if agErr, ok := err.(*aggregator.Error); ok {
		return Error(id, agErr.Code, agErr.Message)
	}
	return ErrorInvalidRequest(id, err.Error())
}
