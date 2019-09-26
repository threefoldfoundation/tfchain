package client

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	erc20api "github.com/threefoldtech/rivine-extension-erc20/http"
	"github.com/threefoldtech/rivine/pkg/cli"
	client "github.com/threefoldtech/rivine/pkg/client"
)

// CreateERC20Cmd creates rootcommand for ERC20 and adds a subcommand
// if rootcommand executed the user will also see the output of the syncing status of ethereum
func CreateERC20Cmd(client *client.CommandLineClient) *cobra.Command {
	erc20SubCmds := &erc20SubCmds{cli: client}

	// define Rootcommand
	var (
		rootCmd = &cobra.Command{
			Use:   "erc20",
			Short: "Perform erc20 actions",
			Long:  "Perform erc20 actions",
			Run:   erc20SubCmds.getSyncingStatus,
		}
		getSyncingStatusCmd = &cobra.Command{
			Use:   "syncstatus",
			Short: "Get the ethereum sync status",
			Long:  `Get the ethereum chain sync status.`,
			Run:   erc20SubCmds.getSyncingStatus,
		}
	)

	rootCmd.AddCommand(getSyncingStatusCmd)

	// register flags
	getSyncingStatusCmd.Flags().Var(
		cli.NewEncodingTypeFlag(cli.EncodingTypeHuman, &erc20SubCmds.getSyncingStatusCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	return rootCmd
}

type erc20SubCmds struct {
	cli                 *client.CommandLineClient
	getSyncingStatusCfg struct {
		EncodingType cli.EncodingType
	}
}

// getSyncingStatus Gets the ethereum blockchain syncing status from the deamon API
func (erc20SubCmds *erc20SubCmds) getSyncingStatus(cmd *cobra.Command, args []string) {
	var syncingStatus erc20api.ERC20SyncingStatus

	err := erc20SubCmds.cli.GetWithResponse("/erc20/downloader/status", &syncingStatus)
	if err != nil {
		cli.DieWithError("error while fetching the syncing status", err)
	}

	// encode depending on the encoding flag
	switch erc20SubCmds.getSyncingStatusCfg.EncodingType {
	case cli.EncodingTypeHuman:
		fmt.Printf(`Starting block height: %d
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
