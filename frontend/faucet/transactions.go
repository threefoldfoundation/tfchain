package main

import (
	"encoding/json"
	"log"

	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/types"
)

func dripCoins(address types.UnlockHash, amount types.Currency) (types.TransactionID, error) {
	data, err := json.Marshal(api.WalletCoinsPOST{
		CoinOutputs: []types.CoinOutput{
			{
				Value:     amount,
				Condition: types.NewCondition(types.NewUnlockHashCondition(address)),
			},
		},
	})
	if err != nil {
		return types.TransactionID{}, err
	}

	log.Println("[DEBUG] Dripping", amount.String(), "coins to address", address.String())

	var resp api.WalletCoinsPOSTResp
	err = httpClient.PostWithResponse("/wallet/coins", string(data), &resp)
	if err != nil {
		log.Println("[ERROR] /wallet/coins - request body:", string(data))
	}
	return resp.TransactionID, err
}
