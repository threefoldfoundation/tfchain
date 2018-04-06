package main

import (
	"github.com/threefoldfoundation/tfchain/pkg/config"

	"github.com/rivine/rivine/pkg/client"
)

func main() {
	defaultClientConfig := client.DefaultConfig()
	defaultClientConfig.Name = config.ThreeFoldTokenChainName
	defaultClientConfig.CurrencyCoinUnit = config.ThreeFoldTokenUnit
	defaultClientConfig.CurrencyUnits = config.GetCurrencyUnits()
	defaultClientConfig.Version = config.Version // blockchain version
	defaultClientConfig.MinimumTransactionFee = config.GetStandardnetGenesis().MinimumTransactionFee

	client.DefaultCLIClient(defaultClientConfig)
}
