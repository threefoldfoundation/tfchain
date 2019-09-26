package config

import (
	"github.com/threefoldtech/rivine/types"
)

// DaemonNetworkConfig defines network-specific tfchain constants.
type DaemonNetworkConfig struct {
	GenesisMintingCondition types.UnlockConditionProxy // used for minting
	FoundationPoolAddress   types.UnlockHash           // used for 3Bot
	ERC20FeePoolAddress     types.UnlockHash           // used for ERC20
}

// GetStandardDaemonNetworkConfig returns the standard network config for the daemon
func GetStandardDaemonNetworkConfig() DaemonNetworkConfig {
	return DaemonNetworkConfig{
		GenesisMintingCondition: types.NewCondition(types.NewMultiSignatureCondition(types.UnlockHashSlice{
			unlockHashFromHex("01434535fd01243c02c277cd58d71423163767a575a8ae44e15807bf545e4a8456a5c4afabad51"),
			unlockHashFromHex("01334cf68f312026ff9df84fc023558db8624bedd717adcc9edc6900488cf6df54ac8e3d1c89a8"),
			unlockHashFromHex("0149a5496fea27315b7db6251e5dfda23bc9d4bf677c5a5c2d70f1382c44357197d8453d9dfa32"),
		}, 2)),
		// TODO: define final address
		FoundationPoolAddress: unlockHashFromHex("017267221ef1947bb18506e390f1f9446b995acfb6d08d8e39508bb974d9830b8cb8fdca788e34"),
		// TODO: define final address
		ERC20FeePoolAddress: unlockHashFromHex("017267221ef1947bb18506e390f1f9446b995acfb6d08d8e39508bb974d9830b8cb8fdca788e34"),
	}
}

// GetTestnetDaemonNetworkConfig returns the testnet network config for the daemon
func GetTestnetDaemonNetworkConfig() DaemonNetworkConfig {
	return DaemonNetworkConfig{
		GenesisMintingCondition: types.NewCondition(types.NewMultiSignatureCondition(types.UnlockHashSlice{
			unlockHashFromHex("016148ac9b17828e0933796eaca94418a376f2aa3fefa15685cea5fa462093f0150e09067f7512"),
			unlockHashFromHex("01d553fab496f3fd6092e25ce60e6f72e24b57950bffc0d372d659e38e5a95e89fb117b4eb3481"),
			unlockHashFromHex("013a787bf6248c518aee3a040a14b0dd3a029bc8e9b19a1823faf5bcdde397f4201ad01aace4c9"),
		}, 2)),
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
		GenesisMintingCondition: types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))),
		// belongs to wallet with mnemonic:
		// carbon boss inject cover mountain fetch fiber fit tornado cloth wing dinosaur proof joy intact fabric thumb rebel borrow poet chair network expire else
		FoundationPoolAddress: unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
		// belongs to wallet with mnemonic:
		// carbon boss inject cover mountain fetch fiber fit tornado cloth wing dinosaur proof joy intact fabric thumb rebel borrow poet chair network expire else
		ERC20FeePoolAddress: unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
	}
}
