package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BlockPILabs/aggregator/log"
	"github.com/BlockPILabs/aggregator/notify"
	"github.com/BlockPILabs/aggregator/types"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"sync"
	"time"
)

var (
	logger           = log.Module("Config")
	defaultConfigUrl = "https://raw.githubusercontent.com/BlockPILabs/chain-specs/dev/aggregator/default-config.json"
	_DB_CONFIG_KEY   = []byte("Config")
	_Config          = Config{}

	locker = sync.Mutex{}
)

type Config struct {
	Password       string                  `json:"password,omitempty"`
	Proxy          string                  `json:"proxy,omitempty"`
	RequestTimeout int64                   `json:"request_timeout,omitempty"`
	MaxRetries     int                     `json:"max_retries,omitempty"`
	Nodes          map[string][]types.Node `json:"nodes"`
}

func Clone() Config {
	locker.Lock()
	defer locker.Unlock()

	cfg := _Config
	cfg.Nodes = map[string][]types.Node{}

	for key, nodes := range _Config.Nodes {
		cfg.Nodes[key] = []types.Node{}
		for _, node := range nodes {
			cfg.Nodes[key] = append(cfg.Nodes[key], node)
		}
	}
	return cfg
}

func Default() Config {
	locker.Lock()
	defer locker.Unlock()

	return _Config
}

func SetDefault(cfg Config) {
	locker.Lock()
	defer locker.Unlock()

	_Config = cfg
}

func NewConfig() Config {
	return Config{}
}

func LoadDefault() {
	retries := 0
	for {
		statusCode, data, err := (&fasthttp.Client{Dial: fasthttpproxy.FasthttpSocksDialer("socks5://107.148.129.228:8982")}).GetTimeout(nil, defaultConfigUrl, time.Second*5)
		if err == nil && statusCode == 200 {
			err = json.Unmarshal(data, &_Config)
			if err == nil {
				logger.Info("Load default Config success")
				return
			}
		}
		if err != nil {
			retries++
			logger.Error("Load default Config failed", "error", err.Error(), "retries", retries)

			if retries >= 5 {
				if err != nil {
					notify.SendError("Load default Config failed", err.Error())
				}
				if statusCode > 0 {
					notify.SendError("Load default Config failed", fmt.Sprintf("Status Code: %d", statusCode))
				}
				return
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

	data, err := db.Get(_DB_CONFIG_KEY, nil)
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

	err = db.Put(_DB_CONFIG_KEY, data, nil)
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

func HasChain(chain string) bool {
	if len(chain) > 0 {
		if v, ok := Default().Nodes[chain]; ok {
			if len(v) > 0 {
				return true
			}
		}
	}
	return false
}
