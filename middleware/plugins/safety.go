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
	"github.com/BlockPILabs/aggregator/safety"
	"github.com/BlockPILabs/aggregator/utils"
	"github.com/valyala/fasthttp"
	"strings"
	"sync"
	"time"
)

var (
	phishingAddressMap          map[string]*phishingAddress
	authorityPhishingAddressMap = map[string]map[string]*phishingAddress{}
	mu                          = sync.Mutex{}
	lastUpdateAt                time.Time
)

type SafetyMiddleware struct {
	nextMiddleware middleware.Middleware
	enabled        bool
}

type phishingAddress struct {
	Address     string
	Description string
	Reporter    string
}

func NewSafetyMiddleware() *SafetyMiddleware {
	m := &SafetyMiddleware{enabled: true}
	m.updatePhishingDb()
	m.updateAuthorityPhishingDb()
	go func() {
		for {
			if time.Since(lastUpdateAt) > time.Second*time.Duration(config.Default().PhishingDbUpdateInterval) {
				m.updatePhishingDb()
				m.updateAuthorityPhishingDb()
			}
			time.Sleep(time.Second * 10)
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
		//logger.Debug("new tx", "method", session.RpcMethod())
		rpcMethod := session.RpcMethod()
		rpcMethod = strings.ToLower(rpcMethod[strings.Index(rpcMethod, "_"):])

		targetAddress := ""

		switch rpcMethod {
		case strings.ToLower("_sendRawTransaction"):
			rawTx, ok := params.([]interface{})[0].(string)
			if !ok {
				return nil
			}
			tx, err := utils.DecodeTx(rawTx)
			if err != nil {
				logger.Warn("Unable to decode tx")
				notify.SendNotice("Unable to decode tx")
			} else {
				targetAddress = tx.To().Hex()
				//phishing, pha := m.isPhishingAddress(receiver)
				//if phishing {
				//	notify.SendError("Transaction is denied", receiver, pha.Description)
				//	logger.Error("transaction is denied", "Receiver", receiver, "Reason", pha.Description)
				//	return aggregator.ErrDenyRequest
				//}
				//session.ChainId = tx.ChainId().Int64()
				//session.Tx = tx
			}
		case strings.ToLower("_call"):
			targetAddress = params.([]interface{})[0].(map[string]interface{})["to"].(string)
		case strings.ToLower("_sendTransaction"):
			targetAddress = params.([]interface{})[0].(map[string]interface{})["to"].(string)
		case strings.ToLower("_sendTransactionAsFeePayer"):
			//targetAddress = params.([]interface{})[0].(map[string]interface{})["to"].(string)
		}

		if len(targetAddress) != 0 {
			phishing, pha := m.isPhishingAddress(targetAddress)
			if phishing {
				reporter := ""
				if len(pha.Reporter) > 0 {
					reporter = "Reporter: " + pha.Reporter
				}
				notify.Send("Option denied - scam address", m.shortAddress(targetAddress), reporter)
				logger.Error("Option denied", "target", targetAddress, "Reason", pha.Description, "reporter", pha.Reporter)
				return aggregator.ErrDenyRequest
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

func (m *SafetyMiddleware) shortAddress(address string) string {
	length := len(address)
	if length > 10 {
		return address[0:6] + "..." + address[length-4:]
	}
	return address
}

func (m *SafetyMiddleware) updateAuthorityPhishingDb() {
	cfg := config.Clone()
	cli := client.NewClient(cfg.RequestTimeout, cfg.Proxy)
	for _, adb := range cfg.AuthorityDB {
		if !adb.Enable {
			logger.Warn("Authority phishing db not enable", "provider", adb.Name)
			continue
		}
		logger.Info("Updating authority phishing db", "provider", adb.Name)
		req := &fasthttp.Request{}
		resp := &fasthttp.Response{}
		req.Header.SetMethod(fasthttp.MethodGet)
		req.Header.Set("Accept-Encoding", "gzip,deflate,br")
		req.SetRequestURI(adb.Url)

		err := cli.Do(req, resp)
		if err != nil {
			log.Error("Phishing db update failed", "url", adb.Url, "err", err)
			continue
		}
		result := map[string]string{}
		body, _ := resp.BodyUncompressed()
		err = json.Unmarshal(body, &result)

		addrMap := map[string]*phishingAddress{}
		for addr, desc := range result {
			pha := &phishingAddress{
				Address:     strings.ToLower(addr),
				Description: desc,
				Reporter:    adb.Name,
			}

			addrMap[addr] = pha
		}
		logger.Info("Updated authority phishing db", "addresses", len(addrMap))
		authorityPhishingAddressMap[adb.Name] = addrMap

	}
}

func (m *SafetyMiddleware) updatePhishingDb() {
	cfg := config.Clone()
	if cfg.PhishingDb == nil || len(cfg.PhishingDb) == 0 {
		return
	}

	cli := client.NewClient(cfg.RequestTimeout, cfg.Proxy)
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}
	req.Header.Set("Accept-Encoding", "gzip,deflate,br")

	hasError := false

	addrMap := map[string]*phishingAddress{}
	for _, dbUrl := range cfg.PhishingDb {
		logger.Info("Updating phishing db", "url", dbUrl)

		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Error update phishing db", "err", err)
				}
			}()

			req.SetRequestURI(dbUrl)
			req.Header.SetMethod(fasthttp.MethodGet)
			err := cli.Do(req, resp)
			if err != nil {
				log.Error("Phishing db update failed", "url", dbUrl, "err", err)
				hasError = true
				return
			}
			result := map[string]string{}
			body, _ := resp.BodyUncompressed()
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Error("Phishing db update failed", "url", dbUrl, "err", err)
				return
			}

			for addr, desc := range result {
				pha := &phishingAddress{
					Address:     strings.ToLower(addr),
					Description: desc,
				}
				addrMap[addr] = pha
			}

		}()
	}

	mu.Lock()
	defer mu.Unlock()
	phishingAddressMap = addrMap

	count := len(phishingAddressMap)
	logger.Info("Updated phishing db", "addresses", count)

	if !hasError {
		lastUpdateAt = time.Now()
	}
}

func (m *SafetyMiddleware) isPhishingAddress(address string) (bool, *phishingAddress) {
	mu.Lock()
	defer mu.Unlock()
	address = strings.ToLower(address)

	var isPhishingAddress bool

	var descs = map[string]string{}
	var reporters []string

	pha, exist := phishingAddressMap[address]
	if exist {
		isPhishingAddress = true
		descs[pha.Description] = pha.Description
		if len(pha.Reporter) > 0 {
			reporters = append(reporters, pha.Reporter)
		}
	}

	for provider, phaMap := range authorityPhishingAddressMap {
		var hash string
		switch provider {
		case "goplus":
			hash = safety.RpcHubAddress(safety.GoPlusAddress(address))
		case "slowmist":
			hash = safety.RpcHubAddress(safety.SlowMistAddress(address))
		}

		if len(hash) > 0 {
			pha, exist = phaMap[hash]
			if exist {
				descs[pha.Description] = pha.Description
				if len(pha.Reporter) > 0 {
					reporters = append(reporters, pha.Reporter)
				}
				isPhishingAddress = true
			}
		}
	}

	if isPhishingAddress {
		var desc []string
		for k, _ := range descs {
			desc = append(desc, k)
		}

		return true, &phishingAddress{
			Address:     address,
			Description: strings.TrimSpace(strings.Join(desc, ", ")),
			Reporter:    strings.TrimSpace(strings.Join(reporters, ", ")),
		}
	}

	return false, nil
}
