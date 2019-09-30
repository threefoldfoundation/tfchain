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
	"github.com/threefoldtech/rivine/types"

	tfconsensus "github.com/threefoldfoundation/tfchain/extensions/tfchain/consensus"
	"github.com/threefoldfoundation/tfchain/extensions/threebot"
	tbapi "github.com/threefoldfoundation/tfchain/extensions/threebot/api"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
	erc20 "github.com/threefoldtech/rivine-extension-erc20"
	erc20daemon "github.com/threefoldtech/rivine-extension-erc20/daemon"
	erc20api "github.com/threefoldtech/rivine-extension-erc20/http"
	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"

	"github.com/julienschmidt/httprouter"
	"github.com/threefoldtech/rivine/extensions/minting"
	mintingapi "github.com/threefoldtech/rivine/extensions/minting/api"
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

const (
	// maxConcurrentRPC is the maximum amount of concurrent RPC's to be handled
	// per peer
	maxConcurrentRPC = 1
)

func runDaemon(cfg ExtendedDaemonConfig, moduleIdentifiers daemon.ModuleIdentifierSet, erc20Cfg erc20daemon.ERC20NodeValidatorConfig) error {
	// Print a startup message.
	fmt.Println("Loading...")
	loadStart := time.Now()

	var (
		i             int
		modulesToLoad = moduleIdentifiers.Len()
	)
	printModuleIsLoading := func(name string) {
		fmt.Printf("Loading %s (%d/%d)...\r\n", name, i+1, modulesToLoad)
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

		// create and validate network config
		networkCfg, err := setupNetwork(cfg)
		if err != nil {
			servErrs <- fmt.Errorf("failed to create network config: %v", err)
			cancel()
			return
		}
		err = networkCfg.NetworkConfig.Constants.Validate()
		if err != nil {
			servErrs <- fmt.Errorf("failed to validate network config: %v", err)
			cancel()
			return
		}

		// Initialize the Rivine modules
		var g modules.Gateway
		if moduleIdentifiers.Contains(daemon.GatewayModule.Identifier()) {
			printModuleIsLoading("gateway")
			g, err = gateway.New(cfg.RPCaddr, !cfg.NoBootstrap, maxConcurrentRPC,
				filepath.Join(cfg.RootPersistentDir, modules.GatewayDir),
				cfg.BlockchainInfo, networkCfg.NetworkConfig.Constants, networkCfg.NetworkConfig.BootstrapPeers, cfg.VerboseLogging)
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
		var mintingPlugin *minting.Plugin
		var threebotPlugin *threebot.Plugin
		var erc20TxValidator erc20types.ERC20TransactionValidator
		var erc20Plugin *erc20.Plugin

		if moduleIdentifiers.Contains(daemon.ConsensusSetModule.Identifier()) {
			printModuleIsLoading("consensus set")
			cs, err = consensus.New(g, !cfg.NoBootstrap,
				filepath.Join(cfg.RootPersistentDir, modules.ConsensusDir),
				cfg.BlockchainInfo, networkCfg.NetworkConfig.Constants, cfg.VerboseLogging, cfg.DebugConsensusDB)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}

			cs.SetTransactionValidators(networkCfg.Validators...)
			for txVersion, validators := range networkCfg.MappedValidators {
				cs.SetTransactionVersionMappedValidators(txVersion, validators...)
			}

			rivineapi.RegisterConsensusHTTPHandlers(router, cs)
			defer func() {
				fmt.Println("Closing consensus set...")
				err := cs.Close()
				if err != nil {
					fmt.Println("Error during consensus set shutdown:", err)
				}
			}()

			// create the minting extension plugin
			mintingPlugin = minting.NewMintingPlugin(
				networkCfg.DaemonNetworkConfig.GenesisMintingCondition,
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
			if cfg.BlockchainInfo.NetworkName != config.NetworkNameStandard {
				// create the 3Bot plugin
				var tbPluginOpts *threebot.PluginOptions
				if cfg.BlockchainInfo.NetworkName == config.NetworkNameTest {
					tbPluginOpts = &threebot.PluginOptions{ // TODO: remove this hack once possible (e.g. a testnet network reset)
						HackMinimumBlockHeightSinceDoubleRegistrationsAreForbidden: 350000,
					}
				}
				threebotPlugin = threebot.NewPlugin(
					networkCfg.DaemonNetworkConfig.FoundationPoolAddress,
					networkCfg.NetworkConfig.Constants.CurrencyUnits.OneCoin,
					tbPluginOpts,
				)
				// add the HTTP handlers for the threebot plugin as well
				tbapi.RegisterConsensusHTTPHandlers(router, threebotPlugin)

				// create the ERC20 Tx Validator, used to validate the ERC20 Coin Creation Transactions
				erc20TxValidator, err = setupERC20TransactionValidator(cfg.RootPersistentDir, cfg.BlockchainInfo.NetworkName, erc20Cfg, ctx.Done())
				if err != nil {
					servErrs <- fmt.Errorf("failed to create ERC20 Transaction validator: %v", err)
					cancel()
					return
				}
				// add the HTTP handlers for the ERC20 plugin as well
				erc20api.RegisterERC20HTTPHandlers(router, erc20TxValidator)

				// create the ERC20 plugin
				erc20Plugin = erc20.NewPlugin(
					networkCfg.DaemonNetworkConfig.ERC20FeePoolAddress,
					networkCfg.NetworkConfig.Constants.CurrencyUnits.OneCoin,
					erc20TxValidator,
					erc20types.TransactionVersions{
						ERC20Conversion:          tftypes.TransactionVersionERC20Conversion,
						ERC20AddressRegistration: tftypes.TransactionVersionERC20AddressRegistration,
						ERC20CoinCreation:        tftypes.TransactionVersionERC20CoinCreation,
					},
				)
				// add the HTTP handlers for the ERC20 plugin as well
				erc20api.RegisterConsensusHTTPHandlers(router, erc20Plugin)

				// register the ERC20 Plugin
				err = cs.RegisterPlugin(ctx, "erc20", erc20Plugin)
				if err != nil {
					servErrs <- fmt.Errorf("failed to register the ERC20 extension: %v", err)
					err = erc20Plugin.Close() //make sure any resources are released
					if err != nil {
						fmt.Println("Error during closing of the erc20Plugin:", err)
					}
					cancel()
					return
				}
				// register the Threebot Plugin
				err = cs.RegisterPlugin(ctx, "threebot", threebotPlugin)
				if err != nil {
					servErrs <- fmt.Errorf("failed to register the threebot extension: %v", err)
					err = threebotPlugin.Close() //make sure any resources are released
					if err != nil {
						fmt.Println("Error during closing of the threebotPlugin:", err)
					}
					cancel()
					return
				}
			}
			// register the Minting Plugin
			err = cs.RegisterPlugin(ctx, "minting", mintingPlugin)
			if err != nil {
				servErrs <- fmt.Errorf("failed to register the minting extension: %v", err)
				err = mintingPlugin.Close() //make sure any resources are released
				if err != nil {
					fmt.Println("Error during closing of the mintingPlugin :", err)
				}
				cancel()
				return
			}
		}

		var tpool modules.TransactionPool
		if moduleIdentifiers.Contains(daemon.TransactionPoolModule.Identifier()) {
			printModuleIsLoading("transaction pool")
			tpool, err = transactionpool.New(cs, g,
				filepath.Join(cfg.RootPersistentDir, modules.TransactionPoolDir),
				cfg.BlockchainInfo, networkCfg.NetworkConfig.Constants, cfg.VerboseLogging)
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
				cfg.BlockchainInfo, networkCfg.NetworkConfig.Constants, cfg.VerboseLogging)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
			rivineapi.RegisterWalletHTTPHandlers(router, w, cfg.APIPassword)
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
				cfg.BlockchainInfo, networkCfg.NetworkConfig.Constants, cfg.VerboseLogging)
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
				cfg.BlockchainInfo, networkCfg.NetworkConfig.Constants, cfg.VerboseLogging)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
			defer func() {
				fmt.Println("Closing explorer...")
				err := e.Close()
				if err != nil {
					fmt.Println("Error during explorer shutdown:", err)
				}
			}()

			// DO NOT register rivineapi for Explorer HTTP Handles,
			// as they are included in the tfchain api already
			//rivineapi.RegisterExplorerHTTPHandlers(router, cs, e, tpool)
			api.RegisterExplorerHTTPHandlers(router, cs, e, tpool, threebotPlugin, erc20Plugin)
			mintingapi.RegisterExplorerMintingHTTPHandlers(router, mintingPlugin)
		}

		fmt.Println("Setting up root HTTP API handler...")

		// register our special daemon HTTP handlers
		router.GET("/daemon/constants", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			constants := modules.NewDaemonConstants(cfg.BlockchainInfo, networkCfg.NetworkConfig.Constants)
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

			cancel()
		})

		// handle all our endpoints over a router,
		// which requires a user agent should one be configured
		srv.Handle("/", rivineapi.RequireUserAgentHandler(router, cfg.RequiredUserAgent))

		// 3Bot and ERC20 is not yet to be used on network standard
		if cfg.BlockchainInfo.NetworkName != config.NetworkNameStandard {
			// Wait for the ethereum network to sync
			err = erc20TxValidator.Wait(ctx)
			if err != nil {
				servErrs <- err
				cancel()
				return
			}
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

type setupNetworkConfig struct {
	NetworkConfig       daemon.NetworkConfig
	DaemonNetworkConfig config.DaemonNetworkConfig

	Validators       []modules.TransactionValidationFunction
	MappedValidators map[types.TransactionVersion][]modules.TransactionValidationFunction
}

// setupNetwork injects the correct chain constants and genesis nodes based on the chosen network,
// it also ensures that features added during the lifetime of the blockchain,
// only get activated on a certain block height, giving everyone sufficient time to upgrade should such features be introduced,
// it also creates the correct tfchain modules based on the given chain.
func setupNetwork(cfg ExtendedDaemonConfig) (setupNetworkConfig, error) {
	// return the network configuration, based on the network name,
	// which includes the genesis block as well as the bootstrap peers
	switch cfg.BlockchainInfo.NetworkName {
	case config.NetworkNameStandard:
		constants := config.GetStandardnetGenesis()
		networkConfig := config.GetStandardDaemonNetworkConfig()

		// Get the bootstrap peers from the hardcoded config, if none are given
		if len(cfg.BootstrapPeers) == 0 {
			cfg.BootstrapPeers = config.GetStandardnetBootstrapPeers()
		}

		// return all info needed to setup the standard network
		return setupNetworkConfig{
			NetworkConfig: daemon.NetworkConfig{
				Constants:      constants,
				BootstrapPeers: cfg.BootstrapPeers,
			},
			Validators:          tfconsensus.GetStandardTransactionValidators(),
			MappedValidators:    tfconsensus.GetStandardTransactionVersionMappedValidators(),
			DaemonNetworkConfig: networkConfig,
		}, nil

	case config.NetworkNameTest:
		constants := config.GetTestnetGenesis()
		networkConfig := config.GetTestnetDaemonNetworkConfig()

		// Get the bootstrap peers from the config
		if len(cfg.BootstrapPeers) == 0 {
			cfg.BootstrapPeers = config.GetTestnetBootstrapPeers()
		}

		// return all info needed to setup the testnet network
		return setupNetworkConfig{
			NetworkConfig: daemon.NetworkConfig{
				Constants:      constants,
				BootstrapPeers: cfg.BootstrapPeers,
			},
			Validators:          tfconsensus.GetTestnetTransactionValidators(),
			MappedValidators:    tfconsensus.GetTestnetTransactionVersionMappedValidators(),
			DaemonNetworkConfig: networkConfig,
		}, nil

	case config.NetworkNameDev:
		constants := config.GetDevnetGenesis()
		networkConfig := config.GetDevnetDaemonNetworkConfig()

		// Get the bootstrap peers from the config
		if len(cfg.BootstrapPeers) == 0 {
			cfg.BootstrapPeers = config.GetDevnetBootstrapPeers()
		}

		// return all info needed to setup the devnet network
		return setupNetworkConfig{
			NetworkConfig: daemon.NetworkConfig{
				Constants:      constants,
				BootstrapPeers: cfg.BootstrapPeers,
			},
			Validators:          tfconsensus.GetDevnetTransactionValidators(),
			MappedValidators:    tfconsensus.GetDevnetTransactionVersionMappedValidators(),
			DaemonNetworkConfig: networkConfig,
		}, nil

	default:
		// network isn't recognised
		return setupNetworkConfig{}, fmt.Errorf(
			"Network name %q not recognized", cfg.BlockchainInfo.NetworkName)
	}
}

func setupERC20TransactionValidator(rootDir, networkName string, erc20Cfg erc20daemon.ERC20NodeValidatorConfig, cancel <-chan struct{}) (erc20types.ERC20TransactionValidator, error) {
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
	return erc20daemon.NewERC20NodeValidator(erc20Cfg, cancel)
}
