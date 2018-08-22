package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivine/rivine/pkg/cli"
	"github.com/rivine/rivine/pkg/daemon"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/bgentry/speakeasy"
	"github.com/spf13/cobra"
)

type commands struct {
	cfg           daemon.Config
	moduleSetFlag daemon.ModuleSetFlag
}

func (cmds *commands) rootCommand(*cobra.Command, []string) {
	// create and validate network config
	networkCfg, err := setupNetworksAndTypes(cmds.cfg.BlockchainInfo.NetworkName)
	if err != nil {
		cli.DieWithError("failed to create network config", err)
	}
	err = networkCfg.Constants.Validate()
	if err != nil {
		cli.DieWithError("failed to validate network config", err)
	}

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
	err = runDaemon(cmds.cfg, networkCfg, cmds.moduleSetFlag.ModuleIdentifiers())
	if err != nil {
		cli.DieWithError("daemon failed", err)
	}
}

// setupNetworksAndTypes injects the correct chain constants and genesis nodes based on the chosen network,
// it also ensures that features added during the lifetime of the blockchain,
// only get activated on a certain block height, giving everyone sufficient time to upgrade should such features be introduced.
func setupNetworksAndTypes(name string) (daemon.NetworkConfig, error) {
	// return the network configuration, based on the network name,
	// which includes the genesis block as well as the bootstrap peers
	switch name {
	case config.NetworkNameStandard:
		// Register the transaction controllers for all transaction versions
		// supported on the standard network
		types.RegisterTransactionTypesForStandardNetwork()
		// Forbid the usage of MultiSignatureCondition (and thus the multisig feature),
		// until the blockchain reached a height of 42000 blocks.
		types.RegisterBlockHeightLimitedMultiSignatureCondition(42000)

		// return the standard genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      config.GetStandardnetGenesis(),
			BootstrapPeers: config.GetStandardnetBootstrapPeers(),
		}, nil

	case config.NetworkNameTest:
		// Register the transaction controllers for all transaction versions
		// supported on the test network
		types.RegisterTransactionTypesForTestNetwork()
		// Use our custom MultiSignatureCondition, just for testing purposes
		types.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		// return the testnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      config.GetTestnetGenesis(),
			BootstrapPeers: config.GetTestnetBootstrapPeers(),
		}, nil

	case config.NetworkNameDev:
		// Register the transaction controllers for all transaction versions
		// supported on the dev network
		types.RegisterTransactionTypesForDevNetwork()
		// Use our custom MultiSignatureCondition, just for testing purposes
		types.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		// return the devnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      config.GetDevnetGenesis(),
			BootstrapPeers: nil,
		}, nil

	default:
		// network isn't recognised
		return daemon.NetworkConfig{}, fmt.Errorf("Netork name %q not recognized", name)
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
}

func (cmds *commands) modulesCommand(*cobra.Command, []string) {
	err := cmds.moduleSetFlag.WriteDescription(os.Stdout)
	if err != nil {
		cli.DieWithError("failed to write usage of the modules flag", err)
	}
}
