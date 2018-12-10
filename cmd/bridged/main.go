package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"sync"

	_ "net/http/pprof"

	"github.com/decred/dcrwallet/version"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldfoundation/tfchain/pkg/persist"

	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/modules/consensus"
	"github.com/threefoldtech/rivine/modules/gateway"
	rivinetypes "github.com/threefoldtech/rivine/types"

	"github.com/spf13/cobra"
)

// used to dump the data of a tfchain network in a meaningful way.
type Bridged struct {
	cs   modules.ConsensusSet
	txdb *persist.TransactionDB

	bcInfo   rivinetypes.BlockchainInfo
	chainCts rivinetypes.ChainConstants

	mut sync.Mutex
}

// Create new Bridged.
func NewBridged(cs modules.ConsensusSet, txdb *persist.TransactionDB, bcInfo rivinetypes.BlockchainInfo, chainCts rivinetypes.ChainConstants, cancel <-chan struct{}) (*Bridged, error) {

	bridged := &Bridged{
		cs:       cs,
		txdb:     txdb,
		bcInfo:   bcInfo,
		chainCts: chainCts,
	}
	err := cs.ConsensusSetSubscribe(bridged, modules.ConsensusChangeRecent, cancel)
	if err != nil {
		return nil, fmt.Errorf("explorer: failed to subscribe to consensus set: %v", err)
	}
	return bridged, nil
}

// Close bridged.
func (bridged *Bridged) Close() {
	bridged.mut.Lock()
	defer bridged.mut.Unlock()
	bridged.cs.Unsubscribe(bridged)
}

// ProcessConsensusChange implements modules.ConsensusSetSubscriber,
// used to apply/revert blocks.
func (bridged *Bridged) ProcessConsensusChange(css modules.ConsensusChange) {
	bridged.mut.Lock()
	defer bridged.mut.Unlock()

	// var err error

	// update reverted blocks
	for _, block := range css.RevertedBlocks {
		fmt.Println("block reverted: ", block)
	}

	// update applied blocks
	for _, block := range css.AppliedBlocks {
		fmt.Println("block applied: ", block)
	}

}

type Commands struct {
	RPCaddr        string
	BlockchainInfo rivinetypes.BlockchainInfo
	ChainConstants rivinetypes.ChainConstants
	BootstrapPeers []modules.NetAddress

	RootPersistentDir string
	transactionDB     *persist.TransactionDB
}

// Root represents the root (`bridged`) command,
// starting a bridged daemon instance, running until the user intervenes.
func (cmd *Commands) Root(_ *cobra.Command, args []string) (cmdErr error) {
	peers := []modules.NetAddress{}
	log.Println("starting bridged v" + version.String() + "...")

	log.Println("loading network config, registering types and loading rivine transaction db (0/3)...")
	switch cmd.BlockchainInfo.NetworkName {
	case config.NetworkNameStandard:
		cmd.transactionDB, cmdErr = persist.NewTransactionDB(cmd.rootPerDir(), config.GetStandardnetGenesisMintCondition())
		if cmdErr != nil {
			return fmt.Errorf("failed to create tfchain transaction DB for tfchain standard: %v", cmdErr)
		}

	case config.NetworkNameTest:
		cmd.transactionDB, cmdErr = persist.NewTransactionDB(cmd.rootPerDir(), config.GetTestnetGenesisMintCondition())
		if cmdErr != nil {
			return fmt.Errorf("failed to create tfchain transaction DB for tfchain testnet: %v", cmdErr)
		}

	case config.NetworkNameDev:
		cmd.transactionDB, cmdErr = persist.NewTransactionDB(cmd.rootPerDir(), config.GetDevnetGenesisMintCondition())
		if cmdErr != nil {
			return fmt.Errorf("failed to create tfchain transaction DB for tfchain devnet: %v", cmdErr)
		}
		// get chain constants and bootstrap peers
		cmd.ChainConstants = config.GetDevnetGenesis()

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

		log.Println("loading rivine gateway module (1/3)...")
		gateway, err := gateway.New(
			cmd.RPCaddr, true, cmd.perDir("gateway"),
			cmd.BlockchainInfo, cmd.ChainConstants, peers)
		if err != nil {
			cmdErr = fmt.Errorf("failed to create gateway module: %v", err)
			log.Println("[ERROR] ", cmdErr)
			cancel()
			return
		}
		defer func() {
			log.Println("Closing gateway module...")
			err := gateway.Close()
			if err != nil {
				cmdErr = err
				log.Println("[ERROR] Closing gateway module resulted in an error: ", err)
			}
		}()

		log.Println("loading rivine consensus module (2/3)...")
		cs, err := consensus.New(
			gateway, true, cmd.perDir("consensus"),
			cmd.BlockchainInfo, cmd.ChainConstants)
		if err != nil {
			cmdErr = fmt.Errorf("failed to create consensus module: %v", err)
			log.Println("[ERROR] ", cmdErr)
			cancel()
			return
		}
		defer func() {
			log.Println("Closing consensus module...")
			err := cs.Close()
			if err != nil {
				cmdErr = err
				log.Println("[ERROR] Closing consensus module resulted in an error: ", err)
			}
		}()
		err = cmd.transactionDB.SubscribeToConsensusSet(cs)
		if err != nil {
			cmdErr = fmt.Errorf("failed to subscribe earlier created transactionDB to the consensus created just now: %v", err)
			log.Println("[ERROR] ", cmdErr)
			cancel()
			return
		}

		log.Println("bridged is up and running...")

		// wait until done
		<-ctx.Done()
	}()

	// stop the server if a kill signal is caught
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// wait for server to be killed or the process to be done
	select {
	case <-sigChan:
		log.Println("Caught stop signal, quitting...")
	case <-ctx.Done():
		log.Println("context is done, quitting...")
	}

	cancel()
	wg.Wait()

	log.Println("Goodbye!")
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
	fmt.Printf("Bridged version            v%s\n", version.String())
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
		cmd.RootPersistentDir,
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
		"the name of the network to which the daemon connects, one of {standard,testnet,devnet}",
	)

	// execute logic
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
