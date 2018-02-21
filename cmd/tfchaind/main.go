package main

import (
	"math/big"

	"github.com/rivine/rivine/modules"

	"github.com/rivine/rivine/rivined"
	"github.com/rivine/rivine/types"
)

func main() {

	setGenesis()
	setBootstrapPeers()

	// Set daemon name for help messages
	rivined.DaemonName = "tfchain"

	defaultDaemonConfig := rivined.DefaultConfig()
	rivined.SetupDefaultDaemon(defaultDaemonConfig)
}

// setGenesis explicitly sets all the required constants in the types package, mainly for the genesis block
func setGenesis() {

	// 10 minute block time
	types.BlockFrequency = 600

	// Payouts take rougly 1 day to mature.
	types.MaturityDelay = 144

	// The genesis timestamp is set to February 21st, 2018
	types.GenesisTimestamp = types.Timestamp(1519200000) // February 21st, 2018 @ 8:00am UTC.

	// 1000 block window for difficulty
	types.TargetWindow = 1e3

	types.MaxAdjustmentUp = big.NewRat(25, 10)
	types.MaxAdjustmentDown = big.NewRat(10, 25)

	types.FutureThreshold = 3 * 60 * 60        // 3 hours.
	types.ExtremeFutureThreshold = 5 * 60 * 60 // 5 hours.

	types.StakeModifierDelay = 2000

	// (2^16s < 1 day < 2^17s)
	types.BlockStakeAging = uint64(1 << 17)

	// Receive 10 coins when you create a block
	types.BlockCreatorFee = types.OneCoin.Mul64(10)

	// Create 1M blockstakes
	bso := types.BlockStakeOutput{
		Value:      types.NewCurrency64(1000000),
		UnlockHash: types.UnlockHash{},
	}

	// Create 100M coins
	co := types.CoinOutput{
		Value: types.OneCoin.Mul64(100 * 1000 * 1000),
	}

	bso.UnlockHash.LoadString("02b1a92f2cb1b2daec2f650717452367273335263136fae0201ddedbbcfe67648572b069c754")
	types.GenesisBlockStakeAllocation = append(types.GenesisBlockStakeAllocation, bso)
	co.UnlockHash.LoadString("02b1a92f2cb1b2daec2f650717452367273335263136fae0201ddedbbcfe67648572b069c754")
	types.GenesisCoinDistribution = append(types.GenesisCoinDistribution, co)

	types.CalculateGenesis()
}

// setBootstrapPeers sets the bootstrap node addresses
func setBootstrapPeers() {
	modules.BootstrapPeers = []modules.NetAddress{}
}
