package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/threefoldfoundation/tfchain/cmd/tfchainc/internal"
	"github.com/threefoldfoundation/tfchain/pkg/api"
	"github.com/threefoldtech/rivine/pkg/cli"
)

// createERC20Cmd creates rootcommand for ERC20 and adds a subcommand
// if rootcommand executed the user will also see the output of the syncing status of ethereum
func createERC20Cmd(client *internal.CommandLineClient) *cobra.Command {
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
		cli.NewEncodingTypeFlag(0, &erc20SubCmds.getSyncingStatusCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	return rootCmd
}

type erc20SubCmds struct {
	cli                 *internal.CommandLineClient
	getSyncingStatusCfg struct {
		EncodingType cli.EncodingType
	}
}

// getSyncingStatus Gets the ethereum blockchain syncing status from the deamon API
func (erc20SubCmds *erc20SubCmds) getSyncingStatus(cmd *cobra.Command, args []string) {
	client := erc20SubCmds.cli
	var syncingStatus api.ERC20SyncingStatus

	err := client.GetAPI("/erc20/downloader/status", &syncingStatus)
	if err != nil {
		cli.DieWithError("error while fetching the syncing status", err)
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch erc20SubCmds.getSyncingStatusCfg.EncodingType {
	case cli.EncodingTypeHuman:
		encode = func(val interface{}) error {
			syncing := syncingStatus.Status.Synchronising
			if !syncing {
				fmt.Println("ERC20 node is not syncronising")
			} else {
				fmt.Printf(`ERC20 node is currently syncronising...
Starting block height: %d
Current block height: %d
Highest block height: %d
`, syncingStatus.Status.StartingBlock, syncingStatus.Status.CurrentBlock, syncingStatus.Status.HighestBlock)
			}
			return nil
		}
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	}
	err = encode(syncingStatus)
	if err != nil {
		cli.DieWithError("failed to encode syncing status", err)
	}
}
