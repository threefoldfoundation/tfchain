package main

import (
	"github.com/rivine/rivine/pkg/client"
	"github.com/threefoldfoundation/tfchain/pkg/config"
)

func main() {
	bchainInfo := config.GetBlockchainInfo()
	client.DefaultCLIClient("", bchainInfo.Name, func(icfg *client.Config) client.Config {
		if icfg == nil {
			constants := config.GetStandardnetGenesis()
			return client.Config{
				ChainName:    bchainInfo.Name,
				NetworkName:  config.NetworkNameStandard,
				ChainVersion: bchainInfo.ChainVersion,

				CurrencyUnits:             config.GetCurrencyUnits(),
				MinimumTransactionFee:     constants.MinimumTransactionFee,
				DefaultTransactionVersion: constants.DefaultTransactionVersion,

				BlockFrequencyInSeconds: int64(constants.BlockFrequency),
				GenesisBlockTimestamp:   constants.GenesisBlock().Timestamp,
			}
		}

		cfg := *icfg
		switch cfg.NetworkName {
		case config.NetworkNameStandard:
			// overwrite standard network genesis block stamp,
			// as the genesis block is way earlier than the actual first block,
			// due to the hard reset at the bumpy/rough start
			cfg.GenesisBlockTimestamp = 1524168391 // timestamp of (standard) block #1
		case config.NetworkNameTest:
			// seems like testnet timestamp wasn't updated last time it was reset
			cfg.GenesisBlockTimestamp = 1522792547 // timestamp of (testnet) block #1
		}
		return cfg
	})
}
