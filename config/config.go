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
	"time"
)

var (
	logger           = log.Module("config")
	defaultConfigUrl = "https://raw.githubusercontent.com/BlockPILabs/chain-specs/dev/aggregator/default-config.json"
	_DB_CONFIG_KEY   = []byte("config")
	Config           = config{}
)

type config struct {
	RequestTimeout int64                    `json:"request_timeout"`
	MaxRetries     int                      `json:"max_retries"`
	Nodes          map[string][]*types.Node `json:"nodes"`
}

func LoadDefault() {
	//todo load from network
	retries := 0
	for {
		statusCode, data, err := (&fasthttp.Client{}).GetTimeout(nil, defaultConfigUrl, time.Second*5)
		if err == nil && statusCode == 200 {
			err = json.Unmarshal(data, &Config)
			if err == nil {
				logger.Info("Load default config success")
				return
			}
		}
		if err != nil || statusCode != 200 {
			retries++
			logger.Error("Load default config failed", "error", err.Error(), "retries", retries)

			if retries >= 5 {
				if err != nil {
					notify.SendError("Load default config failed", err.Error())
				}
				if statusCode > 0 {
					notify.SendError("Load default config failed", fmt.Sprintf("Status Code: %d", statusCode))
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
		logger.Error("Load config failed", "error", err.Error())
		notify.SendError("Load config failed", err.Error())
		return err
	}
	defer db.Close()

	data, err := db.Get(_DB_CONFIG_KEY, nil)
	if err != nil && !errors.Is(err, leveldbErrors.ErrNotFound) {
		logger.Error("Load config failed", "error", err.Error())
		notify.SendError("Load config failed", err.Error())
		return err
	}

	if data != nil {
		err = json.Unmarshal(data, &Config)
		if err != nil {
			logger.Error("Load config failed", "error", err.Error())
			notify.SendError("Load config failed", err.Error())
			return err
		}
	}

	return nil
}

func Save() error {
	db, err := leveldb.OpenFile("data/db", nil)
	if err != nil {
		logger.Error("Save config failed", "error", err.Error())
		notify.SendError("Save config failed", err.Error())
		return err
	}
	defer db.Close()

	data, err := json.Marshal(Config)
	if err != nil {
		logger.Error("Save config failed", "error", err.Error())
		notify.SendError("Save config failed", err.Error())
		return err
	}

	err = db.Put(_DB_CONFIG_KEY, data, nil)
	if err != nil {
		logger.Error("Save config failed", "error", err.Error())
		notify.SendError("Save config failed", err.Error())
		return err
	}

	return nil
}

func HasChain(chain string) bool {
	if len(chain) > 0 {
		if v, ok := Config.Nodes[chain]; ok {
			if len(v) > 0 {
				return true
			}
		}
	}
	return false
}
