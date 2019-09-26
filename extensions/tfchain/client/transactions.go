package client

import (
	"github.com/threefoldfoundation/tfchain/pkg/config"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"

	tbcli "github.com/threefoldfoundation/tfchain/extensions/threebot/client"
	tbtypes "github.com/threefoldfoundation/tfchain/extensions/threebot/types"
	erc20cli "github.com/threefoldtech/rivine-extension-erc20/client"
	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"
	"github.com/threefoldtech/rivine/extensions/minting"
	mintingcli "github.com/threefoldtech/rivine/extensions/minting/client"
	"github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/types"
)

func RegisterStandardTransactions(bc *client.BaseClient) {
	registerTransactions(bc, false, config.GetStandardDaemonNetworkConfig())
}

func RegisterTestnetTransactions(bc *client.BaseClient) {
	registerTransactions(bc, true, config.GetTestnetDaemonNetworkConfig())
}

func RegisterDevnetTransactions(bc *client.BaseClient) {
	registerTransactions(bc, true, config.GetDevnetDaemonNetworkConfig())
}

func registerTransactions(bc *client.BaseClient, extraPlugins bool, daemonCfg config.DaemonNetworkConfig) {
	// create minting plugin client...
	mintingCLI := mintingcli.NewPluginConsensusClient(bc)
	// ...and register minting types
	types.RegisterTransactionVersion(tftypes.TransactionVersionMinterDefinition, minting.MinterDefinitionTransactionController{
		MintingMinerFeeBaseTransactionController: minting.MintingMinerFeeBaseTransactionController{
			MintingBaseTransactionController: minting.MintingBaseTransactionController{
				UseLegacySiaEncoding: true,
			},
			RequireMinerFees: true,
		},
		MintConditionGetter: mintingCLI,
		TransactionVersion:  tftypes.TransactionVersionMinterDefinition,
	})
	types.RegisterTransactionVersion(tftypes.TransactionVersionCoinCreation, minting.CoinCreationTransactionController{
		MintingMinerFeeBaseTransactionController: minting.MintingMinerFeeBaseTransactionController{
			MintingBaseTransactionController: minting.MintingBaseTransactionController{
				UseLegacySiaEncoding: true,
			},
			RequireMinerFees: true,
		},
		MintConditionGetter: mintingCLI,
		TransactionVersion:  tftypes.TransactionVersionCoinCreation,
	})

	if !extraPlugins {
		return // 3Bot and ERC20 transactions are not enabled on all networks
	}

	// register 3Bot Transactions
	tbClient := tbcli.NewPluginConsensusClient(bc)
	types.RegisterTransactionVersion(tbtypes.TransactionVersionBotRegistration, tbtypes.BotRegistrationTransactionController{
		Registry:            tbClient,
		RegistryPoolAddress: daemonCfg.FoundationPoolAddress,
		OneCoin:             bc.Config().CurrencyUnits.OneCoin,
	})
	types.RegisterTransactionVersion(tbtypes.TransactionVersionBotRecordUpdate, tbtypes.BotUpdateRecordTransactionController{
		Registry:            tbClient,
		RegistryPoolAddress: daemonCfg.FoundationPoolAddress,
		OneCoin:             bc.Config().CurrencyUnits.OneCoin,
	})
	types.RegisterTransactionVersion(tbtypes.TransactionVersionBotNameTransfer, tbtypes.BotNameTransferTransactionController{
		Registry:            tbClient,
		RegistryPoolAddress: daemonCfg.FoundationPoolAddress,
		OneCoin:             bc.Config().CurrencyUnits.OneCoin,
	})

	// register ERC20 Transactions
	erc20Client := erc20cli.NewPluginConsensusClient(bc)
	types.RegisterTransactionVersion(tftypes.TransactionVersionERC20Conversion, erc20types.ERC20ConvertTransactionController{
		TransactionVersion: tftypes.TransactionVersionERC20Conversion,
	})
	types.RegisterTransactionVersion(tftypes.TransactionVersionERC20AddressRegistration, erc20types.ERC20AddressRegistrationTransactionController{
		TransactionVersion:   tftypes.TransactionVersionERC20AddressRegistration,
		Registry:             erc20Client,
		BridgeFeePoolAddress: daemonCfg.ERC20FeePoolAddress,
		OneCoin:              bc.Config().CurrencyUnits.OneCoin,
	})
	types.RegisterTransactionVersion(tftypes.TransactionVersionERC20CoinCreation, erc20types.ERC20CoinCreationTransactionController{
		TransactionVersion: tftypes.TransactionVersionERC20CoinCreation,
		Registry:           erc20Client,
		OneCoin:            bc.Config().CurrencyUnits.OneCoin,
	})
}
