package main

import (
	"math/big"

	"github.com/rivine/rivine/build"
	"github.com/rivine/rivine/pkg/client"
	"github.com/rivine/rivine/types"
)

func main() {
	defaultClientConfig := client.DefaultConfig()
	defaultClientConfig.Name = "tfchain"
	defaultClientConfig.CurrencyCoinUnit = "TFT"
	oneCoin := types.NewCurrency(new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil))
	defaultClientConfig.CurrencyUnits = types.CurrencyUnits{
		OneCoin: oneCoin,
	}
	defaultClientConfig.Version = build.NewVersion(1, 0, 1)
	defaultClientConfig.MinimumTransactionFee = oneCoin.Div64(10) // has to stay in sync with config used in tfchaind

	client.DefaultCLIClient(defaultClientConfig)
}
