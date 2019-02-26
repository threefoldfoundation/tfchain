package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	tfchainapi "github.com/threefoldfoundation/tfchain/pkg/api"
	api "github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/cli"
	rivinec "github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/types"
)

// CreateConsensusCmd creates rootcommand for consensus
// if rootcommand executed the user will see the output of the syncing status of tfchain and ethereum network
func createConsensusCmd(client *CommandLineClient) (*consensusCmd, *cobra.Command) {
	consensusCmd := &consensusCmd{cli: client}

	// define Rootcommand
	var (
		rootCmd = &cobra.Command{
			Use:   "consensus",
			Short: "Returns tfchain and ethereum network syncing status.",
			Run:   rivinec.Wrap(consensusCmd.getSyncingStatus),
		}
	)

	// register flags
	rootCmd.Flags().Var(
		cli.NewEncodingTypeFlag(cli.EncodingTypeHuman, &consensusCmd.getSyncingStatusCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	return consensusCmd, rootCmd
}

type consensusCmd struct {
	cli                 *CommandLineClient
	getSyncingStatusCfg struct {
		EncodingType cli.EncodingType
	}
}

// getSyncingStatus Gets the ethereum blockchain syncing status from the deamon API
func (consensusCmd *consensusCmd) getSyncingStatus() {
	var syncingStatus tfchainapi.ERC20SyncingStatus

	err := consensusCmd.cli.GetAPI("/erc20/downloader/status", &syncingStatus)
	if err != nil {
		cli.DieWithError("error while fetching the syncing status", err)
	}

	var cg api.ConsensusGET
	err = consensusCmd.cli.GetAPI("/consensus", &cg)
	if err != nil {
		cli.DieWithError("error while fetching the consensus status", err)
	}

	// encode depending on the encoding flag
	switch consensusCmd.getSyncingStatusCfg.EncodingType {
	case cli.EncodingTypeHuman:
		if cg.Synced {
			fmt.Printf(`Tfchain sync status:
Synced: %v
Block:  %v
Height: %v
Target: %v
`, YesNo(cg.Synced), cg.CurrentBlock, cg.Height, cg.Target)
		} else {
			estimatedHeight := consensusCmd.estimatedHeightAt(time.Now())
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
		fmt.Printf(`
Ethereum sync status:
Starting block height: %d
Current block height: %d
Highest block height: %d
`, syncingStatus.Status.StartingBlock, syncingStatus.Status.CurrentBlock, syncingStatus.Status.HighestBlock)
	case cli.EncodingTypeJSON:
		err = json.NewEncoder(os.Stdout).Encode(syncingStatus.Status)
		if err != nil {
			cli.DieWithError("failed to encode syncing status", err)
		}
	}
}

// YesNo returns "Yes" if b is true, and "No" if b is false.
func YesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// EstimatedHeightAt returns the estimated block height for the given time.
// Block height is estimated by calculating the minutes since a known block in
// the past and dividing by 10 minutes (the block time).
func (consensusCmd *consensusCmd) estimatedHeightAt(t time.Time) types.BlockHeight {
	if consensusCmd.cli.Config.GenesisBlockTimestamp == 0 {
		panic("GenesisBlockTimestamp is undefined")
	}
	return estimatedHeightBetween(
		int64(consensusCmd.cli.Config.GenesisBlockTimestamp),
		t.Unix(),
		consensusCmd.cli.Config.BlockFrequencyInSeconds,
	)
}

func estimatedHeightBetween(from, to, blockFrequency int64) types.BlockHeight {
	lifetimeInSeconds := to - from
	if lifetimeInSeconds < blockFrequency {
		return 0
	}
	estimatedHeight := float64(lifetimeInSeconds) / float64(blockFrequency)
	return types.BlockHeight(estimatedHeight + 0.5) // round to the nearest block
}
