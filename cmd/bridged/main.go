package main

import (
	"context"
	"fmt"
	"math/big"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/threefoldtech/rivine/modules/transactionpool"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	bridgedeth "github.com/threefoldfoundation/tfchain/cmd/bridged/eth"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldfoundation/tfchain/pkg/persist"
	"github.com/threefoldtech/rivine/types"

	"github.com/spf13/cobra"
	tfchaintypes "github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/modules/consensus"
	"github.com/threefoldtech/rivine/modules/gateway"
	rivinetypes "github.com/threefoldtech/rivine/types"
)

// used to dump the data of a tfchain network in a meaningful way.
type Bridged struct {
	cs   modules.ConsensusSet
	txdb *persist.TransactionDB
	tp   modules.TransactionPool

	bcInfo   rivinetypes.BlockchainInfo
	chainCts rivinetypes.ChainConstants

	bridge *ethBridge

	mut sync.Mutex
}

// Create new Bridged.
func NewBridged(cs modules.ConsensusSet, txdb *persist.TransactionDB, tp modules.TransactionPool, bcInfo rivinetypes.BlockchainInfo, chainCts rivinetypes.ChainConstants, ethPort uint16, accountJSON, accountPass string, datadir string, ethNetworkName string, cancel <-chan struct{}) (*Bridged, error) {

	bridge, err := newEthBridge(ethNetworkName, int(ethPort), accountJSON, accountPass, filepath.Join(datadir, "eth"))
	if err != nil {
		return nil, err
	}

	bridged := &Bridged{
		cs:       cs,
		txdb:     txdb,
		tp:       tp,
		bcInfo:   bcInfo,
		chainCts: chainCts,
		bridge:   bridge,
	}
	err = cs.ConsensusSetSubscribe(bridged, txdb.GetLastConsensusChangeID(), cancel)
	if err != nil {
		return nil, fmt.Errorf("bridged: failed to subscribe to consensus set: %v", err)
	}
	networkConfig, err := bridgedeth.GetEthNetworkConfiguration(ethNetworkName)
	ContractAddress := networkConfig.ContractAddress
	go bridged.bridge.loop()
	go bridged.bridge.SubscribeTransfers(ContractAddress)
	go bridged.bridge.SubscribeMint(ContractAddress)
	go bridged.bridge.SubscribeRegisterWithdrawAddress(ContractAddress)

	withdrawChan := make(chan WithdrawEvent)
	go bridged.bridge.SubscribeWithdraw(ContractAddress, withdrawChan)
	go func() {
		addr := types.UnlockHash{}
		addr.LoadString("01fcf7fadb91c63ff6ff9ac5fd1405eee0ae8b29c7f0a95e7e8772bd02dc6e77df4b74a631487b")
		for {
			we := <-withdrawChan
			tx := tfchaintypes.ERC20CoinCreationTransaction{}
			// Todo: dynamic address
			tx.Address = addr
			tx.Value = types.NewCurrency(we.amount)
			tx.TransactionID = tfchaintypes.ERC20TransactionID(we.txHash)
			tx.TransactionFee = types.NewCurrency(OneToken)
			if err := tp.AcceptTransactionSet([]types.Transaction{tx.Transaction()}); err != nil {
				log.Error("Failed to push ERC20 -> TFT transaction", "err", err)
				return
			}
		}
	}()

	return bridged, nil
}

// Close bridged.
func (bridged *Bridged) Close() {
	bridged.mut.Lock()
	defer bridged.mut.Unlock()
	bridged.bridge.close()
	bridged.cs.Unsubscribe(bridged)
}

// ProcessConsensusChange implements modules.ConsensusSetSubscriber,
// used to apply/revert blocks.
func (bridged *Bridged) ProcessConsensusChange(css modules.ConsensusChange) {
	bridged.mut.Lock()
	defer bridged.mut.Unlock()

	// TODO: add delay

	for _, block := range css.AppliedBlocks {
		for _, tx := range block.Transactions {
			if tx.Version == tfchaintypes.TransactionVersionERC20Conversion {
				log.Warn("Found convert transacton")
				txConvert, err := tfchaintypes.ERC20ConvertTransactionFromTransaction(tx)
				if err != nil {
					log.Error("Found a TFT convert transaction version, but can't create a conversion transaction from it")
					return
				}
				// Send the mint transaction, this requires gas
				if err = bridged.Mint(txConvert.Address, txConvert.Value, tx.ID()); err != nil {
					log.Error("Failed to push mint transaction", "error", err)
					return
				}
				log.Info("Created mint transaction on eth network")
			} else if tx.Version == tfchaintypes.TransactionVersionERC20AddressRegistration {
				log.Warn("Found erc20 address registration")
				txRegistration, err := tfchaintypes.ERC20AddressRegistrationTransactionFromTransaction(tx)
				if err != nil {
					log.Error("Found a TFT ERC20 Address registration transaction version, but can't create the right transaction for it")
					return
				}
				// send the address registration transaction
				if err = bridged.RegisterWithdrawalAddress(txRegistration.PublicKey); err != nil {
					log.Error("Failed to push withdrawal address registration transaction", "err", err)
					return
				}
				log.Info("Registered withdrawal address on eth network")
			}
		}
	}
}

var (
	// 18 digit precision
	precision = big.NewInt(1000000000000000000)
)

func (bridged *Bridged) Mint(receiver tfchaintypes.ERC20Address, amount types.Currency, txID types.TransactionID) error {
	return bridged.bridge.Mint(bridged.bridge.GetContractAdress(), common.Address(receiver), big.NewInt(0).Mul(amount.Big(), precision), txID.String())
}

func (bridged *Bridged) RegisterWithdrawalAddress(key types.PublicKey) error {
	// use the first 20 bytes from the key for now
	var addressBytes [20]byte
	copy(addressBytes[:], key.Key[:20])
	return bridged.bridge.RegisterWithdrawalAddress(bridged.bridge.GetContractAdress(), common.Address(addressBytes))
}

type Commands struct {
	RPCaddr        string
	BlockchainInfo rivinetypes.BlockchainInfo
	ChainConstants rivinetypes.ChainConstants
	BootstrapPeers []modules.NetAddress

	EthNetworkName string

	// eth port for light client
	EthPort uint16

	// eth account flags
	accJSON string
	accPass string

	EthLog int

	RootPersistentDir string
	transactionDB     *persist.TransactionDB
}

func getDevnetBootstrapPeers() []modules.NetAddress {
	return []modules.NetAddress{
		"localhost:23112",
	}
}

// Root represents the root (`bridged`) command,
// starting a bridged daemon instance, running until the user intervenes.
func (cmd *Commands) Root(_ *cobra.Command, args []string) (cmdErr error) {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(cmd.EthLog), log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	log.Info("starting bridged  0.1.0")

	log.Info("loading network config, registering types and loading rivine transaction db (0/3)...")
	switch cmd.BlockchainInfo.NetworkName {
	case config.NetworkNameStandard:
		cmd.transactionDB, cmdErr = persist.NewTransactionDB(cmd.rootPerDir(), config.GetStandardnetGenesisMintCondition())
		if cmdErr != nil {
			return fmt.Errorf("failed to create tfchain transaction DB for tfchain standard: %v", cmdErr)
		}
		cmd.ChainConstants = config.GetStandardnetGenesis()
		// Register the transaction controllers for all transaction versions
		// supported on the standard network
		tfchaintypes.RegisterTransactionTypesForStandardNetwork(cmd.transactionDB,
			cmd.ChainConstants.CurrencyUnits.OneCoin, config.GetStandardDaemonNetworkConfig())
		// Forbid the usage of MultiSignatureCondition (and thus the multisig feature),
		// until the blockchain reached a height of 42000 blocks.
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(42000)

		if len(cmd.BootstrapPeers) == 0 {
			cmd.BootstrapPeers = config.GetStandardnetBootstrapPeers()
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
		tfchaintypes.RegisterTransactionTypesForTestNetwork(cmd.transactionDB,
			cmd.ChainConstants.CurrencyUnits.OneCoin, config.GetTestnetDaemonNetworkConfig())
		// Use our custom MultiSignatureCondition, just for testing purposes
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		if len(cmd.BootstrapPeers) == 0 {
			cmd.BootstrapPeers = config.GetTestnetBootstrapPeers()
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
		tfchaintypes.RegisterTransactionTypesForDevNetwork(cmd.transactionDB,
			cmd.ChainConstants.CurrencyUnits.OneCoin, config.GetDevnetDaemonNetworkConfig())
		// Use our custom MultiSignatureCondition, just for testing purposes
		tfchaintypes.RegisterBlockHeightLimitedMultiSignatureCondition(0)

		cmd.BootstrapPeers = getDevnetBootstrapPeers()
	default:
		return fmt.Errorf(
			"%q is an invalid network name, has to be one of {standard,testnet,devnet}",
			cmd.BlockchainInfo.Name)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// load all modules

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Info("loading rivine gateway module (1/4)...")
		gateway, err := gateway.New(
			cmd.RPCaddr, true, cmd.perDir("gateway"),
			cmd.BlockchainInfo, cmd.ChainConstants, cmd.BootstrapPeers)
		if err != nil {
			cmdErr = fmt.Errorf("failed to create gateway module: %v", err)
			log.Error("[ERROR] ", cmdErr)
			cancel()
			return
		}
		defer func() {
			log.Info("Closing gateway module...")
			err := gateway.Close()
			if err != nil {
				cmdErr = err
				log.Error("[ERROR] Closing gateway module resulted in an error: ", err)
			}
		}()

		log.Info("loading rivine consensus module (2/4)...")
		cs, err := consensus.New(
			gateway, true, cmd.perDir("consensus"),
			cmd.BlockchainInfo, cmd.ChainConstants)
		if err != nil {
			cmdErr = fmt.Errorf("failed to create consensus module: %v", err)
			log.Error("[ERROR] ", cmdErr)
			cancel()
			return
		}
		defer func() {
			log.Info("Closing consensus module...")
			err := cs.Close()
			if err != nil {
				cmdErr = err
				log.Error("[ERROR] Closing consensus module resulted in an error: ", err)
			}
		}()
		err = cmd.transactionDB.SubscribeToConsensusSet(cs)
		if err != nil {
			cmdErr = fmt.Errorf("failed to subscribe earlier created transactionDB to the consensus created just now: %v", err)
			log.Error("[ERROR] ", cmdErr)
			cancel()
			return
		}
		log.Info("loading transactionpool module (3/4)...")
		tpool, err := transactionpool.New(cs, gateway, cmd.perDir("transactionpool"), cmd.BlockchainInfo, cmd.ChainConstants)
		if err != nil {
			cmdErr = fmt.Errorf("failed to create transactionpool module")
			log.Error("Failed to create txpool module", "err", err)
			cancel()
			return
		}
		defer func() {
			log.Info("Closing transactionpool module...")
			err := tpool.Close()
			if err != nil {
				cmdErr = err
				log.Error("Failed to close the transactionpool module", "err", err)
			}
		}()

		log.Info("loading bridged module (4/4)...")
		bridged, err := NewBridged(
			cs, cmd.transactionDB, tpool, cmd.BlockchainInfo, cmd.ChainConstants, cmd.EthPort, cmd.accJSON, cmd.accPass, cmd.RootPersistentDir, cmd.EthNetworkName, ctx.Done())
		if err != nil {
			cmdErr = fmt.Errorf("failed to create bridged module: %v", err)
			log.Error("[ERROR] ", cmdErr)
			cancel()
			return
		}
		defer func() {
			log.Info("closing bridged module...")
			bridged.Close()

		}()
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
func (cmd *Commands) Version(_ *cobra.Command, args []string) {
	fmt.Printf("Bridged version            1.2\n")
	fmt.Printf("TFChain Daemon version  v%s\n", cmd.BlockchainInfo.ChainVersion.String())
	fmt.Printf("Rivine protocol version v%s\n", cmd.BlockchainInfo.ProtocolVersion.String())
	fmt.Println()
	fmt.Printf("Go Version   v%s\n", runtime.Version()[2:])
	fmt.Printf("GOOS         %s\n", runtime.GOOS)
	fmt.Printf("GOARCH       %s\n", runtime.GOARCH)

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
		"location of the root diretory used to store persistent data of the daemon of "+cmd.BlockchainInfo.Name,
	)
	cmdRoot.Flags().StringVar(
		&cmd.RPCaddr,
		"rpc-addr",
		cmd.RPCaddr,
		"which port the gateway listens on",
	)

	cmdRoot.Flags().StringVarP(
		&cmd.BlockchainInfo.NetworkName,
		"network", "n",
		cmd.BlockchainInfo.NetworkName,
		"the name of the tfchain network to  connect to  {standard,testnet,devnet}",
	)

	// bridge flags

	cmdRoot.Flags().StringVar(
		&cmd.EthNetworkName,
		"ethnetwork", "main",
		"The ethereum network, {main,rinkeby, ropsten}",
	)
	cmdRoot.Flags().Uint16Var(
		&cmd.EthPort,
		"ethport", 3003,
		"port for the ethereum deamon",
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
		"Ethereum-log-lvl", "e", 3,
		"Log lvl for the ethereum logger",
	)

	// execute logic
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
