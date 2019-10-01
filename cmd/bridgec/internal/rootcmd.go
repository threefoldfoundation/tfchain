package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	erc20api "github.com/threefoldtech/rivine-extension-erc20/http"
	api "github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/cli"
	rivinec "github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/types"
)

// createRootCmd creates root command for consensus
// if root command executed the user will see the output of the syncing status of tfchain and ethereum network
func createRootCmd(binName, clientName string, client *CommandLineClient) {
	rootCmd := &rootCmd{cli: client}

	// create Rootcommand
	client.RootCmd = &cobra.Command{
		Use:     binName,
		Short:   fmt.Sprintf("%s Client", strings.Title(clientName)),
		Run:     rivinec.Wrap(rootCmd.getSyncingStatus),
		PreRunE: client.preRunE,
	}

	// register flags
	client.RootCmd.Flags().Var(
		cli.NewEncodingTypeFlag(cli.EncodingTypeHuman, &rootCmd.getSyncingStatusCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))
}

type rootCmd struct {
	cli                 *CommandLineClient
	getSyncingStatusCfg struct {
		EncodingType cli.EncodingType
	}
}

// getSyncingStatus Gets the ethereum blockchain syncing status from the deamon API
func (rootCmd *rootCmd) getSyncingStatus() {
	var syncingStatus erc20api.ERC20SyncingStatus

	err := rootCmd.cli.GetWithResponse("/erc20/downloader/status", &syncingStatus)
	if err != nil {
		cli.DieWithError("error while fetching the syncing status", err)
	}

	var cg api.ConsensusGET
	err = rootCmd.cli.GetWithResponse("/consensus", &cg)
	if err != nil {
		cli.DieWithError("error while fetching the consensus status", err)
	}

	// encode depending on the encoding flag
	switch rootCmd.getSyncingStatusCfg.EncodingType {
	case cli.EncodingTypeHuman:
		if cg.Synced {
			fmt.Printf(`Tfchain sync status:
Synced: %v
Block:  %v
Height: %v
Target: %v
`, YesNo(cg.Synced), cg.CurrentBlock, cg.Height, cg.Target)
		} else {
			estimatedHeight := rootCmd.estimatedHeightAt(time.Now())
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
		err = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"tfchain":  cg,
			"ethereum": syncingStatus.Status,
		})
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
func (rootCmd *rootCmd) estimatedHeightAt(t time.Time) types.BlockHeight {
	if rootCmd.cli.Config.GenesisBlockTimestamp == 0 {
		panic("GenesisBlockTimestamp is undefined")
	}
	return estimatedHeightBetween(
		int64(rootCmd.cli.Config.GenesisBlockTimestamp),
		t.Unix(),
		rootCmd.cli.Config.BlockFrequencyInSeconds,
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
