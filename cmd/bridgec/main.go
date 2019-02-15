package main

import (
	"fmt"
	"os"

	"github.com/threefoldfoundation/tfchain/cmd/bridgec/internal"
	"github.com/threefoldtech/rivine/pkg/cli"
	"github.com/threefoldtech/rivine/pkg/daemon"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/pkg/client"
)

func main() {
	// create cli
	bchainInfo := config.GetBlockchainInfo()
	cliClient, err := internal.NewCommandLineClient("", bchainInfo.Name, daemon.RivineUserAgent)
	if err != nil {
		panic(err)
	}

	// used for the tfchain-specific controllers
	// txDBReader := internal.NewTransactionDBConsensusClient(cliClient)

	// register root command
	cliClient.ERC20Cmd = createERC20Cmd(cliClient)
	cliClient.RootCmd.AddCommand(cliClient.ERC20Cmd)

	// no ERC20-Tx Validation is done on client-side
	//  nopERC20TxValidator := types.NopERC20TransactionValidator{}

	// define preRun function
	cliClient.PreRunE = func(cfg *client.Config) (*client.Config, error) {
		if cfg == nil {
			bchainInfo := config.GetBlockchainInfo()
			chainConstants := config.GetStandardnetGenesis()
			daemonConstants := modules.NewDaemonConstants(bchainInfo, chainConstants)
			newCfg := client.ConfigFromDaemonConstants(daemonConstants)
			cfg = &newCfg
		}

		switch cfg.NetworkName {
		case config.NetworkNameStandard:
			// networkConfig := config.GetStandardDaemonNetworkConfig()

			// Register the transaction controllers for all transaction versions
			// supported on the standard network
			// types.RegisterTransactionTypesForStandardNetwork(txDBReader, nopERC20TxValidator, cfg.CurrencyUnits.OneCoin, networkConfig)
			// Forbid the usage of MultiSignatureCondition (and thus the multisig feature),
			// until the blockchain reached a height of 42000 blocks.
			types.RegisterBlockHeightLimitedMultiSignatureCondition(42000)

			// overwrite standard network genesis block stamp,
			// as the genesis block is way earlier than the actual first block,
			// due to the hard reset at the bumpy/rough start
			cfg.GenesisBlockTimestamp = 1524168391 // timestamp of (standard) block #1

		case config.NetworkNameTest:
			// networkConfig := config.GetTestnetDaemonNetworkConfig()

			// Register the transaction controllers for all transaction versions
			// supported on the test network
			// types.RegisterTransactionTypesForTestNetwork(txDBReader, nopERC20TxValidator, cfg.CurrencyUnits.OneCoin, networkConfig)
			// Use our custom MultiSignatureCondition, just for testing purposes
			types.RegisterBlockHeightLimitedMultiSignatureCondition(0)

			// seems like testnet timestamp wasn't updated last time it was reset
			cfg.GenesisBlockTimestamp = 1522792547 // timestamp of (testnet) block #1

		case config.NetworkNameDev:
			// networkConfig := config.GetDevnetDaemonNetworkConfig()

			// Register the transaction controllers for all transaction versions
			// supported on the dev network
			// types.RegisterTransactionTypesForDevNetwork(txDBReader, nopERC20TxValidator, cfg.CurrencyUnits.OneCoin, networkConfig)
			// Use our custom MultiSignatureCondition, just for testing purposes
			types.RegisterBlockHeightLimitedMultiSignatureCondition(0)

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
