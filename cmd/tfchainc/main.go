package main

import (
	"fmt"
	"os"

	"github.com/threefoldtech/rivine/pkg/cli"
	"github.com/threefoldtech/rivine/pkg/daemon"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/pkg/client"

	tfcli "github.com/threefoldfoundation/tfchain/extensions/tfchain/client"
	tbcli "github.com/threefoldfoundation/tfchain/extensions/threebot/client"
	erc20cli "github.com/threefoldtech/rivine-extension-erc20/client"
	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"
	authcointxcli "github.com/threefoldtech/rivine/extensions/authcointx/client"
	mintingcli "github.com/threefoldtech/rivine/extensions/minting/client"
)

func main() {
	// create cli
	bchainInfo := config.GetBlockchainInfo()
	cliClient, err := NewCommandLineClient("", bchainInfo.Name, daemon.RivineUserAgent)
	exitIfError(err)

	// register tfchain-specific commands
	err = mintingcli.CreateConsensusCmd(cliClient.CommandLineClient)
	exitIfError(err)
	err = tbcli.CreateConsensusSubCmds(cliClient.CommandLineClient)
	exitIfError(err)
	err = mintingcli.CreateExploreCmd(cliClient.CommandLineClient)
	exitIfError(err)
	err = tbcli.CreateExplorerSubCmds(cliClient.CommandLineClient)
	exitIfError(err)
	err = mintingcli.CreateWalletCmds(
		cliClient.CommandLineClient,
		tftypes.TransactionVersionMinterDefinition, tftypes.TransactionVersionCoinCreation,
		&mintingcli.WalletCmdsOpts{
			CoinDestructionTxVersion: 0,    // disabled
			RequireMinerFees:         true, // require miner fees
		})
	exitIfError(err)
	err = erc20cli.CreateWalletCmds(cliClient.CommandLineClient, erc20types.TransactionVersions{
		ERC20Conversion:          tftypes.TransactionVersionERC20Conversion,
		ERC20AddressRegistration: tftypes.TransactionVersionERC20AddressRegistration,
		ERC20CoinCreation:        tftypes.TransactionVersionERC20CoinCreation,
	})
	exitIfError(err)
	err = tbcli.CreateWalletCmds(cliClient.CommandLineClient)
	exitIfError(err)
	erc20cli.CreateERC20Cmd(cliClient.CommandLineClient)

	err = authcointxcli.CreateConsensusAuthCoinInfoCmd(cliClient.CommandLineClient)
	exitIfError(err)
	err = authcointxcli.CreateExploreAuthCoinInfoCmd(cliClient.CommandLineClient)
	exitIfError(err)
	authcointxcli.CreateWalletCmds(
		cliClient.CommandLineClient,
		tftypes.TransactionVersionAuthConditionUpdate,
		tftypes.TransactionVersionAuthAddressUpdate,
		&authcointxcli.WalletCmdsOpts{
			RequireMinerFees: true, // require miner fees
		},
	)

	// register root command
	cliClient.ERC20Cmd = erc20cli.CreateERC20Cmd(cliClient.CommandLineClient)
	cliClient.RootCmd.AddCommand(cliClient.ERC20Cmd)

	// define preRun function
	cliClient.PreRunE = func(cfg *client.Config) (*client.Config, error) {
		if cfg == nil {
			bchainInfo := config.GetBlockchainInfo()
			chainConstants := config.GetStandardnetGenesis()
			daemonConstants := modules.NewDaemonConstants(bchainInfo, chainConstants, nil)
			newCfg := client.ConfigFromDaemonConstants(daemonConstants)
			cfg = &newCfg
		}

		bc, err := client.NewLazyBaseClient(func() (client.BaseClient, error) {
			return client.NewBaseClient(cliClient.HTTPClient, cfg)
		})
		if err != nil {
			return nil, err
		}

		switch cfg.NetworkName {
		case config.NetworkNameStandard:
			// Register the transaction controllers for all transaction versions
			// supported on the standard network
			err = tfcli.RegisterStandardTransactions(bc)
			if err != nil {
				return nil, err
			}

			// overwrite standard network genesis block stamp,
			// as the genesis block is way earlier than the actual first block,
			// due to the hard reset at the bumpy/rough start
			cfg.GenesisBlockTimestamp = 1524168391 // timestamp of (standard) block #1

		case config.NetworkNameTest:
			// Register the transaction controllers for all transaction versions
			// supported on the test network
			err = tfcli.RegisterTestnetTransactions(bc)
			if err != nil {
				return nil, err
			}

			// seems like testnet timestamp wasn't updated last time it was reset
			cfg.GenesisBlockTimestamp = 1522792547 // timestamp of (testnet) block #1

		case config.NetworkNameDev:
			// Register the transaction controllers for all transaction versions
			// supported on the dev network
			err = tfcli.RegisterDevnetTransactions(bc)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("Netork name %q not recognized", cfg.NetworkName)
		}

		return cfg, nil
	}

	// start cli
	if err := cliClient.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "client exited with an error: ", err)
		// Since no commands return errors (all commands set Command.Run instead of
		// Command.RunE), Command.Execute() should only return an error on an
		// invalid command or flag. Therefore Command.Usage() was called (assuming
		// Command.SilenceUsage is false) and we should exit with exitCodeUsage.
		os.Exit(cli.ExitCodeUsage)
	}
}

func exitIfError(err error) {
	if err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, "client exited during setup with an error:", err)
	os.Exit(cli.ExitCodeGeneral)
}
