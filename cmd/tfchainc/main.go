package main

import (
	"math/big"

	"github.com/rivine/rivine/pkg/client"
	"github.com/rivine/rivine/types"
)

func main() {
	defaultClientConfig := client.DefaultConfig()
	defaultClientConfig.Name = "tfchain"
	defaultClientConfig.CurrencyCoinUnit = "TFT"
	defaultClientConfig.CurrencyUnits = types.CurrencyUnits{
		OneCoin: types.NewCurrency(new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)),
	}

	client.DefaultCLIClient(defaultClientConfig)
}
