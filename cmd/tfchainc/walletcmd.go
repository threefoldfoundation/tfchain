package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/threefoldfoundation/tfchain/cmd/tfchainc/internal"

	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/rivine/rivine/pkg/cli"
	rivinecli "github.com/rivine/rivine/pkg/client"
	rivinetypes "github.com/rivine/rivine/types"

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
	)

	// add commands as wallet sub commands
	client.WalletCmd.RootCmdCreate.AddCommand(
		createMinterDefinitionTxCmd,
		createCoinCreationTxCmd,
	)
	client.WalletCmd.RootCmdSend.AddCommand(
		sendBotRegistrationTxCmd,
	)

	// register flags

	createMinterDefinitionTxCmd.Flags().StringVar(
		&walletSubCmds.minterDefinitionTxCfg.Description, "description", "",
		"optionally add a description to describe the reasons of transfer of minting power, added as arbitrary data")

	createCoinCreationTxCmd.Flags().StringVar(
		&walletSubCmds.coinCreationTxCfg.Description, "description", "",
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
}

type walletSubCmds struct {
	cli                   *rivinecli.CommandLineClient
	minterDefinitionTxCfg struct {
		Description string
	}
	coinCreationTxCfg struct {
		Description string
	}

	sendBotRegistrationTxCfg struct {
		Addresses    []types.NetworkAddress
		Names        []types.BotName
		NrOfMonths   uint8
		PublicKey    types.PublicKey
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

	// if a description is given, use it as arbitrary data
	if n := len(walletSubCmds.minterDefinitionTxCfg.Description); n > 0 {
		tx.ArbitraryData = make([]byte, n)
		copy(tx.ArbitraryData[:], walletSubCmds.minterDefinitionTxCfg.Description[:])
	}

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
	if n := len(walletSubCmds.coinCreationTxCfg.Description); n > 0 {
		tx.ArbitraryData = make([]byte, n)
		copy(tx.ArbitraryData[:], walletSubCmds.coinCreationTxCfg.Description[:])
	}
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
		cli.DieWithError("failed to encode resultd", err)
	}
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
