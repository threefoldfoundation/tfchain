package internal

import (
	"github.com/threefoldtech/rivine/types"
	"github.com/threefoldfoundation/tfchain/pkg/config"
)

// GetFoundationPoolCondition returns the pool condition of the Foundation.
func GetFoundationPoolCondition(network string) types.UnlockConditionProxy {
	switch network {
	case config.NetworkNameStandard:
		return config.GetStandardDaemonNetworkConfig().FoundationPoolCondition
	case config.NetworkNameTest:
		return config.GetTestnetDaemonNetworkConfig().FoundationPoolCondition
	case config.NetworkNameDev:
		return config.GetDevnetDaemonNetworkConfig().FoundationPoolCondition
	default:
		panic("unknown network name: " + network)
	}
}
