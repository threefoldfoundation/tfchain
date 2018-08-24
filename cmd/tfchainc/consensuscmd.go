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

	"github.com/spf13/cobra"
)

func createConsensusSubCmds(client *client.CommandLineClient) {
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
	)

	// add commands as wallet sub commands
	client.ConsensusCmd.AddCommand(
		getMintConditionCmd,
	)

	// register flags
	getMintConditionCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &consensusSubCmds.getMintConditionCfg.EncodingType, 0), "encoding",
		cli.EncodingTypeFlagDescription(0))
}

type consensusSubCmds struct {
	cli                 *client.CommandLineClient
	getMintConditionCfg struct {
		EncodingType cli.EncodingType
	}
}

func (consensusSubCmds *consensusSubCmds) getMintCondition(cmd *cobra.Command, args []string) {
	mintConditionGetter := &cliMintConditionGetter{client: consensusSubCmds.cli}

	var (
		mintCondition rivinetypes.UnlockConditionProxy
		err           error
	)

	switch len(args) {
	case 0:
		// get active mint condition for the latest block height
		mintCondition, err = mintConditionGetter.GetActiveMintCondition()
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
		mintCondition, err = mintConditionGetter.GetMintConditionAt(rivinetypes.BlockHeight(height))
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
