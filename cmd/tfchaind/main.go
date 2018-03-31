package main

import (
	"fmt"
	"math/big"

	"github.com/rivine/rivine/modules"

	"github.com/rivine/rivine/pkg/daemon"
	"github.com/rivine/rivine/types"
)

var (
	testnet     = "testnet"
	standardnet = "standard"
)

func main() {
	defaultDaemonConfig := daemon.DefaultConfig()
	defaultDaemonConfig.BlockchainInfo.Name = "tfchain"
	// Default network name, testnet for now since real network is not live yet
	defaultDaemonConfig.NetworkName = standardnet
	defaultDaemonConfig.CreateNetworConfig = SetupNetworks

	daemon.SetupDefaultDaemon(defaultDaemonConfig)
}

// SetupNetworks injects the correct chain constants and genesis nodes based on the chosen network
func SetupNetworks(name string) (daemon.NetworkConfig, error) {
	switch name {
	case standardnet:
		return daemon.NetworkConfig{
			Constants:      getStandardnetGenesis(),
			BootstrapPeers: getStandardnetBootstrapPeers(),
		}, nil
	case testnet:
		return daemon.NetworkConfig{
			Constants:      getTestnetGenesis(),
			BootstrapPeers: getTestnetBootstrapPeers(),
		}, nil
	default:
		return daemon.NetworkConfig{}, fmt.Errorf("Netork name %q not recognized", name)
	}
}

// getStandardnetGenesis explicitly sets all the required constants for the genesis block of the standard (prod) net
func getStandardnetGenesis() types.ChainConstants {
	cfg := types.DefaultChainConstants()

	// 1 coin = 1 000 000 000 of the smalles possible units
	cfg.CurrencyUnits = types.CurrencyUnits{
		OneCoin: types.NewCurrency(new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)),
	}

	// 2 minute block time
	cfg.BlockFrequency = 120

	// Payouts take rougly 1 day to mature.
	cfg.MaturityDelay = 720

	// The genesis timestamp
	cfg.GenesisTimestamp = types.Timestamp(1522465200) // Human time (GMT): Saturday, 31 of March of 2018 3:00:00

	// 1000 block window for difficulty
	cfg.TargetWindow = 1e3

	cfg.MaxAdjustmentUp = big.NewRat(25, 10)
	cfg.MaxAdjustmentDown = big.NewRat(10, 25)

	cfg.FutureThreshold = 1 * 60 * 60        // 1 hour.
	cfg.ExtremeFutureThreshold = 2 * 60 * 60 // 2 hours.

	cfg.StakeModifierDelay = 2000

	// Blockstake can be used roughly 1 minute after receiving
	cfg.BlockStakeAging = uint64(1 << 6)

	// Receive 1 coins when you create a block
	cfg.BlockCreatorFee = cfg.CurrencyUnits.OneCoin.Mul64(1)

	// Use 0.1 coins as minimum transaction fee
	cfg.MinimumTransactionFee = cfg.CurrencyUnits.OneCoin.Div64(10)

	// distribute initial coins
	cfg.GenesisCoinDistribution = []types.CoinOutput{
		{
			// 695M TFT for Capacity availability until 01/04
			Value: cfg.CurrencyUnits.OneCoin.Mul64(695 * 1000 * 1000),
			// temporary pool
			UnlockHash: unlockHashFromHex("014ab1cf49a331bef9225a51a68623daf7e112fce0e81a91194a7f4fe7af1d9a793bc52d4676d0"),
		},
		{
			// 3K TFT for dev purposes and 695M TFT for Capacity availability until 01/04
			Value: cfg.CurrencyUnits.OneCoin.Mul64(3000),
			// @glendc
			UnlockHash: unlockHashFromHex("01ad4f73417476f8b8350298681dd0fa8640baa53a91915417b1dd8103d118b543c992e6fba1c4"),
		},
		{
			// 1K TFT for dev purposes
			Value: cfg.CurrencyUnits.OneCoin.Mul64(1000),
			// @robvanmieghem
			UnlockHash: unlockHashFromHex("017eb744c7a97443d7927725bfb2a004384e4230386ea0f693f9ce1c1161d771878a1f048c887b"),
		},
		{
			// 1K TFT for dev purposes
			Value: cfg.CurrencyUnits.OneCoin.Mul64(1000),
			// @robvanmieghem
			UnlockHash: unlockHashFromHex("01cc55df18eb3b86670deb6cfbb9b62b8463b62738426f0c14a7ae8926d6b556fbac3aab17f437"),
		},
		{
			// 1K TFT for dev purposes
			Value: cfg.CurrencyUnits.OneCoin.Mul64(1000),
			// @RubenMattan
			UnlockHash: unlockHashFromHex("01e2dc01a686fc0ca25612871a6515f2e3b804aa63244bf19449ecb3c9aaaf36f0cc6839b0af60"),
		},
		{
			// 1K TFT for dev purposes
			Value: cfg.CurrencyUnits.OneCoin.Mul64(1000),
			// @FastGeert
			UnlockHash: unlockHashFromHex("0166da8a4ab39a621637d9e7eb4e1fbaf95f905856af13af1268fc1a79c65b4f6686dec75c9d94"),
		},
		{
			// 1K TFT for dev purposes
			Value: cfg.CurrencyUnits.OneCoin.Mul64(1000),
			// @zaibon
			UnlockHash: unlockHashFromHex("013dfb1c49e8b9a73bc8b460d9ef20fc1f40e0d034742950f70d983c455342719d1e9f656d002b"),
		},
		{
			// 1K TFT for dev purposes
			Value: cfg.CurrencyUnits.OneCoin.Mul64(1000),
			// @leesmet
			UnlockHash: unlockHashFromHex("018a28615b277eb7e7a0e6921e85ad5b3ca378ac210b7f258b0b11ef313ea2ce98bd2e2510472d"),
		},
	}

	// allocate block stakes
	cfg.GenesisBlockStakeAllocation = []types.BlockStakeOutput{
		{
			// 400 BS, initially allocated to a selected group of people (ambassadors)
			Value: types.NewCurrency64(400),
			// Pool address for other people (ambassadors)
			UnlockHash: unlockHashFromHex("01746677df456546d93729066dd88514e2009930f3eebac3c93d43c88a108f8f9aa9e7c6f58893"),
		},
		{
			// 100 BS, one BS for each first-generation TFT node
			Value: types.NewCurrency64(100),
			// @glendc
			UnlockHash: unlockHashFromHex("01ad4f73417476f8b8350298681dd0fa8640baa53a91915417b1dd8103d118b543c992e6fba1c4"),
		},
	}

	return cfg
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

	// distribute initial coins
	cfg.GenesisCoinDistribution = []types.CoinOutput{
		{
			// Create 100M coins
			Value: cfg.CurrencyUnits.OneCoin.Mul64(100 * 1000 * 1000),
			// @leesmet
			UnlockHash: unlockHashFromHex("01fc8714235d549f890f35e52d745b9eeeee34926f96c4b9ef1689832f338d9349b453898f7e51"),
		},
	}

	// allocate block stakes
	cfg.GenesisBlockStakeAllocation = []types.BlockStakeOutput{
		{
			// Create 3K blockstakes
			Value: types.NewCurrency64(3000),
			// @leesmet
			UnlockHash: unlockHashFromHex("01fc8714235d549f890f35e52d745b9eeeee34926f96c4b9ef1689832f338d9349b453898f7e51"),
		},
	}

	return cfg
}

func unlockHashFromHex(hstr string) (uh types.UnlockHash) {
	err := uh.LoadString(hstr)
	if err != nil {
		panic(fmt.Sprintf("func unlockHashFromHex(%s) failed: %v", hstr, err))
	}
	return
}

// getStandardnetBootstrapPeers sets the standard bootstrap node addresses
func getStandardnetBootstrapPeers() []modules.NetAddress {
	return []modules.NetAddress{
		"bootstrap1.threefoldtoken.com:23112",
		"bootstrap2.threefoldtoken.com:23112",
		"bootstrap3.threefoldtoken.com:23112",
		"bootstrap4.threefoldtoken.com:23112",
	}
}

// getTestnetBootstrapPeers sets the testnet bootstrap node addresses
func getTestnetBootstrapPeers() []modules.NetAddress {
	return []modules.NetAddress{
		"bootstrap1.testnet.threefoldtoken.com:23112",
		"bootstrap2.testnet.threefoldtoken.com:23112",
		"bootstrap3.testnet.threefoldtoken.com:23112",
		"bootstrap4.testnet.threefoldtoken.com:23112",
	}
}
