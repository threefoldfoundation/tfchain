package main

import (
	"fmt"

	"github.com/rivine/rivine/modules"
	"github.com/rivine/rivine/pkg/client"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldfoundation/tfchain/pkg/types"
)

func main() {
	bchainInfo := config.GetBlockchainInfo()
	client.DefaultCLIClient("", bchainInfo.Name, func(icfg *client.Config) client.Config {
		cfg := daemonOrDefaultConfig(icfg)
		switch cfg.NetworkName {
		case config.NetworkNameStandard:
			// Register the transaction controllers for all transaction versions
			// supported on the standard network
			types.RegisterTransactionTypesForStandardNetwork()
			// Forbid the usage of MultiSignatureCondition (and thus the multisig feature),
			// until the blockchain reached a height of 42000 blocks.
			types.RegisterBlockHeightLimitedMultiSignatureCondition(42000)

			// overwrite standard network genesis block stamp,
			// as the genesis block is way earlier than the actual first block,
			// due to the hard reset at the bumpy/rough start
			cfg.GenesisBlockTimestamp = 1524168391 // timestamp of (standard) block #1

		case config.NetworkNameTest:
			// Register the transaction controllers for all transaction versions
			// supported on the test network
			types.RegisterTransactionTypesForTestNetwork()
			// Use our custom MultiSignatureCondition, just for testing purposes
			types.RegisterBlockHeightLimitedMultiSignatureCondition(0)

			// seems like testnet timestamp wasn't updated last time it was reset
			cfg.GenesisBlockTimestamp = 1522792547 // timestamp of (testnet) block #1

		case config.NetworkNameDev:
			// Register the transaction controllers for all transaction versions
			// supported on the dev network
			types.RegisterTransactionTypesForDevNetwork()
			// Use our custom MultiSignatureCondition, just for testing purposes
			types.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		default:
			panic(fmt.Sprintf("Netork name %q not recognized", cfg.NetworkName))
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
