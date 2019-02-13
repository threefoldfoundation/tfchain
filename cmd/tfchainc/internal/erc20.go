package internal

import (
	"fmt"

	"github.com/threefoldfoundation/tfchain/pkg/api"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
)

// ERC20Client is used to be able to get the active mint condition,
// the active mint condition at a given block height, as well as any 3bot information
// such that the CLI can also correctly validate a mint-type a 3bot-type transaction,
// without requiring access to the consensus-extended transactiondb,
// normally the validation isn't required on the client side, but it is now possible none the less
type ERC20Client struct {
	client       *CommandLineClient
	rootEndpoint string
}

// NewERC20Client creates a new TransactionDBClient,
// that can be used for easy interaction with the TransactionDB API exposed via the Consensus endpoints
func NewERC20Client(cli *CommandLineClient) *ERC20Client {
	if cli == nil {
		panic("no CommandLineClient given")
	}
	return &ERC20Client{
		client:       cli,
		rootEndpoint: "/erc20",
	}
}

// GetSyncingStatus implements types.MintConditionGetter.GetActiveMintCondition
func (cli *ERC20Client) GetSyncingStatus() (*tftypes.ERC20SyncStatus, error) {
	var result api.ERC20SyncingStatus

	err := cli.client.GetAPI(cli.rootEndpoint+"/downloader/status", &result)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get syncing status from daemon: %v", err)
	}
	return &result.Status, nil
}
