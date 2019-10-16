package client

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	tbtypes "github.com/threefoldfoundation/tfchain/extensions/threebot/types"

	"github.com/threefoldtech/rivine/pkg/cli"
	"github.com/threefoldtech/rivine/pkg/client"
	rivinecli "github.com/threefoldtech/rivine/pkg/client"
	rivinetypes "github.com/threefoldtech/rivine/types"

	"github.com/spf13/cobra"
)

// CreateWalletCmds creates the threebot wallet root command as well as its transaction creation sub commands.
func CreateWalletCmds(ccli *client.CommandLineClient) error {
	bc, err := client.NewLazyBaseClientFromCommandLineClient(ccli)
	if err != nil {
		return err
	}

	walletCmd := &walletCmd{
		cli:          ccli,
		walletClient: rivinecli.NewWalletClient(bc),
		txPoolClient: rivinecli.NewTransactionPoolClient(bc),
		tbClient:     NewPluginExplorerClient(bc),
	}

	// define commands
	var (
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
			Run: rivinecli.Wrap(walletCmd.sendBotRegistrationTxCmd),
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
			Run: rivinecli.Wrap(walletCmd.sendBotRecordUpdateTxCmd),
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
			Run: walletCmd.createBotNameTransferTxCmd,
		}
	)

	// add commands as wallet sub commands
	ccli.WalletCmd.RootCmdCreate.AddCommand(
		createBotNameTransferTxCmd,
	)
	ccli.WalletCmd.RootCmdSend.AddCommand(
		sendBotRegistrationTxCmd,
		sendBotRecordUpdateTxCmd,
	)

	// register flags
	NetworkAddressArrayFlagVar(
		sendBotRegistrationTxCmd.Flags(),
		&walletCmd.sendBotRegistrationTxCfg.Addresses,
		"address",
		"add one or multiple addresses, each address defined as seperate flag arguments",
	)
	BotNameArrayFlagVar(
		sendBotRegistrationTxCmd.Flags(),
		&walletCmd.sendBotRegistrationTxCfg.Names,
		"name",
		"add one or multiple names, each name defined as seperate flag arguments",
	)
	sendBotRegistrationTxCmd.Flags().Uint8VarP(
		&walletCmd.sendBotRegistrationTxCfg.NrOfMonths, "months", "m", 1,
		"the amount of months to prepay, required to be in the inclusive interval [1, 24]")
	PublicKeyFlagVar(
		sendBotRegistrationTxCmd.Flags(),
		&walletCmd.sendBotRegistrationTxCfg.PublicKey,
		"public-key",
		"define a public key to use (of which the private key is loaded in this daemon's wallet)",
	)
	sendBotRegistrationTxCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletCmd.sendBotRegistrationTxCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	NetworkAddressArrayFlagVar(
		sendBotRecordUpdateTxCmd.Flags(),
		&walletCmd.sendBotRecordUpdateTxCfg.AddressesToAdd,
		"add-address",
		"add one or multiple addresses, each address defined as seperate flag arguments",
	)
	NetworkAddressArrayFlagVar(
		sendBotRecordUpdateTxCmd.Flags(),
		&walletCmd.sendBotRecordUpdateTxCfg.AddressesToRemove,
		"remove-address",
		"remove one or multiple addresses, each address defined as seperate flag arguments",
	)
	BotNameArrayFlagVar(
		sendBotRecordUpdateTxCmd.Flags(),
		&walletCmd.sendBotRecordUpdateTxCfg.NamesToAdd,
		"add-name",
		"add one or multiple names, each name defined as seperate flag arguments",
	)
	BotNameArrayFlagVar(
		sendBotRecordUpdateTxCmd.Flags(),
		&walletCmd.sendBotRecordUpdateTxCfg.NamesToRemove,
		"remove-name",
		"remove one or multiple names owned, each name defined as seperate flag arguments",
	)
	sendBotRecordUpdateTxCmd.Flags().Uint8VarP(
		&walletCmd.sendBotRecordUpdateTxCfg.NrOfMonthsToAdd, "add-months", "m", 0,
		"the amount of months to add and pay, required to be in the inclusive interval [0, 24]")
	sendBotRecordUpdateTxCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletCmd.sendBotRecordUpdateTxCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	createBotNameTransferTxCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &walletCmd.createBotNameTransferTxCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))
	createBotNameTransferTxCmd.Flags().BoolVar(
		&walletCmd.createBotNameTransferTxCfg.Sign, "sign", false,
		"optionally sign the transaction (as sender/receiver) prior to printing it")

	return nil
}

type walletCmd struct {
	cli          *rivinecli.CommandLineClient
	walletClient *rivinecli.WalletClient
	txPoolClient *rivinecli.TransactionPoolClient
	tbClient     *PluginClient

	sendBotRegistrationTxCfg struct {
		Addresses    []tbtypes.NetworkAddress
		Names        []tbtypes.BotName
		NrOfMonths   uint8
		PublicKey    rivinetypes.PublicKey
		EncodingType cli.EncodingType
	}

	sendBotRecordUpdateTxCfg struct {
		AddressesToAdd    []tbtypes.NetworkAddress
		AddressesToRemove []tbtypes.NetworkAddress
		NamesToAdd        []tbtypes.BotName
		NamesToRemove     []tbtypes.BotName
		NrOfMonthsToAdd   uint8
		EncodingType      cli.EncodingType
	}

	createBotNameTransferTxCfg struct {
		EncodingType cli.EncodingType
		Sign         bool
	}
}

func (walletCmd *walletCmd) sendBotRegistrationTxCmd() {
	// validate the flags
	if len(walletCmd.sendBotRegistrationTxCfg.Addresses) == 0 && len(walletCmd.sendBotRegistrationTxCfg.Names) == 0 {
		cli.Die("the registration of a 3bot requires at least one name or address to be defined")
		return
	}
	if walletCmd.sendBotRegistrationTxCfg.NrOfMonths == 0 || walletCmd.sendBotRegistrationTxCfg.NrOfMonths > 24 {
		cli.Die("the number of (prepaid) (bot) months has to be in the inclusive interval [1,24]")
		return
	}

	var err error
	pk := walletCmd.sendBotRegistrationTxCfg.PublicKey
	if pk.Algorithm == 0 && len(pk.Key) == 0 {
		pk, err = walletCmd.walletClient.NewPublicKey()
		if err != nil {
			cli.DieWithError("failed to generate new public key", err)
			return
		}
	}

	// create the registration Tx
	tx := tbtypes.BotRegistrationTransaction{
		Addresses:      walletCmd.sendBotRegistrationTxCfg.Addresses,
		Names:          walletCmd.sendBotRegistrationTxCfg.Names,
		NrOfMonths:     walletCmd.sendBotRegistrationTxCfg.NrOfMonths,
		TransactionFee: walletCmd.cli.Config.MinimumTransactionFee,
		Identification: tbtypes.PublicKeySignaturePair{
			PublicKey: pk,
		},
	}
	// compute the additional (bot) fee, such that we can fund it all
	fee := tx.RequiredBotFee(walletCmd.cli.Config.CurrencyUnits.OneCoin)
	// fund the coin inputs
	tx.CoinInputs, tx.RefundCoinOutput, err = walletCmd.walletClient.FundCoins(fee.Add(walletCmd.cli.Config.MinimumTransactionFee), nil, false)
	if err != nil {
		cli.DieWithError("failed to fund the bot registration Tx", err)
		return
	}

	// sign the Tx
	rtx := tx.Transaction(walletCmd.cli.Config.CurrencyUnits.OneCoin)
	err = walletCmd.walletClient.GreedySignTx(&rtx)
	if err != nil {
		cli.DieWithError("failed to sign the bot registration Tx", err)
		return
	}

	// submit the Tx
	txID, err := walletCmd.txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the bot registration Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletCmd.sendBotRegistrationTxCfg.EncodingType {
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

func (walletCmd *walletCmd) sendBotRecordUpdateTxCmd(str string) {
	id, err := walletCmd.botIDFromPosArgStr(str)
	if err != nil {
		cli.DieWithError("failed to parse/fetch unique ID", err)
		return
	}

	// create the record update Tx
	tx := tbtypes.BotRecordUpdateTransaction{
		Identifier: id,
		Addresses: tbtypes.BotRecordAddressUpdate{
			Add:    walletCmd.sendBotRecordUpdateTxCfg.AddressesToAdd,
			Remove: walletCmd.sendBotRecordUpdateTxCfg.AddressesToRemove,
		},
		Names: tbtypes.BotRecordNameUpdate{
			Add:    walletCmd.sendBotRecordUpdateTxCfg.NamesToAdd,
			Remove: walletCmd.sendBotRecordUpdateTxCfg.NamesToRemove,
		},
		NrOfMonths:     walletCmd.sendBotRecordUpdateTxCfg.NrOfMonthsToAdd,
		TransactionFee: walletCmd.cli.Config.MinimumTransactionFee,
	}
	// compute the additional (bot) fee, such that we can fund it all
	fee := tx.RequiredBotFee(walletCmd.cli.Config.CurrencyUnits.OneCoin)
	// fund the coin inputs
	tx.CoinInputs, tx.RefundCoinOutput, err = walletCmd.walletClient.FundCoins(fee.Add(walletCmd.cli.Config.MinimumTransactionFee), nil, false)
	if err != nil {
		cli.DieWithError("failed to fund the bot record update Tx", err)
		return
	}

	// sign the Tx
	rtx := tx.Transaction(walletCmd.cli.Config.CurrencyUnits.OneCoin)
	err = walletCmd.walletClient.GreedySignTx(&rtx)
	if err != nil {
		cli.DieWithError("failed to sign the bot record update Tx", err)
		return
	}

	// submit the Tx
	txID, err := walletCmd.txPoolClient.AddTransactiom(rtx)
	if err != nil {
		b, _ := json.Marshal(rtx)
		fmt.Fprintln(os.Stderr, "bad tx: "+string(b))
		cli.DieWithError("failed to submit the bot record update Tx to the Tx Pool", err)
		return
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletCmd.sendBotRecordUpdateTxCfg.EncodingType {
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
func (walletCmd *walletCmd) createBotNameTransferTxCmd(cmd *cobra.Command, args []string) {
	senderID, err := walletCmd.botIDFromPosArgStr(args[0])
	if err != nil {
		cli.DieWithError("failed to parse/fetch unique (sender bot) ID", err)
		return
	}
	receiverID, err := walletCmd.botIDFromPosArgStr(args[1])
	if err != nil {
		cli.DieWithError("failed to parse/fetch unique (receiver bot) ID", err)
		return
	}

	names := make([]tbtypes.BotName, len(args[2:]))
	for idx, str := range args[2:] {
		err = names[idx].LoadString(str)
		if err != nil {
			cli.DieWithError("failed to parse (pos arg) bot name #"+strconv.Itoa(idx+1), err)
			return
		}
	}

	// create the bot name transfer Tx
	tx := tbtypes.BotNameTransferTransaction{
		Sender: tbtypes.BotIdentifierSignaturePair{
			Identifier: senderID,
		},
		Receiver: tbtypes.BotIdentifierSignaturePair{
			Identifier: receiverID,
		},
		Names:          names,
		TransactionFee: walletCmd.cli.Config.MinimumTransactionFee,
	}
	// compute the additional (bot) fee, such that we can fund it all
	fee := tx.RequiredBotFee(walletCmd.cli.Config.CurrencyUnits.OneCoin)
	// fund the coin inputs
	tx.CoinInputs, tx.RefundCoinOutput, err = walletCmd.walletClient.FundCoins(fee.Add(walletCmd.cli.Config.MinimumTransactionFee), nil, false)
	if err != nil {
		cli.DieWithError("failed to fund the bot name transfer Tx", err)
		return
	}

	rtx := tx.Transaction(walletCmd.cli.Config.CurrencyUnits.OneCoin)

	if walletCmd.createBotNameTransferTxCfg.Sign {
		// optionally sign the Tx
		err = walletCmd.walletClient.GreedySignTx(&rtx)
		if err != nil {
			cli.DieWithError("failed to sign the bot name transfer Tx", err)
			return
		}
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch walletCmd.createBotNameTransferTxCfg.EncodingType {
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

func (walletCmd *walletCmd) botIDFromPosArgStr(str string) (tbtypes.BotID, error) {
	if len(str) < 16 {
		// assume bot ID if the less than 16, seems to short for a public key,
		// so simply return it (as well as the possible parsing error for assuming wrong)
		var botID tbtypes.BotID
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
	record, err := walletCmd.tbClient.GetRecordForKey(pk)
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}
