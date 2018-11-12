package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/threefoldtech/rivine/encoding"
	"github.com/threefoldtech/rivine/pkg/cli"
	rivinecli "github.com/threefoldtech/rivine/pkg/client"
	rivinetypes "github.com/threefoldtech/rivine/types"
	"github.com/threefoldfoundation/tfchain/cmd/tfchainc/internal"
	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/spf13/cobra"
)

func createConsensusSubCmds(client *rivinecli.CommandLineClient) {
	consensusSubCmds := &consensusSubCmds{cli: client}

	// define commands
	var (
		getMintConditionCmd = &cobra.Command{
			Use:   "mintcondition [height]",
			Short: "Get the active mint condition",
			Long: `Get the active mint condition,
either the one active for the current block height,
or the one for the given block height.
`,
			Run: consensusSubCmds.getMintCondition,
		}

		getBotRecordCmd = &cobra.Command{
			Use:   "botrecord (id|pubKey|name)",
			Short: "Get the bot record linked to the given info",
			Long: `Get the bot record linked to the given,
id, public key or name.
`,
			Run: rivinecli.Wrap(consensusSubCmds.getBotRecord),
		}

		getBotTransactionsCmd = &cobra.Command{
			Use:   "bottransactions id",
			Short: "Get the transactions created by the given bot",
			Long:  `Get the transactions that created the given bot's record and updated it.`,
			Run:   rivinecli.Wrap(consensusSubCmds.getBotTransactions),
		}
	)

	// add commands as wallet sub commands
	client.ConsensusCmd.AddCommand(
		getMintConditionCmd,
		getBotRecordCmd,
		getBotTransactionsCmd,
	)

	// register flags
	getMintConditionCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &consensusSubCmds.getMintConditionCfg.EncodingType, 0), "encoding",
		cli.EncodingTypeFlagDescription(0))
	getBotRecordCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &consensusSubCmds.getBotRecordCfg.EncodingType, 0), "encoding",
		cli.EncodingTypeFlagDescription(0))
	getBotTransactionsCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &consensusSubCmds.getBotTransactionsCfg.EncodingType, 0), "encoding",
		cli.EncodingTypeFlagDescription(0))
}

type consensusSubCmds struct {
	cli                 *rivinecli.CommandLineClient
	getMintConditionCfg struct {
		EncodingType cli.EncodingType
	}
	getBotRecordCfg struct {
		EncodingType cli.EncodingType
	}
	getBotTransactionsCfg struct {
		EncodingType cli.EncodingType
	}
}

func (consensusSubCmds *consensusSubCmds) getMintCondition(cmd *cobra.Command, args []string) {
	txDBReader := internal.NewTransactionDBConsensusClient(consensusSubCmds.cli)

	var (
		mintCondition rivinetypes.UnlockConditionProxy
		err           error
	)

	switch len(args) {
	case 0:
		// get active mint condition for the latest block height
		mintCondition, err = txDBReader.GetActiveMintCondition()
		if err != nil {
			cli.DieWithError("failed to get the active mint condition", err)
		}

	case 1:
		// get active mint condition for a given block height
		height, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			cmd.UsageFunc()
			cli.DieWithError("invalid block height given", err)
		}
		mintCondition, err = txDBReader.GetMintConditionAt(rivinetypes.BlockHeight(height))
		if err != nil {
			cli.DieWithError("failed to get the mint condition at the given block height", err)
		}

	default:
		cmd.UsageFunc()
		cli.Die("Invalid amount of arguments. One optional pos argument can be given, a valid block height.")
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch consensusSubCmds.getMintConditionCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	case cli.EncodingTypeHex:
		encode = func(v interface{}) error {
			b := encoding.Marshal(v)
			fmt.Println(hex.EncodeToString(b))
			return nil
		}
	}
	err = encode(mintCondition)
	if err != nil {
		cli.DieWithError("failed to encode mint condition", err)
	}
}

func (consensusSubCmds *consensusSubCmds) getBotRecord(str string) {
	txDBReader := internal.NewTransactionDBConsensusClient(consensusSubCmds.cli)
	record, err := txDBReader.GetRecordForString(str)
	if err != nil {
		cli.DieWithError("error while fetching the 3bot record", err)
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch consensusSubCmds.getBotRecordCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	case cli.EncodingTypeHex:
		encode = func(v interface{}) error {
			b := encoding.Marshal(v)
			fmt.Println(hex.EncodeToString(b))
			return nil
		}
	}
	err = encode(record)
	if err != nil {
		cli.DieWithError("failed to encode 3bot record", err)
	}
}

func (consensusSubCmds *consensusSubCmds) getBotTransactions(str string) {
	var botID types.BotID
	err := botID.LoadString(str)
	if err != nil {
		cli.DieWithError("failed to parse 3bot ID pos arg", err)
	}

	txDBReader := internal.NewTransactionDBConsensusClient(consensusSubCmds.cli)
	ids, err := txDBReader.GetBotTransactionIdentifiers(botID)
	if err != nil {
		cli.DieWithError("error while fetching the 3bot transactions", err)
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch consensusSubCmds.getBotTransactionsCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	case cli.EncodingTypeHex:
		encode = func(v interface{}) error {
			b := encoding.Marshal(v)
			fmt.Println(hex.EncodeToString(b))
			return nil
		}
	}
	err = encode(ids)
	if err != nil {
		cli.DieWithError("failed to encode transactions", err)
	}
}
