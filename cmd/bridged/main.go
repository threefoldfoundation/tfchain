package main

import (
	_ "net/http/pprof"
	"os"

	"github.com/threefoldfoundation/tfchain/pkg/config"

	"github.com/spf13/cobra"
)

// used to dump the data of a tfchain network in a meaningful way.
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
	NetAddressArrayFlagVar(
		cmdRoot.Flags(),
		&cmd.BootstrapPeers,
		"bootstrap-peer",
		"overwrite the bootstrap peers to use, instead of using the default bootstrap peer (requird for devnet)",
	)
	// Example for running on devnet ./bridged --network devnet --bootstrap-peer "localhost:23112" --bootstrap-peer "localhost:23113"

	// execute logic
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
