package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/threefoldfoundation/tfchain/cmd/tfchaint/wallet"

	"github.com/spf13/cobra"
)

const (
	// DefaultUserAgent is the default user agent string for explorer based backends
	DefaultUserAgent = "Rivine-Agent"
	// DefaultKeysToLoad is the default amount of keys to load
	DefaultKeysToLoad = 1
)

type cmds struct {
	KeysToLoad               uint64
	GenerateNewRefundAddress bool
	MultiSig                 bool
	DataString               string
	LockString               string
	Network                  string
	Broker                   string
}

func main() {
	var cmd cmds

	rootCmd := &cobra.Command{
		Use:   "tfchaint [wallet name]",
		Short: "Tfchain command line light wallet",
		Long: `A command line based light wallet for Threefold Chain. This application uses a locally
stored seed, and gets the required blockchain info from public explorers. This way you can manage a wallet
without having to download the entire blockchain.`,
	}

	initCmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new wallet with the given name",
		Long:  `Initialize a new wallet with the given name. A seed for this wallet will be generated for you.`,
		RunE:  cmd.walletInit,
		Args:  cobra.ExactArgs(1),
	}
	recoverCmd := &cobra.Command{
		Use:   "recover [name]",
		Short: "Recovers a wallet from an existing seed",
		Long:  "Recover a wallet from an existing seed. This will add a wallet with the given name and the given seed.",
		RunE:  cmd.walletRecover,
		Args:  cobra.ExactArgs(1),
	}
	initCmd.Flags().Uint64Var(&cmd.KeysToLoad, "key-amount", DefaultKeysToLoad, "Set the default amount of keys to load")
	initCmd.Flags().StringVar(&cmd.Network, "network", "testnet", "Set the network to use for this wallet")
	recoverCmd.Flags().Uint64Var(&cmd.KeysToLoad, "key-amount", DefaultKeysToLoad, "Set the default amount of keys to load")
	recoverCmd.Flags().StringVar(&cmd.Network, "network", "testnet", "Set the network to use for this wallet")

	rootCmd.AddCommand(
		initCmd,
		recoverCmd,
	)

	walletNames, err := listWallets()
	if err != nil {
		fmt.Println("Failed to retrieve wallets:", err)
		return
	}

	for _, walletName := range walletNames {
		walletCmd := &cobra.Command{
			Use:   walletName,
			Short: fmt.Sprintf("Get the balance of wallet %s", walletName),
			Long: fmt.Sprintf(`Get the balance of wallet %v.
Additional actions can be performed on this wallet via subcommands`, walletName),
			RunE: cmd.walletBalance,
			Args: cobra.NoArgs,
		}

		rootCmd.AddCommand(walletCmd)

		seedCmd := &cobra.Command{
			Use:   "seed",
			Short: "Print the seed of this wallet",
			Long:  `Print the seed of this wallet as a mnemonic. This mnemonic can be stored and later used to recover the wallet`,
			RunE:  cmd.walletSeed,
		}

		txCmd := &cobra.Command{
			Use:   "send <amount> <address> ...",
			Short: "Send coins using a transaction",
			Long: `Create a new transaction to send the specified amount of coins to the specified address.
Inputs are selected automatically from the available ones. The transactionfee is set to the lowest permitted value.
In case the sum of the inputs warants a refund output, that is also added automatically as well.

The following formats are supported to identify the receiver:
	- empty string: send to the nil address, anyone can spend
	- single address: send to an address. The owner can spend the funds
	- multiple addresses: send to a multisig address composed of all the addresses given, everyone needs to sign to spend
	- multiple addresses + integer: send to a multisig address composed of all the addresses given, the integer identifies how many parties need to sign to spend`,
			RunE: cmd.walletSend,
			Args: cobra.MinimumNArgs(2),
		}
		txCmd.Flags().BoolVar(&cmd.GenerateNewRefundAddress, "new-refund-addr", false, "Generate a new refund address instead of reusing an existing address")
		txCmd.Flags().BoolVar(&cmd.MultiSig, "multisig", false, "Send coins to a multisignature address")
		txCmd.Flags().StringVarP(&cmd.DataString, "data", "d", "", "Attach this string as arbitrary data to the transaction")
		txCmd.Flags().StringVarP(&cmd.LockString, "lock", "l", "", "Optional time lock. Supported formats are: <integer>, <data>, <date time> <duration>")

		reserveCmd := &cobra.Command{
			Use:   "reserve <type> <size> <email>",
			Short: "Create a reservation transaction",
			Long: `Create a reservation transaction. The exact cost of the reserved
workload is automatically set. The email address is used to receive the connection
info once the reservation has been processed by the broker, identified by the broker
address. For a full overview of the available workloads and their price, see
https://github.com/threefoldtech/grid_broker`,
			// RunE: cmd.walletReserve,
			Args: cobra.ExactArgs(3),
		}
		reserveCmd.PersistentFlags().BoolVar(&cmd.GenerateNewRefundAddress, "new-refund-addr", false, "Generate a new refund address instead of reusing an existing address")
		reserveCmd.PersistentFlags().StringVarP(&cmd.Broker, "broker", "b", "", "Use a custom broker instead of the default public one")

		reserveVMCmd := &cobra.Command{
			Use:   "vm <size> <nodeid> <email>",
			Short: "Reserve a vm on the threefold grid",
			Long: `Create a transaction which attempts to reserve a vm. The exact
cost of the reserved vm is automatically set. The email address is used to receive
the connection info once the vm has been deployed by the broker. For a full overview
of the available sizes and their price, see https://github.com/threefoldtech/grid_broker`,
			RunE: cmd.walletReserveVM,
			Args: cobra.ExactArgs(3),
		}

		reserveS3Cmd := &cobra.Command{
			Use:   "s3 <size> <farm_name> <email>",
			Short: "Reserve an s3 instance on the theefold grid",
			Long: `Create a transaction which attmepts to reserve an s3 instance.
The exact cost of the reserved instance is automatically set. The email address
is used to receive the connection info once the s3 has been deployed by the
broker. For a full overview of the available sizes and their price, see
https://github.com/threefoldtoken/grid_broker`,
			RunE: cmd.walletReserveS3,
			Args: cobra.ExactArgs(3),
		}
		reserveCmd.AddCommand(reserveVMCmd, reserveS3Cmd)

		addressesCmd := &cobra.Command{
			Use:   "addresses",
			Short: "List all loaded addresses",
			Long: `List all loaded addresses. If an address owned by this wallet is not in the list after recovering,
	you can load more addresses using the 'load' sub command`,
			RunE: cmd.walletAddresses,
			Args: cobra.NoArgs,
		}

		generateCmd := &cobra.Command{
			Use:   "generate <amount>",
			Short: "Generate additional addresses",
			Long: `Generate additional keys. The index used to generate the keys is continued from the allready loaded keys. After loading the wallet
persistent data is updated to reflect the additional keys, and they will be available for future use.

If no amount is specified, 1 address will be generated`,
			RunE: cmd.walletLoad,
			Args: cobra.MaximumNArgs(1),
		}
		addressesCmd.AddCommand(generateCmd)
		walletCmd.AddCommand(seedCmd, txCmd, reserveCmd, addressesCmd)
	}

	rootCmd.Execute()
}

func listWallets() ([]string, error) {
	dirs, err := ioutil.ReadDir(wallet.PersistDir())
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	names := []string{}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		names = append(names, dir.Name())
	}
	return names, nil
}
