package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/threefoldtech/rivine/modules/transactionpool"
	"github.com/threefoldtech/rivine/pkg/cli"
	"github.com/threefoldtech/rivine/pkg/daemon"

	"github.com/ethereum/go-ethereum/log"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"
	rivineapi "github.com/threefoldtech/rivine/pkg/api"

	"github.com/threefoldtech/rivine/extensions/minting"
	mintingapi "github.com/threefoldtech/rivine/extensions/minting/api"

	"github.com/threefoldfoundation/tfchain/extensions/threebot"
	bpapi "github.com/threefoldfoundation/tfchain/extensions/threebot/api"
	erc20 "github.com/threefoldtech/rivine-extension-erc20"
	erc20bridge "github.com/threefoldtech/rivine-extension-erc20/api/bridge"
	erc20daemon "github.com/threefoldtech/rivine-extension-erc20/daemon"
	erc20api "github.com/threefoldtech/rivine-extension-erc20/http"

	"github.com/spf13/cobra"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/modules/consensus"
	"github.com/threefoldtech/rivine/modules/gateway"
	rivinetypes "github.com/threefoldtech/rivine/types"
)

// Commands defines the CLI Commands for the Bridge as well as its in-memory state.
type Commands struct {
	RPCaddr        string
	BlockchainInfo rivinetypes.BlockchainInfo
	ChainConstants rivinetypes.ChainConstants
	NetworkConfig  config.DaemonNetworkConfig
	BootstrapPeers []modules.NetAddress
	NoBootstrap    bool

	EthNetworkName string

	// eth port for light client
	EthPort uint16

	// eth bootnodes
	EthBootNodes []string

	// eth account flags
	accJSON string
	accPass string

	EthLog          int
	ContractAddress string

	RootPersistentDir string
	erc20Registry     erc20types.ERC20Registry

	APIaddr   string
	UserAgent string

	VerboseRivineLogging bool
	ConsensusDebugFile   string
}

const (
	// maxConcurrentRPC is the maximum amount of concurrent RPC's to be handled
	// per peer
	maxConcurrentRPC = 1
)

// Root represents the root (`bridged`) command,
// starting a bridged daemon instance, running until the user intervenes.
func (cmd *Commands) Root(_ *cobra.Command, args []string) (cmdErr error) {
	// Define the Ethereum Logger,
	// logging both to a file and the STDERR, with a lower verbosity for the latter.
	logFileDir := cmd.perDir("bridge")
	// Create the directory if it doesn't exist.
	err := os.MkdirAll(logFileDir, 0700)
	if err != nil {
		return err
	}
	ethLogFmtr := log.TerminalFormat(true)
	ethLogFileHandler, err := log.FileHandler(path.Join(logFileDir, "bridge.log"), ethLogFmtr)
	if err != nil {
		return fmt.Errorf("failed to create bridge: error while setting up ETH file-logger: %v", err)
	}

	termLogLvl := log.Lvl(cmd.EthLog)
	fileLogLvl := termLogLvl
	if fileLogLvl < log.LvlTrace {
		fileLogLvl++
	}

	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.Lvl(fileLogLvl), ethLogFileHandler),
		log.LvlFilterHandler(log.Lvl(termLogLvl), log.StreamHandler(os.Stderr, ethLogFmtr))))

	log.Info("starting bridge", "version", cmd.BlockchainInfo.ChainVersion.String())

	log.Info("loading network config (0/4)...")
	switch cmd.BlockchainInfo.NetworkName {
	case config.NetworkNameStandard:
		// cmd.ChainConstants = config.GetStandardnetGenesis()
		// cmd.NetworkConfig = config.GetStandardDaemonNetworkConfig()

		// if len(cmd.BootstrapPeers) == 0 {
		// 	cmd.BootstrapPeers = config.GetStandardnetBootstrapPeers()
		// }

		// if cmd.EthNetworkName == "" {
		// 	// default to main network on standard net
		// 	cmd.EthNetworkName = "main"
		// }
		return errors.New("ERC20 feature is currently not enabled on standard net")

	case config.NetworkNameTest:
		cmd.ChainConstants = config.GetTestnetGenesis()
		cmd.NetworkConfig = config.GetTestnetDaemonNetworkConfig()

		if len(cmd.BootstrapPeers) == 0 {
			cmd.BootstrapPeers = config.GetTestnetBootstrapPeers()
		}

		if cmd.EthNetworkName == "" {
			// default to ropsten network on testnet
			cmd.EthNetworkName = "ropsten"
		}

	case config.NetworkNameDev:
		cmd.ChainConstants = config.GetDevnetGenesis()
		cmd.NetworkConfig = config.GetDevnetDaemonNetworkConfig()

		if len(cmd.BootstrapPeers) == 0 {
			cmd.BootstrapPeers = config.GetDevnetBootstrapPeers()
		}

		if cmd.EthNetworkName == "" {
			// default to rinkeby network on devnet
			cmd.EthNetworkName = "rinkeby"
		}

	default:
		return fmt.Errorf(
			"%q is an invalid network name, has to be one of {standard,testnet,devnet}",
			cmd.BlockchainInfo.Name)
	}

	err = cmd.ChainConstants.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate network config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// create our server already, this way we can fail early if the API addr is already bound
	fmt.Println("Binding API Address and serving the API...")
	srv, err := daemon.NewHTTPServer(cmd.APIaddr)
	if err != nil {
		return err
	}
	servErrs := make(chan error, 32)
	go func() {
		servErrs <- srv.Serve()
	}()
	// load all modules

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// router to register all endpoints to
		router := httprouter.New()

		var cs modules.ConsensusSet

		// handle all our endpoints over a router,
		// which requires a user agent should one be configured
		srv.Handle("/", rivineapi.RequireUserAgentHandler(router, cmd.UserAgent))

		// register our special bridge HTTP handlers
		router.GET("/daemon/constants", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			var pluginNames []string
			if cs != nil {
				pluginNames = cs.LoadedPlugins()
			}
			constants := modules.NewDaemonConstants(cmd.BlockchainInfo, cmd.ChainConstants, pluginNames)
			rivineapi.WriteJSON(w, constants)
		})
		router.POST("/bridge/stop", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			// can't write after we stop the server, so lie a bit.
			rivineapi.WriteSuccess(w)

			// need to flush the response before shutting down the server
			f, ok := w.(http.Flusher)
			if !ok {
				panic("Server does not support flushing")
			}
			f.Flush()

			if err := srv.Close(); err != nil {
				servErrs <- err
			}
		})

		log.Info("loading rivine gateway module (1/4)...")
		gateway, err := gateway.New(
			cmd.RPCaddr, true, maxConcurrentRPC, cmd.perDir("gateway"),
			cmd.BlockchainInfo, cmd.ChainConstants, cmd.BootstrapPeers, cmd.VerboseRivineLogging)
		if err != nil {
			log.Error("Failed to create gateway module", "err", err)
			cancel()
			cmdErr = err
			return
		}
		// Blank password as we are not exposing the bridge HTTP API.
		// TODO: Proper password verification like in the regular daemon
		rivineapi.RegisterGatewayHTTPHandlers(router, gateway, "")
		defer func() {
			log.Info("Closing gateway module...")
			err := gateway.Close()
			if err != nil {
				log.Error("Failed to close gateway module", "err", err)
			}
		}()

		log.Info("loading rivine consensus module (2/4)...")
		cs, err = consensus.New(
			gateway, !cmd.NoBootstrap, cmd.perDir("consensus"),
			cmd.BlockchainInfo, cmd.ChainConstants, cmd.VerboseRivineLogging, cmd.ConsensusDebugFile)
		if err != nil {
			log.Error("Failed to create consensus module", "err", err)
			cancel()
			cmdErr = err
			return
		}
		rivineapi.RegisterConsensusHTTPHandlers(router, cs)
		defer func() {
			log.Info("Closing consensus module...")
			err := cs.Close()
			if err != nil {
				log.Error("Failed to close consensus module", "err", err)
			}
		}()

		var (
			mintingPlugin  *minting.Plugin
			threebotPlugin *threebot.Plugin
		)

		// create the minting extension plugin
		mintingPlugin = minting.NewMintingPlugin(
			cmd.NetworkConfig.GenesisMintingCondition,
			tftypes.TransactionVersionMinterDefinition,
			tftypes.TransactionVersionCoinCreation,
			&minting.PluginOptions{
				UseLegacySiaEncoding: true,
				RequireMinerFees:     true,
			},
		)
		// add the HTTP handlers for the minting plugin as well
		mintingapi.RegisterConsensusMintingHTTPHandlers(router, mintingPlugin)

		// 3Bot and ERC20 is not yet to be used on network standard
		// create the 3Bot plugin
		var tbPluginOpts *threebot.PluginOptions
		if cmd.BlockchainInfo.NetworkName == config.NetworkNameTest {
			tbPluginOpts = &threebot.PluginOptions{ // TODO: remove this hack once possible (e.g. a testnet network reset)
				HackMinimumBlockHeightSinceDoubleRegistrationsAreForbidden: 350000,
			}
		}
		threebotPlugin = threebot.NewPlugin(
			cmd.NetworkConfig.FoundationPoolAddress,
			cmd.ChainConstants.CurrencyUnits.OneCoin,
			tbPluginOpts,
		)
		// add the HTTP handlers for the threebot plugin as well
		bpapi.RegisterConsensusHTTPHandlers(router, threebotPlugin)

		log.Info("loading transactionpool module (3/4)...")
		tpool, err := transactionpool.New(cs, gateway, cmd.perDir("transactionpool"), cmd.BlockchainInfo, cmd.ChainConstants, cmd.VerboseRivineLogging)
		if err != nil {
			log.Error("Failed to create txpool module", "err", err)
			cancel()
			cmdErr = err
			return
		}
		defer func() {
			log.Info("Closing transactionpool module...")
			err := tpool.Close()
			if err != nil {
				log.Error("Failed to close the transactionpool module", "err", err)
			}
		}()

		log.Info("loading bridged module (4/4)...")
		bridged, err := erc20bridge.NewBridge(
			cs, tpool, cmd.EthPort, cmd.accJSON, cmd.accPass, cmd.EthNetworkName, cmd.EthBootNodes, cmd.ContractAddress, cmd.perDir("bridge"),
			cmd.BlockchainInfo, cmd.ChainConstants, erc20types.TransactionVersions{
				ERC20Conversion:          tftypes.TransactionVersionERC20Conversion,
				ERC20AddressRegistration: tftypes.TransactionVersionERC20AddressRegistration,
				ERC20CoinCreation:        tftypes.TransactionVersionERC20CoinCreation,
			}, ctx.Done())
		if err != nil {
			log.Error("Failed to create bridge module", "err", err)
			cancel()
			cmdErr = err
			return
		}
		defer func() {
			log.Info("closing bridged module...")
			err := bridged.Close()
			if err != nil {
				log.Error("Failed to close bridge module", "err", err)
			}
		}()

		erc20Client := bridged.GetClient()
		erc20BridgeContract := bridged.GetBridgeContract()
		erc20NodeValidator, err := erc20daemon.NewERC20NodeValidatorFromBridgeContract(erc20BridgeContract)
		if err != nil {
			log.Error("failed to create ERC20 bridge node validator", "err", err)
			cancel()
			cmdErr = err
			return
		}

		// register the ERC20 plugin
		erc20Plugin := erc20.NewPlugin(
			cmd.NetworkConfig.ERC20FeePoolAddress,
			cmd.ChainConstants.CurrencyUnits.OneCoin,
			erc20NodeValidator, erc20types.TransactionVersions{
				ERC20Conversion:          tftypes.TransactionVersionERC20Conversion,
				ERC20AddressRegistration: tftypes.TransactionVersionERC20AddressRegistration,
				ERC20CoinCreation:        tftypes.TransactionVersionERC20CoinCreation,
			})

		// register ERC20 plugin and HTTP handlers
		// add the HTTP handlers for the ERC20 plugin as well
		erc20api.RegisterERC20HTTPHandlers(router, erc20NodeValidator)

		err = cs.RegisterPlugin(ctx, "erc20", erc20Plugin)
		if err != nil {
			servErrs <- fmt.Errorf("failed to register the ERC20 extension: %v", err)
			err = erc20Plugin.Close() //make sure any resources are released
			if err != nil {
				log.Error("Error during closing of the erc20Plugin", "err", err)
			}
			cancel()
			return
		}
		// add the HTTP handlers for the ERC20 plugin as well
		erc20api.RegisterConsensusHTTPHandlers(router, erc20Plugin)

		// register the threebot plugin
		err = cs.RegisterPlugin(ctx, "threebot", threebotPlugin)
		if err != nil {
			servErrs <- fmt.Errorf("failed to register the threebot extension: %v", err)
			err = threebotPlugin.Close() //make sure any resources are released
			if err != nil {
				log.Error("Error during closing of the threebotPlugin", "err", err)
			}
			cancel()
			return
		}

		// register the minting plugin
		err = cs.RegisterPlugin(ctx, "minting", mintingPlugin)
		if err != nil {
			servErrs <- fmt.Errorf("failed to register the minting extension: %v", err)
			err = mintingPlugin.Close() //make sure any resources are released
			if err != nil {
				log.Error("Error during closing of the mintingPlugin", "err", err)
			}
			cancel()
			return
		}

		// Wait for the ethereum network to sync
		err = erc20Client.Wait(ctx)
		if err != nil {
			log.Error("error while waing for ERC20 client", "err", err)
			cancel()
			cmdErr = err
			return
		}

		// Start the cs after the eth module is synced
		cs.Start()

		// Start the bridge
		err = bridged.Start(cs, erc20Plugin, ctx.Done())
		if err != nil {
			log.Error("error while starting the ERC20 client", "err", err)
			cancel()
			cmdErr = err
			return
		}

		log.Info("bridged is up and running...")

		// wait until done
		<-ctx.Done()
	}()

	// stop the server if a kill signal is caught
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// wait for server to be killed or the process to be done
	select {
	case <-sigChan:
		log.Info("Caught stop signal, quitting...")
	case <-ctx.Done():
		log.Info("context is done, quitting...")
	case err = <-servErrs:
		log.Error("Error while serving API", "err", err)
	}

	cancel()
	wg.Wait()

	log.Info("Goodbye!")
	return
}

func (cmd *Commands) rootPerDir() string {
	return path.Join(
		cmd.RootPersistentDir,
		cmd.BlockchainInfo.Name, cmd.BlockchainInfo.NetworkName)
}

func (cmd *Commands) perDir(module string) string {
	return path.Join(cmd.rootPerDir(), module)
}

// Version represents the version (`bridged version`) command,
// returning the version of the tool, dependencies and Go,
// as well as the OS and Arch type.
func (cmd *Commands) Version(_ *cobra.Command, _ []string) {
	var postfix string
	switch cmd.BlockchainInfo.NetworkName {
	case "devnet":
		postfix = "-dev"
	case "testnet":
		postfix = "-testing"
	case "standard": // ""
	default:
		postfix = "-???"
	}
	fmt.Printf("%s Bridge Daemon v%s%s\n",
		strings.Title(cmd.BlockchainInfo.Name),
		cmd.BlockchainInfo.ChainVersion.String(), postfix)
	fmt.Println("Rivine Protocol v" + cmd.BlockchainInfo.ProtocolVersion.String())

	fmt.Println()
	fmt.Printf("Go Version   v%s\r\n", runtime.Version()[2:])
	fmt.Printf("GOOS         %s\r\n", runtime.GOOS)
	fmt.Printf("GOARCH       %s\r\n", runtime.GOARCH)
}

func main() {
	cmd := new(Commands)
	cmd.RPCaddr = ":23118"
	cmd.BlockchainInfo = config.GetBlockchainInfo()

	// define commands
	cmdRoot := &cobra.Command{
		Use:          "bridged",
		Short:        "start the bridged daemon",
		Long:         `start the bridged daemon`,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE:         cmd.Root,
	}

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "show versions of this tool",
		Args:  cobra.ExactArgs(0),
		Run:   cmd.Version,
	}

	// define command tree
	cmdRoot.AddCommand(
		cmdVersion,
	)

	// define flags
	cmdRoot.Flags().StringVarP(
		&cmd.RootPersistentDir,
		"persistent-directory", "d",
		"bridgedata",
		"location of the root directory used to store persistent data of the daemon of "+cmd.BlockchainInfo.Name,
	)
	cmdRoot.Flags().StringVar(
		&cmd.RPCaddr,
		"rpc-addr",
		cmd.RPCaddr,
		"which port the gateway listens on",
	)
	cmdRoot.Flags().BoolVar(
		&cmd.NoBootstrap,
		"no-bootstrap",
		cmd.NoBootstrap,
		"disable bootstrapping on this run for tfchain",
	)

	cmdRoot.Flags().StringVarP(
		&cmd.BlockchainInfo.NetworkName,
		"network", "n",
		cmd.BlockchainInfo.NetworkName,
		"the name of the tfchain network to  connect to  {standard,testnet,devnet}",
	)

	cli.NetAddressArrayFlagVar(
		cmdRoot.Flags(),
		&cmd.BootstrapPeers,
		"bootstrap-peers",
		"override the default tfchain bootstrap peers",
	)

	// bridge flags

	cmdRoot.Flags().StringVar(
		&cmd.EthNetworkName,
		"ethnetwork", "",
		"The ethereum network, {main, rinkeby, ropsten}, defaults to the TFT-linked network",
	)
	cmdRoot.Flags().Uint16Var(
		&cmd.EthPort,
		"ethport", 30302,
		"port for the ethereum deamon",
	)
	cmdRoot.Flags().StringSliceVar(
		&cmd.EthBootNodes,
		"ethbootnodes", nil,
		"Override the default ethereum bootnodes, a comma seperated list of enode URLs (enode://pubkey1@ip1:port1)",
	)

	// bridge account
	cmdRoot.Flags().StringVar(
		&cmd.accJSON,
		"account-json", "",
		"the path to an account file. If set, the specified account will be loaded",
	)

	cmdRoot.Flags().StringVar(
		&cmd.accPass,
		"account-password", "",
		"Password for the bridge account",
	)

	cmdRoot.Flags().IntVarP(
		&cmd.EthLog,
		"ethereum-log-lvl", "e", 3,
		"Log lvl for the ethereum logger",
	)

	cmdRoot.Flags().StringVar(
		&cmd.ContractAddress,
		"contract-address", "",
		"Use a custom contract",
	)

	cmdRoot.Flags().StringVar(
		&cmd.APIaddr,
		"api-address", "localhost:23111",
		"Set custom api-address for bridged",
	)

	cmdRoot.Flags().StringVar(
		&cmd.UserAgent,
		"user-agent", daemon.RivineUserAgent,
		"Set custom User-Agent",
	)
	cmdRoot.Flags().BoolVarP(&cmd.VerboseRivineLogging, "verboseRivinelogging", "v", false, "enable verboselogging in the logfiles of the rivine modules")
	cmdRoot.Flags().StringVar(&cmd.ConsensusDebugFile, "consensus-db-stats", cmd.ConsensusDebugFile, "file path in which json encoded database stats will be saved")

	// execute logic
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
