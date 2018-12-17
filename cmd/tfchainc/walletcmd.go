package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/threefoldfoundation/tfchain/cmd/tfchainc/internal"

	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/threefoldtech/rivine/pkg/cli"
	rivinecli "github.com/threefoldtech/rivine/pkg/client"
	rivinetypes "github.com/threefoldtech/rivine/types"

	"github.com/spf13/cobra"
)

func createWalletSubCmds(client *rivinecli.CommandLineClient) {
	walletSubCmds := &walletSubCmds{cli: client}

	// define commands
	var (
		createMinterDefinitionTxCmd = &cobra.Command{
			Use:   "minterdefinitiontransaction <dest>|<rawCondition>",
			Short: "Create a new minter definition transaction",
			Long: `Create a new minter definition transaction using the given mint condition.
The mint condition is used to overwrite the current globally defined mint condition,
and can be given as a raw output condition (or address, which resolves to a singlesignature condition).

The returned (raw) MinterDefinitionTransaction still has to be signed, prior to sending.
	`,
			Run: walletSubCmds.createMinterDefinitionTxCmd,
		}
		createCoinCreationTxCmd = &cobra.Command{
			Use:   "coincreationtransaction <dest>|<rawCondition> <amount> [<dest>|<rawCondition> <amount>]...",
			Short: "Create a new coin creation transaction",
			Long: `Create a new coin creation transaction using the given outputs.
The outputs can be given as a pair of value and a raw output condition (or
address, which resolves to a singlesignature condition).

Amounts have to be given expressed in the OneCoin unit, and without the unit of currency.
Decimals are possible and have to be defined using the decimal point.

The Minimum Miner Fee will be added on top of the total given amount automatically.

The returned (raw) CoinCreationTransaction still has to be signed, prior to sending.
	`,
			Run: walletSubCmds.createCoinCreationTxCmd,
		}

		sendBotRegistrationTxCmd = &cobra.Command{
			Use:   "botregistration",
			Short: "Create, sign and send a new 3bot registration transaction",
			Long: `Create, sign and send a new 3bot registration transaction, prepaying 1 month by default.
The coin inputs are funded and signed using the wallet of this daemon.
By default a public key is generated from this wallet's primary seed as well,
however, it is also allowed for you to give a public key that is already loaded in this wallet,
for the creation of the 3bot.

Addresses and names are added as flags, and at least one of both is required.
Multiple addresses and names are allowed as well, of course.

Should you want to prepay more than 1 month, this has to be specified as a flag as well.
One might want to do this, as the ThreefoldFoundation gives 30% discount for 12+ (bot) months,
and 50% discount for 24 (bot) months (the maximum).

All fees are automatically added.

If this command returns without errors, the Tx is signed and sent,
and you'll receive the TxID and PublicKey which will allow you to look it up in an explorer.
The public key is to be used to get to know the unique ID assigned to your registered bot (if succesfull).
`,
			Run: rivinecli.Wrap(walletSubCmds.sendBotRegistrationTxCmd),
		}

		sendBotRecordUpdateTxCmd = &cobra.Command{
			Use:   "botupdate (id|publickey)",
			Short: "Create, sign and send a 3bot record update transaction",
			Long: `Create, sign and send a 3bot record update transaction, updating an existing 3bot.
The coin inputs are funded and signed using the wallet of this daemon.
The Public key linked to the 3bot has to be loaded into the wallet in order to be able to sign.

Addresses and names to be removed/added are defined as flags, and at least one
update is required (defining NrOfMonths to add (and pay) to the 3bot record counts as an update as well).

> NOTE: a name can only be removed if owned (which implies the 3bot has to be active at the point of the update).

Should you want to prepay more than 1 month at once, this is possible and
the ThreefoldFoundation gives 30% discount for 12+ (bot) months,
and 50% discount for 24 (bot) months (the maximum).

All fees are automatically added.

If this command returns without errors, the Tx is signed and sent,
and you'll receive the TxID which will allow you to look it up in an explorer.
`,
			Run: rivinecli.Wrap(walletSubCmds.sendBotRecordUpdateTxCmd),
		}

		createBotNameTransferTxCmd = &cobra.Command{
			Use:   "botnametransfer (id|publickey) (id|publickey) names...",
			Args:  cobra.MinimumNArgs(3),
			Short: "Create and optionally sign a 3bot name transfer transaction",
			Long: `Create and optionally sign a 3bot name transfer transaction, involving two active 3bots.
The coin inputs are funded and signed using the wallet of this daemon.
The Public key linked to the 3bot has to be loaded into the wallet in order to be able to sign.

The first positional argument identifies the sender, nad the second positional argument identifies the receiver.
All other positional arguments (at least one more is required) define the names to be transfered.
At least one name has to be transferred.

All fees are automatically added.

If this command returns without errors, the Tx (optionally signed)
is printed to the STDOUT.
`,
			Run: walletSubCmds.createBotNameTransferTxCmd,
		}

		sendERC20FundsCmd = &cobra.Command{
			Use:   "erc20funds erc20_address amount",
			Short: "Convert TFT to ERC20 funds and send those to an ERC20 adddress (minus fees)",
			Run:   rivinecli.Wrap(walletSubCmds.sendERC20Funds),
		}

		sendERC20FundsClaimCmd = &cobra.Command{
			Use:   "erc20fundsclaim tft_address amount erc20_txid",
			Short: "Convert ERC20 funds to TFT",
			Run:   rivinecli.Wrap(walletSubCmds.sendERC20FundsClaim),
		}

		sendERC20AddressRegistrationCmd = &cobra.Command{
			Use:   "erc20address [public_key]",
			Short: "Register an ERC20 address linked to a TFT Public Key",
			Long: `Register an ERC20 address linked to a TFT Public Key

If no Public Key is given, a new one will be generated
using the unlocked wallet from the tfchain daemon.
`,
			Run: walletSubCmds.sendERC20AddressRegistration,
		}
	)

	// add commands as wallet sub commands
	client.WalletCmd.RootCmdCreate.AddCommand(
		createMinterDefinitionTxCmd,
		createCoinCreationTxCmd,
		createBotNameTransferTxCmd,
	)
	client.WalletCmd.RootCmdSend.AddCommand(
		sendBotRegistrationTxCmd,
		sendBotRecordUpdateTxCmd,
		sendERC20FundsCmd,
		sendERC20FundsClaimCmd,
		sendERC20AddressRegistrationCmd,
	)

	// register flags

	createMinterDefinitionTxCmd.Flags().Var(
		&walletSubCmds.minterDefinitionTxCfg.Data, "data",
		"optionally add a description to describe the reasons of transfer of minting power, added as arbitrary data")

	createCoinCreationTxCmd.Flags().Var(
		&walletSubCmds.coinCreationTxCfg.Data, "data",
		"optionally add a description to describe the origins of the coin creation, added as arbitrary data")

	internal.NetworkAddressArrayFlagVar(
		sendBotRegistrationTxCmd.Flags(),
		&walletSubCmds.sendBotRegistrationTxCfg.Addresses,
		"address",
		"add one or multiple addresses, each address defined as seperate flag arguments",
	)
	internal.BotNameArrayFlagVar(
		sendBotRegistrationTxCmd.Flags(),
		&walletSubCmds.sendBotRegistrationTxCfg.Names,
		"name",
		"add one or multiple names, each name defined as seperate flag arguments",
	)
	sendBotRegistrationTxCmd.Flags().Uint8VarP(
		&walletSubCmds.sendBotRegistrationTxCfg.NrOfMonths, "months", "m", 1,
		"the amount of months to prepay, required to be in the inclusive interval [1, 24]")
	internal.PublicKeyFlagVar(
		sendBotRegistrationTxCmd.Flags(),
		&walletSubCmds.sendBotRegistrationTxCfg.PublicKey,
		"public-key",
		"define a public key to use (of which the private key is loaded in this daemon's wallet)",
	)
	sendBotRegistrationTxCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletSubCmds.sendBotRegistrationTxCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	internal.NetworkAddressArrayFlagVar(
		sendBotRecordUpdateTxCmd.Flags(),
		&walletSubCmds.sendBotRecordUpdateTxCfg.AddressesToAdd,
		"add-address",
		"add one or multiple addresses, each address defined as seperate flag arguments",
	)
	internal.NetworkAddressArrayFlagVar(
		sendBotRecordUpdateTxCmd.Flags(),
		&walletSubCmds.sendBotRecordUpdateTxCfg.AddressesToRemove,
		"remove-address",
		"remove one or multiple addresses, each address defined as seperate flag arguments",
	)
	internal.BotNameArrayFlagVar(
		sendBotRecordUpdateTxCmd.Flags(),
		&walletSubCmds.sendBotRecordUpdateTxCfg.NamesToAdd,
		"add-name",
		"add one or multiple names, each name defined as seperate flag arguments",
	)
	internal.BotNameArrayFlagVar(
		sendBotRecordUpdateTxCmd.Flags(),
		&walletSubCmds.sendBotRecordUpdateTxCfg.NamesToRemove,
		"remove-name",
		"remove one or multiple names owned, each name defined as seperate flag arguments",
	)
	sendBotRecordUpdateTxCmd.Flags().Uint8VarP(
		&walletSubCmds.sendBotRecordUpdateTxCfg.NrOfMonthsToAdd, "add-months", "m", 0,
		"the amount of months to add and pay, required to be in the inclusive interval [0, 24]")
	sendBotRecordUpdateTxCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletSubCmds.sendBotRecordUpdateTxCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	createBotNameTransferTxCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletSubCmds.createBotNameTransferTxCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))
	createBotNameTransferTxCmd.Flags().BoolVar(
		&walletSubCmds.createBotNameTransferTxCfg.Sign, "sign", false,
		"optionally sign the transaction (as sender/receiver) prior to printing it")

	sendERC20FundsCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletSubCmds.sendERC20FundsCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	sendERC20FundsClaimCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletSubCmds.sendERC20FundsClaimCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	sendERC20AddressRegistrationCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletSubCmds.sendERC20AddressRegistrationCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))
}

type walletSubCmds struct {
	cli                   *rivinecli.CommandLineClient
	minterDefinitionTxCfg struct {
		Data cli.ArbitraryDataFlag
	}
	coinCreationTxCfg struct {
		Data cli.ArbitraryDataFlag
	}

	sendBotRegistrationTxCfg struct {
		Addresses    []types.NetworkAddress
		Names        []types.BotName
		NrOfMonths   uint8
		PublicKey    rivinetypes.PublicKey
		EncodingType cli.EncodingType
	}

	sendBotRecordUpdateTxCfg struct {
		AddressesToAdd    []types.NetworkAddress
		AddressesToRemove []types.NetworkAddress
		NamesToAdd        []types.BotName
		NamesToRemove     []types.BotName
		NrOfMonthsToAdd   uint8
		EncodingType      cli.EncodingType
	}

	createBotNameTransferTxCfg struct {
		EncodingType cli.EncodingType
		Sign         bool
	}

	sendERC20FundsCfg struct {
		EncodingType cli.EncodingType
	}
	sendERC20FundsClaimCfg struct {
		EncodingType cli.EncodingType
	}
	sendERC20AddressRegistrationCfg struct {
		EncodingType cli.EncodingType
	}
}

func (walletSubCmds *walletSubCmds) createMinterDefinitionTxCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.UsageFunc()
		cli.Die("Invalid amount of arguments. One argume has to be given: <dest>|<rawCondition>")
	}

	// create a minter definition tx with a random nonce and the minimum required miner fee
	tx := types.MinterDefinitionTransaction{
		Nonce:     types.RandomTransactionNonce(),
		MinerFees: []rivinetypes.Currency{walletSubCmds.cli.Config.MinimumTransactionFee},
	}

	// parse the given mint condition
	var err error
	tx.MintCondition, err = parseConditionString(args[0])
	if err != nil {
		cmd.UsageFunc()(cmd)
		cli.Die(err)
	}

	// if data is given, use it as arbitrary data
	tx.ArbitraryData.Data = walletSubCmds.minterDefinitionTxCfg.Data.Data
	tx.ArbitraryData.Type = walletSubCmds.minterDefinitionTxCfg.Data.DataType

	// encode the transaction as a JSON-encoded string and print it to the STDOUT
	json.NewEncoder(os.Stdout).Encode(tx.Transaction())
}

func (walletSubCmds *walletSubCmds) createCoinCreationTxCmd(cmd *cobra.Command, args []string) {
	currencyConvertor := walletSubCmds.cli.CreateCurrencyConvertor()

	// Check that the remaining args are condition + value pairs
	if len(args)%2 != 0 {
		cmd.UsageFunc()
		cli.Die("Invalid arguments. Arguments must be of the form <dest>|<rawCondition> <amount> [<dest>|<rawCondition> <amount>]...")
	}

	// parse the remainder as output coditions and values
	pairs, err := parsePairedOutputs(args, currencyConvertor.ParseCoinString)
	if err != nil {
		cmd.UsageFunc()(cmd)
		cli.Die(err)
	}

	tx := types.CoinCreationTransaction{
		Nonce:     types.RandomTransactionNonce(),
		MinerFees: []rivinetypes.Currency{walletSubCmds.cli.Config.MinimumTransactionFee},
	}
	tx.ArbitraryData.Data = walletSubCmds.coinCreationTxCfg.Data.Data
	tx.ArbitraryData.Type = walletSubCmds.coinCreationTxCfg.Data.DataType
	for _, pair := range pairs {
		tx.CoinOutputs = append(tx.CoinOutputs, rivinetypes.CoinOutput{
			Value:     pair.Value,
			Condition: pair.Condition,
		})
	}
	json.NewEncoder(os.Stdout).Encode(tx.Transaction())
}

func (walletSubCmds *walletSubCmds) sendBotRegistrationTxCmd() {
	// validate the flags
	if len(walletSubCmds.sendBotRegistrationTxCfg.Addresses) == 0 && len(walletSubCmds.sendBotRegistrationTxCfg.Names) == 0 {
		cli.Die("the registration of a 3bot requires at least one name or address to be defined")
		return
	}
	if walletSubCmds.sendBotRegistrationTxCfg.NrOfMonths == 0 || walletSubCmds.sendBotRegistrationTxCfg.NrOfMonths > 24 {
		cli.Die("the number of (prepaid) (bot) months has to be in the inclusive interval [1,24]")
		return
	}

	// start the registration process
	walletClient := internal.NewWalletClient(walletSubCmds.cli)

	var err error
	pk := walletSubCmds.sendBotRegistrationTxCfg.PublicKey
	if pk.Algorithm == 0 && len(pk.Key) == 0 {
		pk, err = walletClient.NewPublicKey()
		if err != nil {
			cli.DieWithError("failed to generate new public key", err)
			return
		}
	}

	// create the registration Tx
	tx := types.BotRegistrationTransaction{
		Addresses:      walletSubCmds.sendBotRegistrationTxCfg.Addresses,
		Names:          walletSubCmds.sendBotRegistrationTxCfg.Names,
		NrOfMonths:     walletSubCmds.sendBotRegistrationTxCfg.NrOfMonths,
		TransactionFee: walletSubCmds.cli.Config.MinimumTransactionFee,
		Identification: types.PublicKeySignaturePair{
			PublicKey: pk,
		},
	}
	// compute the additional (bot) fee, such that we can fund it all
	fee := tx.RequiredBotFee(walletSubCmds.cli.Config.CurrencyUnits.OneCoin)
	// fund the coin inputs
	tx.CoinInputs, tx.RefundCoinOutput, err = walletClient.FundCoins(fee.Add(walletSubCmds.cli.Config.MinimumTransactionFee))
	if err != nil {
		cli.DieWithError("failed to fund the bot registration Tx", err)
		return
	}

	// sign the Tx
	rtx := tx.Transaction(
		walletSubCmds.cli.Config.CurrencyUnits.OneCoin,
		internal.GetFoundationPoolCondition(walletSubCmds.cli.Config.NetworkName))
	err = walletClient.GreedySignTx(&rtx)
	if err != nil {
		cli.DieWithError("failed to sign the bot registration Tx", err)
		return
	}

	// submit the Tx
	txPoolClient := internal.NewTransactionPoolClient(walletSubCmds.cli)
	txID, err := txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the bot registration Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletSubCmds.sendBotRegistrationTxCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	}
	err = encode(map[string]interface{}{
		"publickey":     pk,
		"transactionid": txID,
	})
	if err != nil {
		cli.DieWithError("failed to encode result", err)
	}
}

func (walletSubCmds *walletSubCmds) sendBotRecordUpdateTxCmd(str string) {
	id, err := walletSubCmds.botIDFromPosArgStr(str)
	if err != nil {
		cli.DieWithError("failed to parse/fetch unique ID", err)
		return
	}

	// start the record update process
	walletClient := internal.NewWalletClient(walletSubCmds.cli)

	// create the record update Tx
	tx := types.BotRecordUpdateTransaction{
		Identifier: id,
		Addresses: types.BotRecordAddressUpdate{
			Add:    walletSubCmds.sendBotRecordUpdateTxCfg.AddressesToAdd,
			Remove: walletSubCmds.sendBotRecordUpdateTxCfg.AddressesToRemove,
		},
		Names: types.BotRecordNameUpdate{
			Add:    walletSubCmds.sendBotRecordUpdateTxCfg.NamesToAdd,
			Remove: walletSubCmds.sendBotRecordUpdateTxCfg.NamesToRemove,
		},
		NrOfMonths:     walletSubCmds.sendBotRecordUpdateTxCfg.NrOfMonthsToAdd,
		TransactionFee: walletSubCmds.cli.Config.MinimumTransactionFee,
	}
	// compute the additional (bot) fee, such that we can fund it all
	fee := tx.RequiredBotFee(walletSubCmds.cli.Config.CurrencyUnits.OneCoin)
	// fund the coin inputs
	tx.CoinInputs, tx.RefundCoinOutput, err = walletClient.FundCoins(fee.Add(walletSubCmds.cli.Config.MinimumTransactionFee))
	if err != nil {
		cli.DieWithError("failed to fund the bot record update Tx", err)
		return
	}

	// sign the Tx
	rtx := tx.Transaction(
		walletSubCmds.cli.Config.CurrencyUnits.OneCoin,
		internal.GetFoundationPoolCondition(walletSubCmds.cli.Config.NetworkName))
	err = walletClient.GreedySignTx(&rtx)
	if err != nil {
		cli.DieWithError("failed to sign the bot record update Tx", err)
		return
	}

	// submit the Tx
	txPoolClient := internal.NewTransactionPoolClient(walletSubCmds.cli)
	txID, err := txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the bot record update Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletSubCmds.sendBotRecordUpdateTxCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	}
	err = encode(map[string]interface{}{
		"transactionid": txID,
	})
	if err != nil {
		cli.DieWithError("failed to encode result", err)
	}
}

// create botnametransfer (publickey|id) (publickey|id) names...
// arguments in order: sender, receiver and a slice of names (at least one name is required),
// hence this command requires a minimum of 3 arguments
func (walletSubCmds *walletSubCmds) createBotNameTransferTxCmd(cmd *cobra.Command, args []string) {
	senderID, err := walletSubCmds.botIDFromPosArgStr(args[0])
	if err != nil {
		cli.DieWithError("failed to parse/fetch unique (sender bot) ID", err)
		return
	}
	receiverID, err := walletSubCmds.botIDFromPosArgStr(args[1])
	if err != nil {
		cli.DieWithError("failed to parse/fetch unique (receiver bot) ID", err)
		return
	}

	// start the record update process
	walletClient := internal.NewWalletClient(walletSubCmds.cli)

	names := make([]types.BotName, len(args[2:]))
	for idx, str := range args[2:] {
		err = names[idx].LoadString(str)
		if err != nil {
			cli.DieWithError("failed to parse (pos arg) bot name #"+strconv.Itoa(idx+1), err)
			return
		}
	}

	// create the bot name transfer Tx
	tx := types.BotNameTransferTransaction{
		Sender: types.BotIdentifierSignaturePair{
			Identifier: senderID,
		},
		Receiver: types.BotIdentifierSignaturePair{
			Identifier: receiverID,
		},
		Names:          names,
		TransactionFee: walletSubCmds.cli.Config.MinimumTransactionFee,
	}
	// compute the additional (bot) fee, such that we can fund it all
	fee := tx.RequiredBotFee(walletSubCmds.cli.Config.CurrencyUnits.OneCoin)
	// fund the coin inputs
	tx.CoinInputs, tx.RefundCoinOutput, err = walletClient.FundCoins(fee.Add(walletSubCmds.cli.Config.MinimumTransactionFee))
	if err != nil {
		cli.DieWithError("failed to fund the bot name transfer Tx", err)
		return
	}

	rtx := tx.Transaction(
		walletSubCmds.cli.Config.CurrencyUnits.OneCoin,
		internal.GetFoundationPoolCondition(walletSubCmds.cli.Config.NetworkName))

	if walletSubCmds.createBotNameTransferTxCfg.Sign {
		// optionally sign the Tx
		err = walletClient.GreedySignTx(&rtx)
		if err != nil {
			cli.DieWithError("failed to sign the bot name transfer Tx", err)
			return
		}
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletSubCmds.createBotNameTransferTxCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	}
	err = encode(rtx)
	if err != nil {
		cli.DieWithError("failed to encode result", err)
	}
}

func (walletSubCmds *walletSubCmds) sendERC20Funds(hexAddress, strAmount string) {
	// load ERC20 address
	var address types.ERC20Address
	err := address.LoadString(hexAddress)
	if err != nil {
		cli.DieWithError("failed to parse hex-encoded ERC20 address", err)
		return
	}
	// load amount (in TFT)
	currencyConvertor := walletSubCmds.cli.CreateCurrencyConvertor()
	amount, err := currencyConvertor.ParseCoinString(strAmount)
	if err != nil {
		cli.DieWithError("failed to parse coin (TFT) string", err)
		return
	}

	// start the ER20 fund convert process
	walletClient := internal.NewWalletClient(walletSubCmds.cli)

	// create the ERC20 Convert Tx
	tx := types.ERC20ConvertTransaction{
		Address:        address,
		Value:          amount,
		TransactionFee: walletSubCmds.cli.Config.MinimumTransactionFee,
	}
	// fund the coin inputs
	tx.CoinInputs, tx.RefundCoinOutput, err = walletClient.FundCoins(tx.TransactionFee.Add(tx.Value))
	if err != nil {
		cli.DieWithError("failed to fund the ERC20 Convert Txd", err)
		return
	}

	// sign the Tx
	rtx := tx.Transaction()
	err = walletClient.GreedySignTx(&rtx)
	if err != nil {
		cli.DieWithError("failed to sign the ERC20 Convert Tx", err)
		return
	}

	// submit the Tx
	txPoolClient := internal.NewTransactionPoolClient(walletSubCmds.cli)
	txID, err := txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the ERC20 Convert Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletSubCmds.sendERC20FundsCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	}
	err = encode(map[string]interface{}{
		"transactionid": txID,
	})
	if err != nil {
		cli.DieWithError("failed to encode result", err)
	}
}

func (walletSubCmds *walletSubCmds) sendERC20FundsClaim(hexAddress, strAmount, hexTransctionID string) {
	// load TFT address
	var address rivinetypes.UnlockHash
	err := address.LoadString(hexAddress)
	if err != nil {
		cli.DieWithError("failed to parse hex-encoded TFT address", err)
		return
	}
	// load amount (in TFT)
	currencyConvertor := walletSubCmds.cli.CreateCurrencyConvertor()
	amount, err := currencyConvertor.ParseCoinString(strAmount)
	if err != nil {
		cli.DieWithError("failed to parse coin (TFT) string", err)
		return
	}
	// load ERC20 TransactionID
	var transactionID types.ERC20TransactionID
	err = transactionID.LoadString(hexTransctionID)
	if err != nil {
		cli.DieWithError("failed to parse hex-encoded ERC20 TransactionID", err)
		return
	}

	// create the ERC20 CoinCreation Tx
	tx := types.ERC20CoinCreationTransaction{
		Address:        address,
		Value:          amount,
		TransactionFee: walletSubCmds.cli.Config.MinimumTransactionFee,
		TransactionID:  transactionID,
	}
	rtx := tx.Transaction()

	// submit the Tx
	txPoolClient := internal.NewTransactionPoolClient(walletSubCmds.cli)
	txID, err := txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the ERC20 CoinCreation Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletSubCmds.sendERC20FundsCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	}
	err = encode(map[string]interface{}{
		"transactionid": txID,
	})
	if err != nil {
		cli.DieWithError("failed to encode result", err)
	}
}

func (walletSubCmds *walletSubCmds) sendERC20AddressRegistration(_ *cobra.Command, args []string) {
	var pubkey rivinetypes.PublicKey

	// required to fund as well as to create a public key if needed
	walletClient := internal.NewWalletClient(walletSubCmds.cli)

	switch len(args) {
	case 0:
		// generate a new public key
		var err error
		pubkey, err = walletClient.NewPublicKey()
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
	regFee := walletSubCmds.cli.Config.CurrencyUnits.OneCoin.Mul64(types.HardcodedERC20AddressRegistrationFeeOneCoinMultiplier)

	// create the ERC20 Address Registration Tx
	tx := types.ERC20AddressRegistrationTransaction{
		PublicKey:       pubkey,
		Signature:       nil, // will be signed later by the daemon
		RegistrationFee: regFee,
		TransactionFee:  walletSubCmds.cli.Config.MinimumTransactionFee,
	}
	// fund the coin inputs
	var err error
	tx.CoinInputs, tx.RefundCoinOutput, err = walletClient.FundCoins(tx.TransactionFee.Add(regFee))
	if err != nil {
		cli.DieWithError("failed to fund the ERC20 Address Registration Tx", err)
		return
	}

	// sign the Tx
	rtx := tx.Transaction()
	err = walletClient.GreedySignTx(&rtx)
	if err != nil {
		cli.DieWithError("failed to sign the ERC20 Address Registration Tx", err)
		return
	}

	// submit the Tx
	txPoolClient := internal.NewTransactionPoolClient(walletSubCmds.cli)
	txID, err := txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the Address Registration Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletSubCmds.sendERC20FundsCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	}
	err = encode(map[string]interface{}{
		"transactionid": txID,
	})
	if err != nil {
		cli.DieWithError("failed to encode result", err)
	}
}

func (walletSubCmds *walletSubCmds) botIDFromPosArgStr(str string) (types.BotID, error) {
	if len(str) < 16 {
		// assume bot ID if the less than 16, seems to short for a public key,
		// so simply return it (as well as the possible parsing error for assuming wrong)
		var botID types.BotID
		err := botID.LoadString(str)
		return botID, err
	}

	// assume a public key was meant,
	// so we need to get the (bot) record in order to know the (unique) ID
	var pk rivinetypes.PublicKey
	err := pk.LoadString(str)
	if err != nil {
		return 0, err
	}
	txDB := internal.NewTransactionDBConsensusClient(walletSubCmds.cli)
	record, err := txDB.GetRecordForKey(pk)
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}

type (
	// parseCurrencyString takes the string representation of a currency value
	parseCurrencyString func(string) (rivinetypes.Currency, error)

	outputPair struct {
		Condition rivinetypes.UnlockConditionProxy
		Value     rivinetypes.Currency
	}
)

func parsePairedOutputs(args []string, parseCurrency parseCurrencyString) (pairs []outputPair, err error) {
	argn := len(args)
	if argn < 2 {
		err = errors.New("not enough arguments, at least 2 required")
		return
	}
	if argn%2 != 0 {
		err = errors.New("arguments have to be given in pairs of '<dest>|<rawCondition>'+'<value>'")
		return
	}

	for i := 0; i < argn; i += 2 {
		// parse value first, as it's the one without any possibility of ambiguity
		var pair outputPair
		pair.Value, err = parseCurrency(args[i+1])
		if err != nil {
			err = fmt.Errorf("failed to parse amount/value for output #%d: %v", i/2, err)
			return
		}

		// parse condition second
		pair.Condition, err = parseConditionString(args[i])
		if err != nil {
			err = fmt.Errorf("failed to parse condition for output #%d: %v", i/2, err)
			return
		}

		// append succesfully parsed pair
		pairs = append(pairs, pair)
	}
	return
}

// try to parse the string first as an unlock hash,
// if that fails parse it as a
func parseConditionString(str string) (condition rivinetypes.UnlockConditionProxy, err error) {
	// try to parse it as an unlock hash
	var uh rivinetypes.UnlockHash
	err = uh.LoadString(str)
	if err == nil {
		// parsing as an unlock hash was succesfull, store the pair and continue to the next pair
		condition = rivinetypes.NewCondition(rivinetypes.NewUnlockHashCondition(uh))
		return
	}

	// try to parse it as a JSON-encoded unlock condition
	err = condition.UnmarshalJSON([]byte(str))
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, fmt.Errorf(
			"condition has to be UnlockHash or JSON-encoded UnlockCondition, output %q is neither", str)
	}
	return
}
