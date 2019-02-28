package main

import (
	"context"
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
	"github.com/threefoldtech/rivine/pkg/daemon"

	"github.com/ethereum/go-ethereum/log"
	"github.com/threefoldfoundation/tfchain/pkg/api"
	"github.com/threefoldfoundation/tfchain/pkg/cli"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldfoundation/tfchain/pkg/eth/erc20"
	"github.com/threefoldfoundation/tfchain/pkg/persist"
	rivineapi "github.com/threefoldtech/rivine/pkg/api"

	"github.com/spf13/cobra"
	tfchaintypes "github.com/threefoldfoundation/tfchain/pkg/types"
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
	transactionDB     *persist.TransactionDB

	APIaddr   string
	UserAgent string

	VerboseRivineLogging bool
}

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

	log.Info("loading network config, registering types and loading rivine transaction db (0/4)...")
	switch cmd.BlockchainInfo.NetworkName {
	case config.NetworkNameStandard:
		cmd.transactionDB, cmdErr = persist.NewTransactionDB(cmd.rootPerDir(), config.GetStandardnetGenesisMintCondition())
		if cmdErr != nil {
			return fmt.Errorf("failed to create tfchain transaction DB for tfchain standard: %v", cmdErr)
		}
		cmd.ChainConstants = config.GetStandardnetGenesis()
		// Register the transaction controllers for all transaction versions
		// supported on the standard network
		tfchaintypes.RegisterTransactionTypesForStandardNetwork(cmd.transactionDB, &tfchaintypes.NopERC20TransactionValidator{},
			cmd.ChainConstants.CurrencyUnits.OneCoin, config.GetStandardDaemonNetworkConfig())
		// Forbid the usage of MultiSignatureCondition (and thus the multisig feature),
		// until the blockchain reached a height of 42000 blocks.
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(42000)

		if len(cmd.BootstrapPeers) == 0 {
			cmd.BootstrapPeers = config.GetStandardnetBootstrapPeers()
		}

		if cmd.EthNetworkName == "" {
			// default to main network on standard net
			cmd.EthNetworkName = "main"
		}

	case config.NetworkNameTest:
		cmd.transactionDB, cmdErr = persist.NewTransactionDB(cmd.rootPerDir(), config.GetTestnetGenesisMintCondition())
		if cmdErr != nil {
			return fmt.Errorf("failed to create tfchain transaction DB for tfchain testnet: %v", cmdErr)
		}
		// get chain constants and bootstrap peers
		cmd.ChainConstants = config.GetTestnetGenesis()
		// Register the transaction controllers for all transaction versions
		// supported on the test network
		tfchaintypes.RegisterTransactionTypesForTestNetwork(cmd.transactionDB, &tfchaintypes.NopERC20TransactionValidator{},
			cmd.ChainConstants.CurrencyUnits.OneCoin, config.GetTestnetDaemonNetworkConfig())
		// Use our custom MultiSignatureCondition, just for testing purposes
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		if len(cmd.BootstrapPeers) == 0 {
			cmd.BootstrapPeers = config.GetTestnetBootstrapPeers()
		}

		if cmd.EthNetworkName == "" {
			// default to ropsten network on testnet
			cmd.EthNetworkName = "ropsten"
		}

	case config.NetworkNameDev:
		cmd.transactionDB, cmdErr = persist.NewTransactionDB(cmd.rootPerDir(), config.GetDevnetGenesisMintCondition())
		if cmdErr != nil {
			return fmt.Errorf("failed to create tfchain transaction DB for tfchain devnet: %v", cmdErr)
		}
		// get chain constants and bootstrap peers
		cmd.ChainConstants = config.GetDevnetGenesis()
		// Register the transaction controllers for all transaction versions
		// supported on the dev network
		tfchaintypes.RegisterTransactionTypesForDevNetwork(cmd.transactionDB, &tfchaintypes.NopERC20TransactionValidator{},
			cmd.ChainConstants.CurrencyUnits.OneCoin, config.GetDevnetDaemonNetworkConfig())
		// Use our custom MultiSignatureCondition, just for testing purposes
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(0)

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

		log.Info("loading rivine gateway module (1/4)...")
		gateway, err := gateway.New(
			cmd.RPCaddr, true, cmd.perDir("gateway"),
			cmd.BlockchainInfo, cmd.ChainConstants, cmd.BootstrapPeers, cmd.VerboseRivineLogging)
		if err != nil {
			log.Error("Failed to create gateway module", "err", err)
			cancel()
			cmdErr = err
			return
		}
		defer func() {
			log.Info("Closing gateway module...")
			err := gateway.Close()
			if err != nil {
				log.Error("Failed to close gateway module", "err", err)
			}
		}()

		log.Info("loading rivine consensus module (2/4)...")
		cs, err := consensus.New(
			gateway, !cmd.NoBootstrap, cmd.perDir("consensus"),
			cmd.BlockchainInfo, cmd.ChainConstants, cmd.VerboseRivineLogging)
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
		err = cmd.transactionDB.SubscribeToConsensusSet(cs)
		if err != nil {
			log.Error("Failed to subscribe transactionDB module to consensus module", "err", err)
			cancel()
			cmdErr = err
			return
		}
		log.Info("loading transactionpool module (3/4)...")
		tpool, err := transactionpool.New(cs, gateway, cmd.perDir("transactionpool"), cmd.BlockchainInfo, cmd.ChainConstants)
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
		bridged, err := erc20.NewBridge(
			cs, cmd.transactionDB, tpool, cmd.EthPort, cmd.accJSON, cmd.accPass, cmd.EthNetworkName, cmd.EthBootNodes, cmd.ContractAddress, cmd.perDir("bridge"),
			cmd.BlockchainInfo, cmd.ChainConstants, ctx.Done())
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

		// Register ERC20 http handlers
		api.RegisterERC20HTTPHandlers(router, erc20Client)

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

		// handle all our endpoints over a router,
		// which requires a user agent should one be configured
		srv.Handle("/", rivineapi.RequireUserAgentHandler(router, cmd.UserAgent))

		// Wait for the ethereum network to sync
		err = erc20Client.Wait(ctx)
		if err != nil {
			log.Error("error while waing for ERC20 client", "err", err)
			cancel()
			cmdErr = err
			return
		}

		// Start the bridge
		err = bridged.Start(cs, cmd.transactionDB, ctx.Done())
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
		"ethport", 3003,
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

	// execute logic
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
