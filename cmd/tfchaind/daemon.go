package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/threefoldfoundation/tfchain/pkg/api"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldfoundation/tfchain/pkg/persist"

	tfchaintypes "github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/julienschmidt/httprouter"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/modules/blockcreator"
	"github.com/threefoldtech/rivine/modules/consensus"
	"github.com/threefoldtech/rivine/modules/explorer"
	"github.com/threefoldtech/rivine/modules/gateway"
	"github.com/threefoldtech/rivine/modules/transactionpool"
	"github.com/threefoldtech/rivine/modules/wallet"
	rivineapi "github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/daemon"
)

func runDaemon(cfg ExtendedDaemonConfig, moduleIdentifiers daemon.ModuleIdentifierSet, erc20Cfg ERC20NodeValidatorConfig) error {
	// Print a startup message.
	fmt.Println("Loading...")
	loadStart := time.Now()

	var (
		i             int
		modulesToLoad = moduleIdentifiers.Len()
	)
	printModuleIsLoading := func(name string) {
		fmt.Printf("Loading %s (%d/%d)...\r\n", name, i, modulesToLoad)
		i++
	}

	// create our server already, this way we can fail early if the API addr is already bound
	fmt.Println("Binding API Address and serving the API...")
	srv, err := daemon.NewHTTPServer(cfg.APIaddr)
	if err != nil {
		return err
	}
	servErrs := make(chan error, 32)
	go func() {
		servErrs <- srv.Serve()
	}()

	ctx, cancel := context.WithCancel(context.Background())

	// load all modules

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// router to register all endpoints to
		router := httprouter.New()

		// create the ERC20 Tx Validator, used to validate the ERC20 Coin Creation Transactions
		erc20TxValidator, err := setupERC20TransactionValidator(cfg.RootPersistentDir, cfg.BlockchainInfo.NetworkName, erc20Cfg, ctx.Done())
		if err != nil {
			servErrs <- fmt.Errorf("failed to create ERC20 Transaction validator: %v", err)
			cancel()
			return
		}
		api.RegisterERC20HTTPHandlers(router, erc20TxValidator)

		// create and validate network config, and the transactionDB as well
		// txdb is on index 0, as it is not manually loaded
		printModuleIsLoading("(auto) transaction db")
		networkCfg, txdb, err := setupNetwork(cfg, erc20TxValidator)
		if err != nil {
			servErrs <- fmt.Errorf("failed to create network config: %v", err)
			cancel()
			return
		}
		defer func() {
			fmt.Println("Closing transaction db...")
			err := txdb.Close()
			if err != nil {
				fmt.Println("Error during transactiondb shutdown:", err)
			}
		}()
		err = networkCfg.Constants.Validate()
		if err != nil {
			servErrs <- fmt.Errorf("failed to validate network config: %v", err)
			cancel()
			return
		}
		api.RegisterTransactionDBHTTPHandlers(router, txdb)

		// Initialize the Rivine modules
		var g modules.Gateway
		if moduleIdentifiers.Contains(daemon.GatewayModule.Identifier()) {
			printModuleIsLoading("gateway")
			g, err = gateway.New(cfg.RPCaddr, !cfg.NoBootstrap,
				filepath.Join(cfg.RootPersistentDir, modules.GatewayDir),
				cfg.BlockchainInfo, networkCfg.Constants, networkCfg.BootstrapPeers, cfg.VerboseLogging)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
			rivineapi.RegisterGatewayHTTPHandlers(router, g, cfg.APIPassword)
			defer func() {
				fmt.Println("Closing gateway...")
				err := g.Close()
				if err != nil {
					fmt.Println("Error during gateway shutdown:", err)
				}
			}()

		}
		var cs modules.ConsensusSet
		if moduleIdentifiers.Contains(daemon.ConsensusSetModule.Identifier()) {
			printModuleIsLoading("consensus set")
			cs, err = consensus.New(g, !cfg.NoBootstrap,
				filepath.Join(cfg.RootPersistentDir, modules.ConsensusDir),
				cfg.BlockchainInfo, networkCfg.Constants, cfg.VerboseLogging)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
			rivineapi.RegisterConsensusHTTPHandlers(router, cs)
			defer func() {
				fmt.Println("Closing consensus set...")
				err := cs.Close()
				if err != nil {
					fmt.Println("Error during consensus set shutdown:", err)
				}
			}()
			err = txdb.SubscribeToConsensusSet(cs)
			if err != nil {
				servErrs <- fmt.Errorf("failed to subscribe earlier created transactionDB to the consensus created just now: %v", err)
				cancel()
				return
			}
		}
		var tpool modules.TransactionPool
		if moduleIdentifiers.Contains(daemon.TransactionPoolModule.Identifier()) {
			printModuleIsLoading("transaction pool")
			tpool, err = transactionpool.New(cs, g,
				filepath.Join(cfg.RootPersistentDir, modules.TransactionPoolDir),
				cfg.BlockchainInfo, networkCfg.Constants)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
			rivineapi.RegisterTransactionPoolHTTPHandlers(router, cs, tpool, cfg.APIPassword)
			defer func() {
				fmt.Println("Closing transaction pool...")
				err := tpool.Close()
				if err != nil {
					fmt.Println("Error during transaction pool shutdown:", err)
				}
			}()
		}
		var w modules.Wallet
		if moduleIdentifiers.Contains(daemon.WalletModule.Identifier()) {
			printModuleIsLoading("wallet")
			w, err = wallet.New(cs, tpool,
				filepath.Join(cfg.RootPersistentDir, modules.WalletDir),
				cfg.BlockchainInfo, networkCfg.Constants, cfg.VerboseLogging)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
			rivineapi.RegisterWalletHTTPHandlers(router, w, cfg.APIPassword)
			api.RegisterWalletHTTPHandlers(router, w, cfg.APIPassword)
			defer func() {
				fmt.Println("Closing wallet...")
				err := w.Close()
				if err != nil {
					fmt.Println("Error during wallet shutdown:", err)
				}
			}()

		}
		var b modules.BlockCreator
		if moduleIdentifiers.Contains(daemon.BlockCreatorModule.Identifier()) {
			printModuleIsLoading("block creator")
			b, err = blockcreator.New(cs, tpool, w,
				filepath.Join(cfg.RootPersistentDir, modules.BlockCreatorDir),
				cfg.BlockchainInfo, networkCfg.Constants, cfg.VerboseLogging)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
			// block creator has no API endpoints to register
			defer func() {
				fmt.Println("Closing block creator...")
				err := b.Close()
				if err != nil {
					fmt.Println("Error during block creator shutdown:", err)
				}
			}()
		}
		var e modules.Explorer
		if moduleIdentifiers.Contains(daemon.ExplorerModule.Identifier()) {
			printModuleIsLoading("explorer")
			e, err = explorer.New(cs,
				filepath.Join(cfg.RootPersistentDir, modules.ExplorerDir),
				cfg.BlockchainInfo, networkCfg.Constants)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
			// DO NOT register rivineapi for Explorer HTTP Handles,
			// as they are included in the tfchain api already
			//rivineapi.RegisterExplorerHTTPHandlers(router, cs, e, tpool)
			api.RegisterExplorerHTTPHandlers(router, cs, e, tpool, txdb)
			defer func() {
				fmt.Println("Closing explorer...")
				err := e.Close()
				if err != nil {
					fmt.Println("Error during explorer shutdown:", err)
				}
			}()
		}

		fmt.Println("Setting up root HTTP API handler...")

		// register our special daemon HTTP handlers
		router.GET("/daemon/constants", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			constants := modules.NewDaemonConstants(cfg.BlockchainInfo, networkCfg.Constants)
			rivineapi.WriteJSON(w, constants)
		})
		router.GET("/daemon/version", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			rivineapi.WriteJSON(w, daemon.Version{
				ChainVersion:    cfg.BlockchainInfo.ChainVersion,
				ProtocolVersion: cfg.BlockchainInfo.ProtocolVersion,
			})
		})
		router.POST("/daemon/stop", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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
		srv.Handle("/", rivineapi.RequireUserAgentHandler(router, cfg.RequiredUserAgent))

		// Wait for the ethereum network to sync
		err = erc20TxValidator.Wait(ctx)
		if err != nil {
			servErrs <- err
			cancel()
			return
		}

		if cs != nil {
			cs.Start()
		}

		// Print a 'startup complete' message.
		startupTime := time.Since(loadStart)
		fmt.Println("Finished loading in", startupTime.Seconds(), "seconds")

		// wait until done
		<-ctx.Done()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// wait for server to be killed or the process to be done
	select {
	case <-sigChan:
		fmt.Println("\rCaught stop signal, quitting...")
		srv.Close()
	case <-ctx.Done():
		fmt.Println("\rContext is done, quitting...")
		fmt.Println("context is done, quitting...")
	}

	cancel()
	wg.Wait()

	// return the first error which is returned
	return <-servErrs
}

// setupNetwork injects the correct chain constants and genesis nodes based on the chosen network,
// it also ensures that features added during the lifetime of the blockchain,
// only get activated on a certain block height, giving everyone sufficient time to upgrade should such features be introduced,
// it also creates the correct tfchain modules based on the given chain.
func setupNetwork(cfg ExtendedDaemonConfig, erc20TxValidator tfchaintypes.ERC20TransactionValidator) (daemon.NetworkConfig, *persist.TransactionDB, error) {
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
		tfchaintypes.RegisterTransactionTypesForStandardNetwork(txdb, erc20TxValidator, constants.CurrencyUnits.OneCoin, networkConfig)
		// Forbid the usage of MultiSignatureCondition (and thus the multisig feature),
		// until the blockchain reached a height of 42000 blocks.
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(42000)

		if len(cfg.BootstrapPeers) == 0 {
			cfg.BootstrapPeers = config.GetStandardnetBootstrapPeers()
		}

		// return the standard genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      constants,
			BootstrapPeers: cfg.BootstrapPeers,
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
		tfchaintypes.RegisterTransactionTypesForTestNetwork(txdb, erc20TxValidator, constants.CurrencyUnits.OneCoin, networkConfig)
		// Use our custom MultiSignatureCondition, just for testing purposes
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		if len(cfg.BootstrapPeers) == 0 {
			cfg.BootstrapPeers = config.GetTestnetBootstrapPeers()
		}

		// return the testnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      constants,
			BootstrapPeers: cfg.BootstrapPeers,
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
		tfchaintypes.RegisterTransactionTypesForDevNetwork(txdb, erc20TxValidator, constants.CurrencyUnits.OneCoin, networkConfig)
		// Use our custom MultiSignatureCondition, just for testing purposes
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		if len(cfg.BootstrapPeers) == 0 {
			cfg.BootstrapPeers = config.GetDevnetBootstrapPeers()
		}

		// return the devnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      constants,
			BootstrapPeers: cfg.BootstrapPeers,
		}, txdb, nil

	default:
		// network isn't recognised
		return daemon.NetworkConfig{}, nil, fmt.Errorf(
			"Netork name %q not recognized", cfg.BlockchainInfo.NetworkName)
	}
}

func setupERC20TransactionValidator(rootDir, networkName string, erc20Cfg ERC20NodeValidatorConfig, cancel <-chan struct{}) (tfchaintypes.ERC20TransactionValidator, error) {
	if erc20Cfg.NetworkName == "" {
		switch networkName {
		case config.NetworkNameStandard:
			erc20Cfg.NetworkName = "mainnet"
		case config.NetworkNameTest:
			erc20Cfg.NetworkName = "ropsten"
		default:
			erc20Cfg.NetworkName = "rinkeby"
		}
	}
	erc20Cfg.DataDir = path.Join(rootDir, "leth")
	return NewERC20NodeValidator(erc20Cfg, cancel)
}
