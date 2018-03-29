package main

import (
	"fmt"
	"math/big"

	"github.com/rivine/rivine/modules"

	"github.com/rivine/rivine/pkg/daemon"
	"github.com/rivine/rivine/types"
)

var (
	testnet = "testnet"
)

func main() {

	// setGenesis()
	// setBootstrapPeers()

	defaultDaemonConfig := daemon.DefaultConfig()
	defaultDaemonConfig.BlockchainInfo.Name = "tfchain"
	// Default network name, testnet for now since real network is not live yet
	defaultDaemonConfig.NetworkName = testnet
	defaultDaemonConfig.CreateNetworConfig = SetupNetworks

	daemon.SetupDefaultDaemon(defaultDaemonConfig)
}

// SetupNetworks injects the correct chain constants and genesis nodes based on the chosen network
func SetupNetworks(name string) (daemon.NetworkConfig, error) {
	if name == testnet {
		return daemon.NetworkConfig{
			Constants:      getTestnetGenesis(),
			BootstrapPeers: getTestnetBootstrapPeers(),
		}, nil
	}

	return daemon.NetworkConfig{}, fmt.Errorf("Netork name \"%v\" not recognized", name)
}

// getTestnetGenesis explicitly sets all the required constants for the genesis block of the testnet
func getTestnetGenesis() types.ChainConstants {

	cfg := types.DefaultChainConstants()

	// 1 coin = 1 000 000 000 of the smalles possible units
	cfg.CurrencyUnits = types.CurrencyUnits{
		OneCoin: types.NewCurrency(new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)),
	}

	// 2 minute block time
	cfg.BlockFrequency = 120

	// Payouts take rougly 1 day to mature.
	cfg.MaturityDelay = 720

	// The genesis timestamp is set to February 21st, 2018
	cfg.GenesisTimestamp = types.Timestamp(1519200000) // February 21st, 2018 @ 8:00am UTC.

	// 1000 block window for difficulty
	cfg.TargetWindow = 1e3

	cfg.MaxAdjustmentUp = big.NewRat(25, 10)
	cfg.MaxAdjustmentDown = big.NewRat(10, 25)

	cfg.FutureThreshold = 1 * 60 * 60        // 1 hour.
	cfg.ExtremeFutureThreshold = 2 * 60 * 60 // 2 hours.

	cfg.StakeModifierDelay = 2000

	// Blockstake can be used roughly 1 minute after receiving
	cfg.BlockStakeAging = uint64(1 << 6)

	// Receive 10 coins when you create a block
	cfg.BlockCreatorFee = cfg.CurrencyUnits.OneCoin.Mul64(10)

	// Use 0.1 coins as minimum transaction fee
	cfg.MinimumTransactionFee = cfg.CurrencyUnits.OneCoin.Div64(10)

	// Create 3K blockstakes
	bso := types.BlockStakeOutput{
		Value:      types.NewCurrency64(3000),
		UnlockHash: types.UnlockHash{},
	}

	// Create 100M coins
	co := types.CoinOutput{
		Value: cfg.CurrencyUnits.OneCoin.Mul64(100 * 1000 * 1000),
	}

	bso.UnlockHash.LoadString("01fc8714235d549f890f35e52d745b9eeeee34926f96c4b9ef1689832f338d9349b72d12744e14")
	cfg.GenesisBlockStakeAllocation = []types.BlockStakeOutput{}
	cfg.GenesisBlockStakeAllocation = append(cfg.GenesisBlockStakeAllocation, bso)
	co.UnlockHash.LoadString("01fc8714235d549f890f35e52d745b9eeeee34926f96c4b9ef1689832f338d9349b72d12744e14")
	cfg.GenesisCoinDistribution = []types.CoinOutput{}
	cfg.GenesisCoinDistribution = append(cfg.GenesisCoinDistribution, co)

	return cfg
}

// getTestnetBootstrapPeers sets the bootstrap node addresses
func getTestnetBootstrapPeers() []modules.NetAddress {
	return []modules.NetAddress{
		"bootstrap1.testnet.threefoldtoken.com:23112",
		"bootstrap2.testnet.threefoldtoken.com:23112",
		"bootstrap3.testnet.threefoldtoken.com:23112",
		"bootstrap4.testnet.threefoldtoken.com:23112",
	}
}
