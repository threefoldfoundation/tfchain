package config

import (
	"github.com/threefoldtech/rivine/types"
)

// DaemonNetworkConfig defines network-specific tfchain constants.
type DaemonNetworkConfig struct {
	FoundationPoolAddress types.UnlockHash
	ERC20FeePoolAddress   types.UnlockHash
}

// GetStandardDaemonNetworkConfig returns the standard network config for the daemon
func GetStandardDaemonNetworkConfig() DaemonNetworkConfig {
	return DaemonNetworkConfig{
		// TODO: define final address
		FoundationPoolAddress: unlockHashFromHex("017267221ef1947bb18506e390f1f9446b995acfb6d08d8e39508bb974d9830b8cb8fdca788e34"),
		// TODO: define final address
		ERC20FeePoolAddress: unlockHashFromHex("017267221ef1947bb18506e390f1f9446b995acfb6d08d8e39508bb974d9830b8cb8fdca788e34"),
	}
}

// GetTestnetDaemonNetworkConfig returns the testnet network config for the daemon
func GetTestnetDaemonNetworkConfig() DaemonNetworkConfig {
	return DaemonNetworkConfig{
		// TODO: define final address
		FoundationPoolAddress: unlockHashFromHex("016148ac9b17828e0933796eaca94418a376f2aa3fefa15685cea5fa462093f0150e09067f7512"),
		// TODO: define final address
		ERC20FeePoolAddress: unlockHashFromHex("016148ac9b17828e0933796eaca94418a376f2aa3fefa15685cea5fa462093f0150e09067f7512"),
	}
}

// GetDevnetDaemonNetworkConfig returns the devnet network config for the daemon
func GetDevnetDaemonNetworkConfig() DaemonNetworkConfig {
	return DaemonNetworkConfig{
		// belongs to wallet with mnemonic:
		// carbon boss inject cover mountain fetch fiber fit tornado cloth wing dinosaur proof joy intact fabric thumb rebel borrow poet chair network expire else
		FoundationPoolAddress: unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
		// belongs to wallet with mnemonic:
		// carbon boss inject cover mountain fetch fiber fit tornado cloth wing dinosaur proof joy intact fabric thumb rebel borrow poet chair network expire else
		ERC20FeePoolAddress: unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
	}
}
