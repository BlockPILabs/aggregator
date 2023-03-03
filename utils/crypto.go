package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func DecodeTx(hexTx string) (*types.Transaction, error) {
	rawTx := common.FromHex(hexTx)
	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(rawTx)
	if err != nil {
		return nil, err
	}

	//txType := tx.Type()

	//signer := &types.HomesteadSigner{}

	//signer := types.NewEIP155Signer(tx.ChainId())
	//
	//addr, err := types.Sender(signer, tx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//hexAddr := addr.Hex()
	//fmt.Println(txType, hexAddr)

	return tx, nil
}
