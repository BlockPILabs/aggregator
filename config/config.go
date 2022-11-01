package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BlockPILabs/aggregator/aggregator"
	"github.com/BlockPILabs/aggregator/log"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/valyala/fasthttp"
	"sync"
	"time"
)

var (
	logger           = log.Module("config")
	defaultConfigUrl = "https://raw.githubusercontent.com/BlockPILabs/chain-specs/dev/aggregator/default-config.json"
	_Config          = &Config{Password: "blockpi", RequestTimeout: 30, MaxRetries: 3}

	locker = sync.Mutex{}
)

type Config struct {
	Password       string                       `json:"password,omitempty"`
	Proxy          string                       `json:"proxy,omitempty"`
	RequestTimeout int64                        `json:"request_timeout,omitempty"`
	MaxRetries     int                          `json:"max_retries,omitempty"`
	Nodes          map[string][]aggregator.Node `json:"nodes"`
}

func (c Config) HasChain(chain string) bool {
	if len(chain) > 0 {
		if v, ok := c.Nodes[chain]; ok {
			if len(v) > 0 {
				return true
			}
		}
	}
	return false
}

func Clone() Config {
	locker.Lock()
	defer locker.Unlock()

	cfg := *_Config
	cfg.Nodes = map[string][]aggregator.Node{}

	for key, nodes := range _Config.Nodes {
		cfg.Nodes[key] = []aggregator.Node{}
		for _, node := range nodes {
			cfg.Nodes[key] = append(cfg.Nodes[key], node)
		}
	}
	return cfg
}

func Default() *Config {
	locker.Lock()
	defer locker.Unlock()

	return _Config
}

func SetDefault(cfg *Config) {
	locker.Lock()
	defer locker.Unlock()

	_Config = cfg
}

func LoadDefault() {
	retries := 0
	for {
		statusCode, data, err := (&fasthttp.Client{}).GetTimeout(nil, defaultConfigUrl, time.Second*5)
		if err == nil && statusCode == 200 {
			err = json.Unmarshal(data, &_Config)
			if err == nil {
				logger.Info("Load default Config success")
				return
			}
		}
		if err != nil || statusCode != 200 {
			retries++
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			logger.Error("Load default Config failed", "statusCode", statusCode, "error", errStr, "retries", retries)

			if retries >= 5 {
				notify.SendError("Load default Config failed", fmt.Sprintf("Status Code: %d\nError: %s", statusCode, errStr))
				logger.Warn("Load default config failed, you may manually update the config by postman or via blockpi aggregator manager [https://aggregator.blockpi.io]")
				//if err != nil {
				//	notify.SendError("Load default Config failed", err.Error())
				//}
				//if statusCode > 0 {
				//	notify.SendError("Load default Config failed", fmt.Sprintf("Status Code: %d", statusCode))
				//}
				return
			} else {
				time.Sleep(time.Second * 3)
			}
		}
	}

	//Config.RequestTimeout = 30
	//Config.Nodes = map[string][]*types.Node{}
	//Config.Nodes["klaytn"] = []*types.Node{
	//	{
	//		Name: "blockpi",
	//		//Endpoint: "http://168.119.1.188:9645/",
	//		Endpoint: "https://public-rpc.blockpi.me/klaytn",
	//		Weight:   100,
	//		ReadOnly: false,
	//		Disabled: false,
	//	},
	//}
}

func Load() error {
	LoadDefault()

	db, err := leveldb.OpenFile("data/db", nil)
	if err != nil {
		logger.Error("Load Config failed", "error", err.Error())
		notify.SendError("Load Config failed", err.Error())
		return err
	}
	defer db.Close()

	data, err := db.Get(aggregator.KeyDbConfig, nil)
	if err != nil && !errors.Is(err, leveldbErrors.ErrNotFound) {
		logger.Error("Load Config failed", "error", err.Error())
		notify.SendError("Load Config failed", err.Error())
		return err
	}

	if data != nil {
		err = json.Unmarshal(data, &_Config)
		if err != nil {
			logger.Error("Load Config failed", "error", err.Error())
			notify.SendError("Load Config failed", err.Error())
			return err
		}
	}

	return nil
}

func Save() error {
	db, err := leveldb.OpenFile("data/db", nil)
	if err != nil {
		logger.Error("Save Config failed", "error", err.Error())
		notify.SendError("Save Config failed", err.Error())
		return err
	}
	defer db.Close()

	data, err := json.Marshal(Default())
	if err != nil {
		logger.Error("Save Config failed", "error", err.Error())
		notify.SendError("Save Config failed", err.Error())
		return err
	}

	err = db.Put(aggregator.KeyDbConfig, data, nil)
	if err != nil {
		logger.Error("Save Config failed", "error", err.Error())
		notify.SendError("Save Config failed", err.Error())
		return err
	}

	return nil
}

func Chains() []string {
	var chains []string
	for key, _ := range Default().Nodes {
		chains = append(chains, key)
	}
	return chains
}
