package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/rivine/rivine/pkg/client"
	rivinetypes "github.com/rivine/rivine/types"

	"github.com/spf13/cobra"
)

func createWalletSubCmds(cli *client.CommandLineClient) {
	walletSubCmds := &walletSubCmds{cli: cli}

	// define commands
	var (
		createCoinCreationTxCmd = &cobra.Command{
			Use:   "coincreationtransaction <dest>|<rawCondition> <amount> [<dest>|<rawCondition> <amount>]...",
			Short: "Create a new coin creation transaction",
			Long: `Create a new coin creation transaction using the given outputs.
The outputs can be given as a pair of value and a raw output condition (or
address, which resolved to a singlesignature condition).

Amounts have to be given expressed in the OneCoin unit, and without the unit of currency.
Decimals are possible and have to be defined using the decimal point.

The Minimum Miner Fee will be added on top of the total given amount automatically.

The returned (raw) CoinCreationTransaction still has to be signed, prior to sending.
	`,
			Run: walletSubCmds.createCoinCreationTxCmd,
		}
	)

	// add commands as wallet sub commands
	var walletCreateCmd *cobra.Command
	for _, cmd := range cli.WalletCmd.Commands() {
		if cmd.Name() == "create" {
			walletCreateCmd = cmd
			break
		}
	}
	if walletCreateCmd == nil {
		panic("wallet create cmd does not exist")
	}
	walletCreateCmd.AddCommand(
		createCoinCreationTxCmd,
	)

	// register flags
	createCoinCreationTxCmd.Flags().StringVar(
		&walletSubCmds.coinCreationTxCfg.Description, "description", "",
		"optionally add a description to describe the origins of the coin creation, added as arbitrary data")
}

type walletSubCmds struct {
	cli               *client.CommandLineClient
	coinCreationTxCfg struct {
		Description string
	}
}

func (walletSubCmds *walletSubCmds) createCoinCreationTxCmd(cmd *cobra.Command, args []string) {
	currencyConvertor := walletSubCmds.cli.CreateCurrencyConvertor()

	// Check that the remaining args are condition + value pairs
	if len(args)%2 != 0 {
		cmd.UsageFunc()
		client.Die("Invalid arguments. Arguments must be of the form <dest>|<rawCondition> <amount> [<dest>|<rawCondition> <amount>]...")
	}

	// parse the remainder as output coditions and values
	pairs, err := parsePairedOutputs(args, currencyConvertor.ParseCoinString)
	if err != nil {
		cmd.UsageFunc()(cmd)
		client.Die(err)
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

		// try to parse it as an unlock hash
		var uh rivinetypes.UnlockHash
		err = uh.LoadString(args[i])
		if err == nil {
			// parsing as an unlock hash was succesfull, store the pair and continue to the next pair
			pair.Condition = rivinetypes.NewCondition(rivinetypes.NewUnlockHashCondition(uh))
			pairs = append(pairs, pair)
			continue
		}

		// try to parse it as a JSON-encoded unlock condition
		err = pair.Condition.UnmarshalJSON([]byte(args[i]))
		if err != nil {
			err = fmt.Errorf("condition has to be UnlockHash or JSON-encoded UnlockCondition, output #%d's was neither", i/2)
			return
		}
		pairs = append(pairs, pair)
	}
	return
}
