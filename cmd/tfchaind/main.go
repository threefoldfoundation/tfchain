package main

import (
	"fmt"

	"github.com/threefoldfoundation/tfchain/pkg/config"

	"github.com/rivine/rivine/pkg/daemon"
)

func main() {
	defaultDaemonConfig := daemon.DefaultConfig()
	defaultDaemonConfig.BlockchainInfo = config.GetBlockchainInfo()
	defaultDaemonConfig.CreateNetworkConfig = SetupNetworksAndTypes

	daemon.SetupDefaultDaemon(defaultDaemonConfig)
}

// SetupNetworksAndTypes injects the correct chain constants and genesis nodes based on the chosen network,
// it also ensures that features added during the lifetime of the blockchain,
// only get activated on a certain block height, giving everyone sufficient time to upgrade should such features be introduced.
func SetupNetworksAndTypes(name string) (daemon.NetworkConfig, error) {
	// return the network configuration, based on the network name,
	// which includes the genesis block as well as the bootstrap peers
	switch name {
	case config.NetworkNameStandard:
		// Forbid the usage of MultiSignatureCondition (and thus the multisig feature),
		// until the blockchain reached a height of 42000 blocks.
		RegisteredBlockHeightLimitedMultiSignatureCondition()

		// return the standard genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      config.GetStandardnetGenesis(),
			BootstrapPeers: config.GetStandardnetBootstrapPeers(),
		}, nil

	case config.NetworkNameTest:
		// return the testnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      config.GetTestnetGenesis(),
			BootstrapPeers: config.GetTestnetBootstrapPeers(),
		}, nil

	case config.NetworkNameDev:
		// return the devnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      config.GetDevnetGenesis(),
			BootstrapPeers: nil,
		}, nil

	default:
		// network isn't recognised
		return daemon.NetworkConfig{}, fmt.Errorf("Netork name %q not recognized", name)
	}
}
