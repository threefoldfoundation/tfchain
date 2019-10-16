package client

import (
	"encoding/json"
	"fmt"
	"os"

	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"

	rivineapi "github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/cli"
	"github.com/threefoldtech/rivine/pkg/client"
	rivinecli "github.com/threefoldtech/rivine/pkg/client"
	rivinetypes "github.com/threefoldtech/rivine/types"

	"github.com/spf13/cobra"
)

// CreateWalletCmds creates the ERC20 wallet root command as well as its transaction creation sub commands.
func CreateWalletCmds(ccli *client.CommandLineClient, txVersions erc20types.TransactionVersions) error {
	bc, err := client.NewLazyBaseClientFromCommandLineClient(ccli)
	if err != nil {
		return err
	}
	walletCmd := &walletCmd{
		cli:          ccli,
		txVersions:   txVersions,
		walletClient: rivinecli.NewWalletClient(bc),
		txPoolClient: rivinecli.NewTransactionPoolClient(bc),
		erc20Client:  NewPluginExplorerClient(bc),
	}

	// define commands
	var (
		sendERC20FundsCmd = &cobra.Command{
			Use:   "erc20funds erc20_address amount",
			Short: "Convert TFT to ERC20 funds and send those to an ERC20 adddress (minus fees)",
			Run:   rivinecli.Wrap(walletCmd.sendERC20Funds),
		}

		sendERC20FundsClaimCmd = &cobra.Command{
			Use:   "erc20fundsclaim tft_address amount erc20_blockid erc20_txid",
			Short: "Convert ERC20 funds to TFT",
			Run:   rivinecli.Wrap(walletCmd.sendERC20FundsClaim),
		}

		sendERC20AddressRegistrationCmd = &cobra.Command{
			Use:   "erc20address [public_key]",
			Short: "Register an ERC20 address linked to a TFT Public Key",
			Long: `Register an ERC20 address linked to a TFT Public Key

If no Public Key is given, a new one will be generated
using the unlocked wallet from the tfchain daemon.
`,
			Run: walletCmd.sendERC20AddressRegistration,
		}

		listERC20AddressesCmd = &cobra.Command{
			Use:   "erc20addresses",
			Short: "List all known ERC20 addresses for this wallet",
			Long: `List all known ERC20 addresses for this wallet.

An address is considered as owned by this wallet, if the public key used to derive the TFT address,
from which the ERC20 address is then derived, is known by this wallet.`,

			Run: walletCmd.listERC20AddressRegistrations,
		}
	)

	// add commands as wallet sub commands
	ccli.WalletCmd.RootCmdSend.AddCommand(
		sendERC20FundsCmd,
		sendERC20FundsClaimCmd,
		sendERC20AddressRegistrationCmd,
	)

	ccli.WalletCmd.RootCmdList.AddCommand(
		listERC20AddressesCmd,
	)

	// register flags
	sendERC20FundsCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletCmd.sendERC20FundsCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	sendERC20FundsClaimCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletCmd.sendERC20FundsClaimCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	sendERC20AddressRegistrationCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletCmd.sendERC20AddressRegistrationCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	listERC20AddressesCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletCmd.listERC20AddressRegistrationsCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	return nil
}

type walletCmd struct {
	cli        *rivinecli.CommandLineClient
	txVersions erc20types.TransactionVersions

	walletClient *rivinecli.WalletClient
	txPoolClient *rivinecli.TransactionPoolClient
	erc20Client  *PluginClient

	sendERC20FundsCfg struct {
		EncodingType cli.EncodingType
	}
	sendERC20FundsClaimCfg struct {
		EncodingType cli.EncodingType
	}
	sendERC20AddressRegistrationCfg struct {
		EncodingType cli.EncodingType
	}
	listERC20AddressRegistrationsCfg struct {
		EncodingType cli.EncodingType
	}
}

func (walletCmd *walletCmd) sendERC20Funds(hexAddress, strAmount string) {
	// load ERC20 address
	var address erc20types.ERC20Address
	err := address.LoadString(hexAddress)
	if err != nil {
		cli.DieWithError("failed to parse hex-encoded ERC20 address", err)
		return
	}
	// load amount (in TFT)
	currencyConvertor := walletCmd.cli.CreateCurrencyConvertor()
	amount, err := currencyConvertor.ParseCoinString(strAmount)
	if err != nil {
		cli.DieWithError("failed to parse coin (TFT) string", err)
		return
	}

	// create the ERC20 Convert Tx
	tx := erc20types.ERC20ConvertTransaction{
		Address:        address,
		Value:          amount,
		TransactionFee: walletCmd.cli.Config.MinimumTransactionFee,
	}
	// fund the coin inputs
	tx.CoinInputs, tx.RefundCoinOutput, err = walletCmd.walletClient.FundCoins(tx.TransactionFee.Add(tx.Value), nil, false)
	if err != nil {
		cli.DieWithError("failed to fund the ERC20 Convert Txd", err)
		return
	}

	// sign the Tx
	rtx := tx.Transaction(walletCmd.txVersions.ERC20Conversion)
	err = walletCmd.walletClient.GreedySignTx(&rtx)
	if err != nil {
		cli.DieWithError("failed to sign the ERC20 Convert Tx", err)
		return
	}

	// submit the Tx
	txID, err := walletCmd.txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the ERC20 Convert Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	switch walletCmd.sendERC20FundsCfg.EncodingType {
	case cli.EncodingTypeHuman:
		fmt.Println("Transaction ID:", txID)
	case cli.EncodingTypeJSON:
		err = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"transactionid": txID,
		})
		if err != nil {
			cli.DieWithError("failed to encode result", err)
		}
	}
}

func (walletCmd *walletCmd) sendERC20FundsClaim(hexAddress, strAmount, hexBlockID, hexTransctionID string) {
	// load TFT address
	var address rivinetypes.UnlockHash
	err := address.LoadString(hexAddress)
	if err != nil {
		cli.DieWithError("failed to parse hex-encoded TFT address", err)
		return
	}
	// load amount (in TFT)
	currencyConvertor := walletCmd.cli.CreateCurrencyConvertor()
	amount, err := currencyConvertor.ParseCoinString(strAmount)
	if err != nil {
		cli.DieWithError("failed to parse coin (TFT) string", err)
		return
	}
	// load ERC20 BlockID
	var blockID erc20types.ERC20Hash
	err = blockID.LoadString(hexBlockID)
	if err != nil {
		cli.DieWithError("failed to parse hex-encoded ERC20 BlockID", err)
		return
	}
	// load ERC20 TransactionID
	var transactionID erc20types.ERC20Hash
	err = transactionID.LoadString(hexTransctionID)
	if err != nil {
		cli.DieWithError("failed to parse hex-encoded ERC20 TransactionID", err)
		return
	}

	txFee := walletCmd.cli.Config.MinimumTransactionFee
	value := amount.Sub(txFee)

	// create the ERC20 CoinCreation Tx
	tx := erc20types.ERC20CoinCreationTransaction{
		Address:        address,
		Value:          value,
		TransactionFee: txFee,
		BlockID:        blockID,
		TransactionID:  transactionID,
	}
	rtx := tx.Transaction(walletCmd.txVersions.ERC20CoinCreation)

	// submit the Tx
	txID, err := walletCmd.txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the ERC20 CoinCreation Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	switch walletCmd.sendERC20FundsCfg.EncodingType {
	case cli.EncodingTypeHuman:
		fmt.Println("TransactionID:", txID)
	case cli.EncodingTypeJSON:
		err = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"transactionid": txID,
		})
		if err != nil {
			cli.DieWithError("failed to encode result", err)
		}
	}
}

func (walletCmd *walletCmd) sendERC20AddressRegistration(_ *cobra.Command, args []string) {
	var pubkey rivinetypes.PublicKey
	switch len(args) {
	case 0:
		// generate a new public key
		var err error
		pubkey, err = walletCmd.walletClient.NewPublicKey()
		if err != nil {
			cli.DieWithError("failed to generate new public key", err)
			return
		}
	case 1:
		err := pubkey.LoadString(args[0])
		if err != nil {
			cli.DieWithError("failed to parse stringified public key", err)
			return
		}
	default:
		cli.Die("only one pos. argument is allowed, an optional public key")
		return
	}

	// compute the hardcoded Tx fee
	regFee := walletCmd.cli.Config.CurrencyUnits.OneCoin.Mul64(erc20types.HardcodedERC20AddressRegistrationFeeOneCoinMultiplier)

	// create the ERC20 Address Registration Tx
	tx := erc20types.ERC20AddressRegistrationTransaction{
		PublicKey:       pubkey,
		Signature:       nil, // will be signed later by the daemon
		RegistrationFee: regFee,
		TransactionFee:  walletCmd.cli.Config.MinimumTransactionFee,
	}
	// fund the coin inputs
	var err error
	tx.CoinInputs, tx.RefundCoinOutput, err = walletCmd.walletClient.FundCoins(tx.TransactionFee.Add(regFee), nil, false)
	if err != nil {
		cli.DieWithError("failed to fund the ERC20 Address Registration Tx", err)
		return
	}

	// sign the Tx
	rtx := tx.Transaction(walletCmd.txVersions.ERC20AddressRegistration)
	err = walletCmd.walletClient.GreedySignTx(&rtx)
	if err != nil {
		cli.DieWithError("failed to sign the ERC20 Address Registration Tx", err)
		return
	}

	// submit the Tx
	txID, err := walletCmd.txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the Address Registration Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	tftAddr, err := rivinetypes.NewPubKeyUnlockHash(pubkey)
	if err != nil {
		cli.DieWithError("failed to create public key unlock hash", err)
		return
	}
	switch walletCmd.sendERC20AddressRegistrationCfg.EncodingType {
	case cli.EncodingTypeHuman:
		erc20Addr, err := erc20types.ERC20AddressFromUnlockHash(tftAddr)
		if err != nil {
			cli.DieWithError("failed to create ERC20 address from unlockhash", err)
			return
		}
		fmt.Println("Transaction ID:", txID)
		fmt.Println("TFT address:   ", tftAddr)
		fmt.Println("ERC20 address: ", erc20Addr)
	case cli.EncodingTypeJSON:
		erc20Addr, err := erc20types.ERC20AddressFromUnlockHash(tftAddr)
		if err != nil {
			cli.DieWithError("failed to create ERC20 address from unlockhash", err)
			return
		}
		err = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"transactionid": txID,
			"tft_address":   tftAddr,
			"erc20_address": erc20Addr,
		})
		if err != nil {
			cli.DieWithError("failed to encode result", err)
		}
	}
}

func (walletCmd *walletCmd) listERC20AddressRegistrations(_ *cobra.Command, args []string) {
	// fetch all known addresses
	addrs := new(rivineapi.WalletAddressesGET)
	err := walletCmd.cli.GetWithResponse("/wallet/addresses", addrs)
	if err != nil {
		cli.DieWithError("Failed to fetch addresses:", err)
	}

	erc20Addresses := []erc20types.ERC20Address{}
	// for every address check if it has a known link
	for _, addr := range addrs.Addresses {
		erc20Addr, exists, err := walletCmd.erc20Client.GetERC20AddressForTFTAddress(addr)
		if err != nil {
			cli.DieWithError("Failed to verify if address is linked to a known erc20 withdraw address: ", err)
		}
		if !exists {
			continue
		}
		erc20Addresses = append(erc20Addresses, erc20Addr)
	}

	switch walletCmd.listERC20AddressRegistrationsCfg.EncodingType {
	case cli.EncodingTypeHuman:
		if len(erc20Addresses) > 0 {
			fmt.Println("Known ERC20 withdrawal addresses:")
			fmt.Println()
			for _, address := range erc20Addresses {
				fmt.Println("\t", address)
			}
		} else {
			fmt.Println("No known ERC20 withdrawal addresses")
		}
	case cli.EncodingTypeJSON:
		err = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"erc20_withdraw_addresses": erc20Addresses,
		})
		if err != nil {
			cli.DieWithError("failed to encode result", err)
		}
	}
}
