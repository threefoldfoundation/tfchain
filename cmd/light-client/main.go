package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/threefoldfoundation/tfchain/cmd/light-client/wallet"

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
	DataString               string
}

func main() {

	var cmd cmds

	rootCmd := &cobra.Command{
		Use:   "tfchain-light [wallet name]",
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
	recoverCmd.Flags().Uint64Var(&cmd.KeysToLoad, "key-amount", DefaultKeysToLoad, "Set the default amount of keys to load")

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
			Long: `Print the seed of this wallet as a mnemonic. This mnemonic can be stored and later used to recover
	the wallet`,
			RunE: cmd.walletSeed,
		}

		txCmd := &cobra.Command{
			Use:   "send <address> <amount>",
			Short: "Send coins using a transaction",
			Long: `Create a new transaction to send the specified amount of coins to the specified address.
	Inputs are selected automatically from the available ones. The transactionfee is set to the lowest permitted value.
	In case the sum of the inputs warants a refund output, that is also added automatically as well.`,
			RunE: cmd.walletSend,
			Args: cobra.ExactArgs(2),
		}
		txCmd.Flags().BoolVar(&cmd.GenerateNewRefundAddress, "new-refund-addr", false, "Generate a new refund address instead of reusing an existing address")
		txCmd.Flags().StringVarP(&cmd.DataString, "data", "d", "", "Attach this string as arbitrary data to the transaction")

		addressesCmd := &cobra.Command{
			Use:   "addresses",
			Short: "List all loaded addresses",
			Long: `List all loaded addresses. If an address owned by this wallet is not in the list after recovering,
	you can load more addresses using the 'load' sub command`,
			RunE: cmd.walletAddresses,
			Args: cobra.NoArgs,
		}

		loadCmd := &cobra.Command{
			Use:   "load <amount>",
			Short: "Generate additional keys",
			Long: `Generate additional keys. The index used to generate the keys is continued from the allready loaded keys. After loading the wallet
	persistent data is updated to reflect the additional keys, and they will be available for future use.`,
			RunE: cmd.walletLoad,
			Args: cobra.ExactArgs(1),
		}

		walletCmd.AddCommand(seedCmd, txCmd, addressesCmd, loadCmd)
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
