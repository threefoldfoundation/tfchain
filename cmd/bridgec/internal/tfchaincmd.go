package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/cli"
	rivinec "github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/types"
)

// createTFChainCommand creates a tfchain syncstatus command
func createTFChainCommand(client *CommandLineClient) *cobra.Command {
	tfchaincmd := &tfchaincmd{cli: client}

	// define Rootcommand
	var (
		rootCmd = &cobra.Command{
			Use:   "tfchain",
			Short: "Returns tfchain network syncing status.",
			Run:   rivinec.Wrap(tfchaincmd.getSyncingStatus),
		}
	)

	// register flags
	rootCmd.Flags().Var(
		cli.NewEncodingTypeFlag(cli.EncodingTypeHuman, &tfchaincmd.getSyncingStatusCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	return rootCmd
}

type tfchaincmd struct {
	cli                 *CommandLineClient
	getSyncingStatusCfg struct {
		EncodingType cli.EncodingType
	}
}

// getSyncingStatus Gets the ethereum blockchain syncing status from the deamon API
func (tfchaincmd *tfchaincmd) getSyncingStatus() {
	var cg api.ConsensusGET
	err := tfchaincmd.cli.GetWithResponse("/consensus", &cg)
	if err != nil {
		cli.DieWithError("error while fetching the consensus status", err)
	}

	// encode depending on the encoding flag
	switch tfchaincmd.getSyncingStatusCfg.EncodingType {
	case cli.EncodingTypeHuman:
		if cg.Synced {
			fmt.Printf(`Tfchain sync status:
Synced: %v
Block:  %v
Height: %v
Target: %v
`, YesNo(cg.Synced), cg.CurrentBlock, cg.Height, cg.Target)
		} else {
			estimatedHeight := tfchaincmd.estimatedHeightAt(time.Now())
			estimatedProgress := float64(cg.Height) / float64(estimatedHeight) * 100
			if estimatedProgress > 99 {
				estimatedProgress = 99
			}
			fmt.Printf(`Tfchain sync status:
Synced: %v
Height: %v
Progress (estimated): %.2f%%
`, YesNo(cg.Synced), cg.Height, estimatedProgress)
		}
	case cli.EncodingTypeJSON:
		err = json.NewEncoder(os.Stdout).Encode(cg)
		if err != nil {
			cli.DieWithError("failed to encode syncing status", err)
		}
	}
}

// EstimatedHeightAt returns the estimated block height for the given time.
// Block height is estimated by calculating the minutes since a known block in
// the past and dividing by 10 minutes (the block time).
func (tfchaincmd *tfchaincmd) estimatedHeightAt(t time.Time) types.BlockHeight {
	if tfchaincmd.cli.Config.GenesisBlockTimestamp == 0 {
		panic("GenesisBlockTimestamp is undefined")
	}
	return estimatedHeightBetween(
		int64(tfchaincmd.cli.Config.GenesisBlockTimestamp),
		t.Unix(),
		tfchaincmd.cli.Config.BlockFrequencyInSeconds,
	)
}
