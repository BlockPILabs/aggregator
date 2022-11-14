package plugins

import (
	"encoding/json"
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/client"
	"github.com/BlockPILabs/aggregator/config"
	"github.com/BlockPILabs/aggregator/log"
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/BlockPILabs/aggregator/rpc"
	"github.com/BlockPILabs/aggregator/utils"
	"github.com/valyala/fasthttp"
	"strings"
	"sync"
	"time"
)

var (
	phishingAddressMap = map[string]*phishingAddress{}
	mu                 = sync.Mutex{}
	nextUpdateAt       = time.Now()
)

type SafetyMiddleware struct {
	nextMiddleware middleware.Middleware
	enabled        bool
}

type phishingAddress struct {
	Address     string
	Description string
}

func NewSafetyMiddleware() *SafetyMiddleware {
	m := &SafetyMiddleware{enabled: true}
	m.updatePhishingDb()
	go func() {
		for {
			if nextUpdateAt.Sub(time.Now()) <= 0 {
				m.updatePhishingDb()
			}
			time.Sleep(time.Second)
		}
	}()

	return m
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
		logger.Debug("new tx", "method", session.RpcMethod())

		if strings.HasSuffix(strings.ToLower(session.RpcMethod()), strings.ToLower("_sendRawTransaction")) {
			rawTx, ok := params.([]interface{})[0].(string)
			if !ok {
				return nil
			}
			msg := utils.DecodeTx(rawTx)
			if msg != nil {
				receiver := msg.To().Hex()

				phishing, pha := m.isPhishingAddress(receiver)
				if phishing {
					notify.SendError("Transaction is denied", receiver, pha.Description)
					logger.Error("transaction is denied", "Receiver", receiver, "Reason", pha.Description)
					return aggregator.ErrDenyRequest
				}
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

func (m *SafetyMiddleware) updatePhishingDb() {
	cfg := config.Clone()
	if cfg.PhishingDb == nil || len(cfg.PhishingDb) == 0 {
		return
	}
	nextUpdateAt = time.Now().Add(time.Second * 10)

	cli := client.NewClient(cfg.RequestTimeout, cfg.Proxy)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	var phishingAddresses []*phishingAddress

	for _, dbUrl := range cfg.PhishingDb {
		logger.Info("Updating phishing db", "url", dbUrl)

		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Error update phishing db", "err", err)
				}
			}()

			resp.Reset()
			req.SetRequestURI(dbUrl)
			req.Header.SetMethod(fasthttp.MethodGet)
			err := cli.Do(req, resp)
			if err != nil {
				log.Error("Phishing db update failed", "url", dbUrl, "err", err)
				return
			}
			result := map[string]interface{}{}
			err = json.Unmarshal(resp.Body(), &result)
			if err != nil {
				log.Error("Phishing db update failed", "url", dbUrl, "err", err)
				return
			}

			if result["success"].(bool) {
				for addr, val := range result["result"].(map[string]interface{}) {
					pha := &phishingAddress{
						Address: strings.ToLower(addr),
					}
					description := (val.([]interface{})[0]).(map[string]interface{})["description"]
					if description != nil {
						pha.Description = description.(string)
					}
					phishingAddresses = append(phishingAddresses, pha)
				}
			}
		}()
	}

	mu.Lock()
	defer mu.Unlock()
	if phishingAddresses != nil && len(phishingAddresses) > 0 {
		for addr, _ := range phishingAddressMap {
			found := false
			for _, pha := range phishingAddresses {
				if addr == pha.Address {
					found = true
					break
				}
			}
			if !found {
				delete(phishingAddressMap, addr)
			}
		}
		for _, pha := range phishingAddresses {
			phishingAddressMap[pha.Address] = pha
		}
	}

	count := len(phishingAddressMap)
	logger.Info("Updated phishing db", "addresses", count)

	nextUpdateAt = time.Now().Add(time.Second * time.Duration(cfg.PhishingDbUpdateInterval))
}

func (m *SafetyMiddleware) isPhishingAddress(address string) (exist bool, pha *phishingAddress) {
	mu.Lock()
	defer mu.Unlock()
	pha, exist = phishingAddressMap[strings.ToLower(address)]
	return
}
