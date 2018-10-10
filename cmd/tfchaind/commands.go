package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldfoundation/tfchain/pkg/persist"
	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/rivine/rivine/pkg/cli"
	"github.com/rivine/rivine/pkg/daemon"

	"github.com/bgentry/speakeasy"
	"github.com/spf13/cobra"
)

type commands struct {
	cfg           daemon.Config
	moduleSetFlag daemon.ModuleSetFlag
}

func (cmds *commands) rootCommand(*cobra.Command, []string) {
	var err error

	// Silently append a subdirectory for storage with the name of the network so we don't create conflicts
	cmds.cfg.RootPersistentDir = filepath.Join(cmds.cfg.RootPersistentDir, cmds.cfg.BlockchainInfo.NetworkName)

	// Check if we require an api password
	if cmds.cfg.AuthenticateAPI {
		// if its not set, ask one now
		if cmds.cfg.APIPassword == "" {
			// Prompt user for API password.
			cmds.cfg.APIPassword, err = speakeasy.Ask("Enter API password: ")
			if err != nil {
				cli.DieWithError("failed to ask for API password", err)
			}
		}
		if cmds.cfg.APIPassword == "" {
			cli.DieWithError("failed to configure daemon", errors.New("password cannot be blank"))
		}
	} else {
		// If authenticateAPI is not set, explicitly set the password to the empty string.
		// This way the api server maintains consistency with the authenticateAPI var, even if apiPassword is set (possibly by mistake)
		cmds.cfg.APIPassword = ""
	}

	// Process the config variables, cleaning up slightly invalid values
	cmds.cfg = daemon.ProcessConfig(cmds.cfg)

	// run daemon
	err = runDaemon(cmds.cfg, cmds.moduleSetFlag.ModuleIdentifiers())
	if err != nil {
		cli.DieWithError("daemon failed", err)
	}
}

// setupNetwork injects the correct chain constants and genesis nodes based on the chosen network,
// it also ensures that features added during the lifetime of the blockchain,
// only get activated on a certain block height, giving everyone sufficient time to upgrade should such features be introduced,
// it also creates the correct tfchain modules based on the given chain.
func setupNetwork(cfg daemon.Config) (daemon.NetworkConfig, *persist.TransactionDB, error) {
	// return the network configuration, based on the network name,
	// which includes the genesis block as well as the bootstrap peers
	switch cfg.BlockchainInfo.NetworkName {
	case config.NetworkNameStandard:
		txdb, err := persist.NewTransactionDB(cfg.RootPersistentDir, config.GetStandardnetGenesisMintCondition())
		if err != nil {
			return daemon.NetworkConfig{}, nil, err
		}

		constants := config.GetStandardnetGenesis()
		networkConfig := config.GetStandardDaemonNetworkConfig()

		// Register the transaction controllers for all transaction versions
		// supported on the standard network
		types.RegisterTransactionTypesForStandardNetwork(txdb, constants.CurrencyUnits.OneCoin, networkConfig)
		// Forbid the usage of MultiSignatureCondition (and thus the multisig feature),
		// until the blockchain reached a height of 42000 blocks.
		types.RegisterBlockHeightLimitedMultiSignatureCondition(42000)

		// return the standard genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      constants,
			BootstrapPeers: config.GetStandardnetBootstrapPeers(),
		}, txdb, nil

	case config.NetworkNameTest:
		txdb, err := persist.NewTransactionDB(cfg.RootPersistentDir, config.GetTestnetGenesisMintCondition())
		if err != nil {
			return daemon.NetworkConfig{}, nil, err
		}

		constants := config.GetTestnetGenesis()
		networkConfig := config.GetTestnetDaemonNetworkConfig()

		// Register the transaction controllers for all transaction versions
		// supported on the test network
		types.RegisterTransactionTypesForTestNetwork(txdb, constants.CurrencyUnits.OneCoin, networkConfig)
		// Use our custom MultiSignatureCondition, just for testing purposes
		types.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		// return the testnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      constants,
			BootstrapPeers: config.GetTestnetBootstrapPeers(),
		}, txdb, nil

	case config.NetworkNameDev:
		txdb, err := persist.NewTransactionDB(cfg.RootPersistentDir, config.GetDevnetGenesisMintCondition())
		if err != nil {
			return daemon.NetworkConfig{}, nil, err
		}

		constants := config.GetDevnetGenesis()
		networkConfig := config.GetDevnetDaemonNetworkConfig()

		// Register the transaction controllers for all transaction versions
		// supported on the dev network
		types.RegisterTransactionTypesForDevNetwork(txdb, constants.CurrencyUnits.OneCoin, networkConfig)
		// Use our custom MultiSignatureCondition, just for testing purposes
		types.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		// return the devnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      constants,
			BootstrapPeers: nil,
		}, txdb, nil

	default:
		// network isn't recognised
		return daemon.NetworkConfig{}, nil, fmt.Errorf(
			"Netork name %q not recognized", cfg.BlockchainInfo.NetworkName)
	}
}

func (cmds *commands) versionCommand(*cobra.Command, []string) {
	var postfix string
	switch cmds.cfg.BlockchainInfo.NetworkName {
	case "devnet":
		postfix = "-dev"
	case "testnet":
		postfix = "-testing"
	case "standard": // ""
	default:
		postfix = "-???"
	}
	fmt.Printf("%s Daemon v%s%s\n",
		strings.Title(cmds.cfg.BlockchainInfo.Name),
		cmds.cfg.BlockchainInfo.ChainVersion.String(), postfix)
	fmt.Println("Rivine Protocol v" + cmds.cfg.BlockchainInfo.ProtocolVersion.String())

	fmt.Println()
	fmt.Printf("Go Version   v%s\r\n", runtime.Version()[2:])
	fmt.Printf("GOOS         %s\r\n", runtime.GOOS)
	fmt.Printf("GOARCH       %s\r\n", runtime.GOARCH)
}

func (cmds *commands) modulesCommand(*cobra.Command, []string) {
	err := cmds.moduleSetFlag.WriteDescription(os.Stdout)
	if err != nil {
		cli.DieWithError("failed to write usage of the modules flag", err)
	}
}
