package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/rivine/rivine/encoding"
	"github.com/rivine/rivine/pkg/cli"
	"github.com/rivine/rivine/pkg/client"
	rivinetypes "github.com/rivine/rivine/types"
	"github.com/threefoldfoundation/tfchain/pkg/api"

	"github.com/spf13/cobra"
)

func createExplorerSubCmds(client *client.CommandLineClient) {
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
	)

	// add commands as wallet sub commands
	client.ExploreCmd.AddCommand(
		getMintConditionCmd,
	)

	// register flags
	getMintConditionCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &explorerSubCmds.getMintConditionCfg.EncodingType, 0), "encoding",
		cli.EncodingTypeFlagDescription(0))
}

type explorerSubCmds struct {
	cli                 *client.CommandLineClient
	getMintConditionCfg struct {
		EncodingType cli.EncodingType
	}
}

func (explorerSubCmds *explorerSubCmds) getMintCondition(cmd *cobra.Command, args []string) {
	var (
		mintCondition rivinetypes.UnlockConditionProxy
		err           error
	)

	switch len(args) {
	case 0:
		// get active mint condition for the latest block height
		var result api.TransactionDBGetMintCondition
		err := explorerSubCmds.cli.GetAPI("/explorer/mintcondition", &result)
		if err != nil {
			cli.DieWithError("failed to get the active mint condition from the explorer", err)
		}
		mintCondition = result.MintCondition

	case 1:
		// get active mint condition for a given block height
		height, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			cmd.UsageFunc()
			cli.DieWithError("invalid block height given", err)
		}
		var result api.TransactionDBGetMintCondition
		err = explorerSubCmds.cli.GetAPI(fmt.Sprintf("/explorer/mintcondition/%d", height), &result)
		if err != nil {
			cli.DieWithError("failed to get the mint condition from explorer at the given block height", err)
		}
		mintCondition = result.MintCondition

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
