package main

import (
	"github.com/rivine/rivine/modules"
	"github.com/rivine/rivine/pkg/client"
	"github.com/threefoldfoundation/tfchain/pkg/config"
)

func main() {
	bchainInfo := config.GetBlockchainInfo()
	client.DefaultCLIClient("", bchainInfo.Name, func(icfg *client.Config) client.Config {
		cfg := daemonOrDefaultConfig(icfg)
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

// daemonOrDefaultConfig uses a default config
// if a config was not returned originating from the used daemon's constants.
func daemonOrDefaultConfig(icfg *client.Config) client.Config {
	if icfg != nil {
		return *icfg
	}

	bchainInfo := config.GetBlockchainInfo()
	chainConstants := config.GetStandardnetGenesis()
	daemonConstants := modules.NewDaemonConstants(bchainInfo, chainConstants)
	return client.ConfigFromDaemonConstants(daemonConstants)
}
