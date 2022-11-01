package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func DecodeTx(hexTx string) *types.Message {
	rawTx := common.FromHex(hexTx)
	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(rawTx)
	if err != nil {
		return nil
	}
	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), tx.GasPrice())
	if err != nil {
		return nil
	}
	return &msg
}
