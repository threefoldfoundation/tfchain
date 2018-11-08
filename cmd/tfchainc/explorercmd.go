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

	"github.com/spf13/cobra"
)

func createExplorerSubCmds(client *rivinecli.CommandLineClient) {
	explorerSubCmds := &explorerSubCmds{cli: client}

	// define commands
	var (
		getMintConditionCmd = &cobra.Command{
			Use:   "mintcondition [height]",
			Short: "Get the active mint condition",
			Long: `Get the active mint condition,
either the one active for the current block height,
or the one for the given block height.
`,
			Run: explorerSubCmds.getMintCondition,
		}

		getBotRecordCmd = &cobra.Command{
			Use:   "botrecord (id|pubKey|name)",
			Short: "Get the bot record linked to the given info",
			Long: `Get the bot record linked to the given,
id, public key or name.
`,
			Run: rivinecli.Wrap(explorerSubCmds.getBotRecord),
		}
	)

	// add commands as wallet sub commands
	client.ExploreCmd.AddCommand(
		getMintConditionCmd,
		getBotRecordCmd,
	)

	// register flags
	getMintConditionCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &explorerSubCmds.getMintConditionCfg.EncodingType, 0), "encoding",
		cli.EncodingTypeFlagDescription(0))
	getBotRecordCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &explorerSubCmds.getBotRecordCfg.EncodingType, 0), "encoding",
		cli.EncodingTypeFlagDescription(0))
}

type explorerSubCmds struct {
	cli                 *rivinecli.CommandLineClient
	getMintConditionCfg struct {
		EncodingType cli.EncodingType
	}
	getBotRecordCfg struct {
		EncodingType cli.EncodingType
	}
}

func (explorerSubCmds *explorerSubCmds) getMintCondition(cmd *cobra.Command, args []string) {
	txDBReader := internal.NewTransactionDBExplorerClient(explorerSubCmds.cli)

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
	switch explorerSubCmds.getMintConditionCfg.EncodingType {
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

func (explorerSubCmds *explorerSubCmds) getBotRecord(str string) {
	txDBReader := internal.NewTransactionDBExplorerClient(explorerSubCmds.cli)
	record, err := txDBReader.GetRecordForString(str)
	if err != nil {
		cli.DieWithError("error while fetching the 3bot record", err)
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch explorerSubCmds.getMintConditionCfg.EncodingType {
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
