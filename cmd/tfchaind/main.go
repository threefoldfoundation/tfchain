package main

import (
	"fmt"

	"github.com/threefoldfoundation/tfchain/pkg/config"

	"github.com/rivine/rivine/pkg/daemon"
	"github.com/rivine/rivine/types"
)

var (
	devnet      = "devnet"
	testnet     = "testnet"
	standardnet = "standard"
)

func main() {
	defaultDaemonConfig := daemon.DefaultConfig()
	defaultDaemonConfig.BlockchainInfo = config.GetBlockchainInfo()
	// Default network name, testnet for now since real network is not live yet
	defaultDaemonConfig.NetworkName = standardnet
	defaultDaemonConfig.CreateNetworConfig = SetupNetworks

	daemon.SetupDefaultDaemon(defaultDaemonConfig)
}

// SetupNetworks injects the correct chain constants and genesis nodes based on the chosen network
func SetupNetworks(name string) (daemon.NetworkConfig, error) {
	// configure any tfchain-specific logic related to how
	// unlock conditions/fulfillments are to be interpreted/used.
	ConfigureUnlockConditions(name)

	// return the network configuration, based on the network name,
	// which includes the genesis block as well as the bootstrap peers
	switch name {
	case standardnet:
		return daemon.NetworkConfig{
			Constants:      config.GetStandardnetGenesis(),
			BootstrapPeers: config.GetStandardnetBootstrapPeers(),
		}, nil
	case testnet:
		return daemon.NetworkConfig{
			Constants:      config.GetTestnetGenesis(),
			BootstrapPeers: config.GetTestnetBootstrapPeers(),
		}, nil
	case devnet:
		return daemon.NetworkConfig{
			Constants:      config.GetDevnetGenesis(),
			BootstrapPeers: nil,
		}, nil

	default:
		return daemon.NetworkConfig{}, fmt.Errorf("Netork name %q not recognized", name)
	}
}

// ConfigureUnlockConditions configures the unlock conditions,
// as to define the unlock conditions and fulfillments.
// For the most parts the configuration follows the standard Rivine unlock condition configuration,
// but there are a couple of differences:
//
// * if the networkName equals "standard",
//   the MultiSignatureCondition is only considered standard since block height 42000,
//   giving peers on the standard network more or less 7 days, since 1.0.6 was released
func ConfigureUnlockConditions(networkName string) {
	if networkName != standardnet {
		return // only apply the blockheight-based restrictions for standard net
	}
	// Overwrite the Rivine-standard MultiSignatureCondition with our wrapped version,
	// as to ensure that this condition is only introduced in the
	// standard tfchain since blockheight 42000
	overwriteMultiSignatureConditionType()
	// MultiSignatureFulfillment doesn't need to be overriden,
	// as it is impossible to use it as long as there isn't an unspend output
	// which uses the MultiSignatureCondition
}

func overwriteMultiSignatureConditionType() {
	types.RegisterUnlockConditionType(types.ConditionTypeMultiSignature,
		func() types.MarshalableUnlockCondition { return new(MultiSignatureCondition) })
}

// MultiSignatureCondition wraps around the Rivine-standard MultiSignatureCondition type,
// as to ensure that in the standard network of tfchain, it can only be used since blockheight 42000
type MultiSignatureCondition struct {
	*types.MultiSignatureCondition
}

const (
	// MinimumBlockHeightForMultiSignatureConditions defines the blockheight
	// since when MultiSignatureConditions are considered standard on the tfchain standard network.
	MinimumBlockHeightForMultiSignatureConditions = types.BlockHeight(42000)
)

// IsStandardCondition implements UnlockCondition.IsStandardCondition,
// wrapping around the internal MultiSignatureCondition's IsStandardCondition check,
// adding a pre-check of the blockheight
func (msc MultiSignatureCondition) IsStandardCondition(ctx types.StandardCheckContext) error {
	if ctx.BlockHeight < MinimumBlockHeightForMultiSignatureConditions {
		return fmt.Errorf(
			"multisignature conditions are only allowed since blockheight %d",
			MinimumBlockHeightForMultiSignatureConditions)
	}
	return msc.MultiSignatureCondition.IsStandardCondition(ctx)
}
