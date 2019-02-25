package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/threefoldfoundation/tfchain/cmd/bridgec/internal"
	"github.com/threefoldfoundation/tfchain/pkg/api"
	erc20 "github.com/threefoldfoundation/tfchain/pkg/eth/erc20"
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
			Run:   erc20SubCmds.getSyncingStatus,
		}
		getSyncingStatusCmd = &cobra.Command{
			Use:   "syncstatus",
			Short: "Get the ethereum dowloader sync status",
			Run:   erc20SubCmds.getSyncingStatus,
		}
		getBalanceInfoCmd = &cobra.Command{
			Use:   "balance",
			Short: "Get the ethereum balance and address information",
			Run:   erc20SubCmds.getBalanceInfo,
		}
	)

	rootCmd.AddCommand(getSyncingStatusCmd)
	rootCmd.AddCommand(getBalanceInfoCmd)

	// register flags
	getSyncingStatusCmd.Flags().Var(
		cli.NewEncodingTypeFlag(cli.EncodingTypeHuman, &erc20SubCmds.getSyncingStatusCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))
	getBalanceInfoCmd.Flags().Var(
		cli.NewEncodingTypeFlag(cli.EncodingTypeHuman, &erc20SubCmds.getBalanceInfoCfg.EncodingType, cli.EncodingTypeHuman|cli.EncodingTypeJSON), "encoding",
		cli.EncodingTypeFlagDescription(cli.EncodingTypeHuman|cli.EncodingTypeJSON))

	return rootCmd
}

type erc20SubCmds struct {
	cli                 *internal.CommandLineClient
	getSyncingStatusCfg struct {
		EncodingType cli.EncodingType
	}
	getBalanceInfoCfg struct {
		EncodingType cli.EncodingType
	}
}

// getSyncingStatus Gets the ethereum blockchain syncing status from the deamon API
func (erc20SubCmds *erc20SubCmds) getSyncingStatus(cmd *cobra.Command, args []string) {
	var syncingStatus api.ERC20SyncingStatus

	err := erc20SubCmds.cli.GetAPI("/erc20/downloader/status", &syncingStatus)
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

// getSyncingStatus Gets the ethereum blockchain syncing status from the deamon API
func (erc20SubCmds *erc20SubCmds) getBalanceInfo(cmd *cobra.Command, args []string) {
	var balanceInfo api.ERC20BalanceInformation

	err := erc20SubCmds.cli.GetAPI("/erc20/account/balance", &balanceInfo)
	if err != nil {
		cli.DieWithError("error while fetching the balance information", err)
	}

	// encode depending on the encoding flag
	switch erc20SubCmds.getSyncingStatusCfg.EncodingType {
	case cli.EncodingTypeHuman:
		ether := erc20.Denominate(balanceInfo.BalanceInfo.Balance)
		fmt.Printf(`Address: %s
Balance: %s
`, balanceInfo.BalanceInfo.Address.String(), ether)
	case cli.EncodingTypeJSON:
		err = json.NewEncoder(os.Stdout).Encode(balanceInfo.BalanceInfo)
		if err != nil {
			cli.DieWithError("failed to encode balance information", err)
		}
	}
}
