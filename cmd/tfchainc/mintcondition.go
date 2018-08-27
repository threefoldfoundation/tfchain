package main

import (
	"fmt"

	"github.com/threefoldfoundation/tfchain/pkg/api"
	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/rivine/rivine/pkg/client"
	rivinetypes "github.com/rivine/rivine/types"
)

// cliMintConditionGetter is used to be able to get the active mint condition,
// as well as the active mint condition at a given block height,
// such that the CLI can also correctly validate a mint-type transaction,
// without requiring access to the consensus-extended transactiondb,
// normally the validation isn't required on the client side, but it is now possible none the less
type cliMintConditionGetter struct {
	client *client.CommandLineClient
}

var (
	// ensure cliMintConditionGetter implements the MintConditionGetter interface
	_ types.MintConditionGetter = (*cliMintConditionGetter)(nil)
)

// GetActiveMintCondition implements types.MintConditionGetter.GetActiveMintCondition
func (cli *cliMintConditionGetter) GetActiveMintCondition() (rivinetypes.UnlockConditionProxy, error) {
	var result api.TransactionDBGetMintCondition
	err := cli.client.GetAPI("/consensus/mintcondition", &result)
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, fmt.Errorf(
			"failed to get active mint condition from daemon: %v", err)
	}
	return result.MintCondition, nil
}

// GetMintConditionAt implements types.MintConditionGetter.GetMintConditionAt
func (cli *cliMintConditionGetter) GetMintConditionAt(height rivinetypes.BlockHeight) (rivinetypes.UnlockConditionProxy, error) {
	var result api.TransactionDBGetMintCondition
	err := cli.client.GetAPI(fmt.Sprintf("/consensus/mintcondition/%d", height), &result)
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, fmt.Errorf(
			"failed to get mint condition at height %d from daemon: %v", height, err)
	}
	return result.MintCondition, nil
}
