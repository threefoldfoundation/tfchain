package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/threefoldfoundation/tfchain/pkg/api"

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

		// Wait for the ethereum network to sync
		err = erc20TxValidator.Wait(ctx)
		if err != nil {
			servErrs <- err
			cancel()
			return
		}

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
