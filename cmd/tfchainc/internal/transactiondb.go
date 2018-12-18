package internal

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/threefoldfoundation/tfchain/pkg/api"
	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/threefoldtech/rivine/pkg/client"
	rivinetypes "github.com/threefoldtech/rivine/types"
)

// TransactionDBClient is used to be able to get the active mint condition,
// the active mint condition at a given block height, as well as any 3bot information
// such that the CLI can also correctly validate a mint-type a 3bot-type transaction,
// without requiring access to the consensus-extended transactiondb,
// normally the validation isn't required on the client side, but it is now possible none the less
type TransactionDBClient struct {
	client       *client.CommandLineClient
	rootEndpoint string
}

// NewTransactionDBConsensusClient creates a new TransactionDBClient,
// that can be used for easy interaction with the TransactionDB API exposed via the Consensus endpoints
func NewTransactionDBConsensusClient(cli *client.CommandLineClient) *TransactionDBClient {
	if cli == nil {
		panic("no CommandLineClient given")
	}
	return &TransactionDBClient{
		client:       cli,
		rootEndpoint: "/consensus",
	}
}

// NewTransactionDBExplorerClient creates a new TransactionDBClient,
// that can be used for easy interaction with the TransactionDB API exposed via the Explorer endpoints
func NewTransactionDBExplorerClient(cli *client.CommandLineClient) *TransactionDBClient {
	if cli == nil {
		panic("no CommandLineClient given")
	}
	return &TransactionDBClient{
		client:       cli,
		rootEndpoint: "/explorer",
	}
}

var (
	// ensure TransactionDBClient implements the MintConditionGetter interface
	_ types.MintConditionGetter = (*TransactionDBClient)(nil)
	// ensure TransactionDBClient implements the BotRecordReadRegistry interface
	_ types.BotRecordReadRegistry = (*TransactionDBClient)(nil)
)

// GetActiveMintCondition implements types.MintConditionGetter.GetActiveMintCondition
func (cli *TransactionDBClient) GetActiveMintCondition() (rivinetypes.UnlockConditionProxy, error) {
	var result api.TransactionDBGetMintCondition
	err := cli.client.GetAPI(cli.rootEndpoint+"/mintcondition", &result)
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, fmt.Errorf(
			"failed to get active mint condition from daemon: %v", err)
	}
	return result.MintCondition, nil
}

// GetMintConditionAt implements types.MintConditionGetter.GetMintConditionAt
func (cli *TransactionDBClient) GetMintConditionAt(height rivinetypes.BlockHeight) (rivinetypes.UnlockConditionProxy, error) {
	var result api.TransactionDBGetMintCondition
	err := cli.client.GetAPI(fmt.Sprintf("%s/mintcondition/%d", cli.rootEndpoint, height), &result)
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, fmt.Errorf(
			"failed to get mint condition at height %d from daemon: %v", height, err)
	}
	return result.MintCondition, nil
}

// GetRecordForID implements types.BotRecordReadRegistry.GetRecordForID
func (cli *TransactionDBClient) GetRecordForID(id types.BotID) (*types.BotRecord, error) {
	var result api.TransactionDBGetBotRecord
	err := cli.client.GetAPI(fmt.Sprintf("%s/3bot/%d", cli.rootEndpoint, id), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot record for ID %d from daemon: %v", id, err)
	}
	return &result.Record, nil
}

// GetRecordForKey implements types.BotRecordReadRegistry.GetRecordForKey
func (cli *TransactionDBClient) GetRecordForKey(key rivinetypes.PublicKey) (*types.BotRecord, error) {
	var result api.TransactionDBGetBotRecord
	err := cli.client.GetAPI(fmt.Sprintf("%s/3bot/%s", cli.rootEndpoint, key.String()), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot record for crypto public-key %s from daemon: %v", key.String(), err)
	}
	return &result.Record, nil
}

// GetRecordForName implements types.BotRecordReadRegistry.GetRecordForName
func (cli *TransactionDBClient) GetRecordForName(name types.BotName) (*types.BotRecord, error) {
	var result api.TransactionDBGetBotRecord
	err := cli.client.GetAPI(fmt.Sprintf("%s/whois/3bot/%s", cli.rootEndpoint, name.String()), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot record for botname %s from daemon: %v", name.String(), err)
	}
	return &result.Record, nil
}

// GetBotTransactionIdentifiers implements types.BotRecordReadRegistry.GetBotTransactionIdentifiers
func (cli *TransactionDBClient) GetBotTransactionIdentifiers(id types.BotID) ([]rivinetypes.TransactionID, error) {
	var result api.TransactionDBGetBotTransactions
	err := cli.client.GetAPI(fmt.Sprintf("%s/3bot/%s/transactions", cli.rootEndpoint, id.String()), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot transactions for ID %s from daemon: %v", id.String(), err)
	}
	return result.Identifiers, nil
}

// GetRecordForString gets a bot record for either a ID, (public) key or name,
// as long as it is referenced to a registered Bot.
func (cli *TransactionDBClient) GetRecordForString(str string) (*types.BotRecord, error) {
	// try str as a BotID
	id, err := strconv.ParseUint(str, 10, 32)
	if err == nil {
		return cli.GetRecordForID(types.BotID(id))
	}

	// try str as a BotName
	var name types.BotName
	err = name.LoadString(str)
	if err == nil {
		return cli.GetRecordForName(name)
	}

	// should be a public key, last choice
	var publicKey rivinetypes.PublicKey
	err = publicKey.LoadString(str)
	if err != nil {
		return nil, errors.New("argument should be a valid BotID, BotName or PublicKey")
	}
	return cli.GetRecordForKey(publicKey)
}

// GetERC20AddressForTFTAddress implements types.ERC20Registry.GetERC20AddressForTFTAddress
func (cli *TransactionDBClient) GetERC20AddressForTFTAddress(uh rivinetypes.UnlockHash) (types.ERC20Address, error) {
	var result api.TransactionDBGetERC20RelatedAddress
	err := cli.client.GetAPI(fmt.Sprintf("%s/erc20/addresses/%s", cli.rootEndpoint, uh.String()), &result)
	if err != nil {
		return types.ERC20Address{}, fmt.Errorf("failed to get ERC20 Info for TFT address %s from daemon: %v", uh.String(), err)
	}
	return result.ERC20Address, nil
}

// GetTFTTransactionIDForERC20TransactionID implements types.ERC20Registry.GetTFTTransactionIDForERC20TransactionID
func (cli *TransactionDBClient) GetTFTTransactionIDForERC20TransactionID(txid types.ERC20TransactionID) (rivinetypes.TransactionID, error) {
	var result api.TransactionDBGetERC20TransactionID
	err := cli.client.GetAPI(fmt.Sprintf("%s/erc20/transactions/%s", cli.rootEndpoint, txid.String()), &result)
	if err != nil {
		return rivinetypes.TransactionID{}, fmt.Errorf("failed to get info linked to ERC20 Transaction ID %s from daemon: %v", txid.String(), err)
	}
	return result.TfchainTransactionID, nil
}
