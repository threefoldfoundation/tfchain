package client

import (
	"errors"
	"fmt"
	"strconv"

	tbapi "github.com/threefoldfoundation/tfchain/extensions/threebot/api"
	tbtypes "github.com/threefoldfoundation/tfchain/extensions/threebot/types"
	"github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/types"
)

// PluginClient is used to be able to get auth information from
// a daemon that has the authcointx extension enabled and running.
type PluginClient struct {
	bc           client.BaseClient
	rootEndpoint string
}

var (
	_ tbtypes.BotRecordReadRegistry = (*PluginClient)(nil)
)

// NewPluginConsensusClient creates a new PluginClient,
// that can be used for easy interaction with the API exposed via the Consensus endpoints
func NewPluginConsensusClient(bc client.BaseClient) *PluginClient {
	if bc == nil {
		panic("no BaseClient given")
	}
	return &PluginClient{
		bc:           bc,
		rootEndpoint: "/consensus",
	}
}

// NewPluginExplorerClient creates a new PluginClient,
// that can be used for easy interaction with the API exposed via the Explorer endpoints
func NewPluginExplorerClient(bc client.BaseClient) *PluginClient {
	if bc == nil {
		panic("no BaseClient given")
	}
	return &PluginClient{
		bc:           bc,
		rootEndpoint: "/explorer",
	}
}

func (client *PluginClient) GetRecordForKey(publicKey types.PublicKey) (*tbtypes.BotRecord, error) {
	var result tbapi.GetBotRecord
	err := client.bc.HTTP().GetWithResponse(fmt.Sprintf("%s/3bot/%s", client.rootEndpoint, publicKey.String()), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot record for crypto public-key %s from daemon: %v", publicKey.String(), err)
	}
	return &result.Record, nil
}

func (client *PluginClient) GetRecordForName(name tbtypes.BotName) (*tbtypes.BotRecord, error) {
	var result tbapi.GetBotRecord
	err := client.bc.HTTP().GetWithResponse(fmt.Sprintf("%s/whois/3bot/%s", client.rootEndpoint, name.String()), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot record for botname %s from daemon: %v", name.String(), err)
	}
	return &result.Record, nil
}

func (client *PluginClient) GetRecordForID(id tbtypes.BotID) (*tbtypes.BotRecord, error) {
	var result tbapi.GetBotRecord
	err := client.bc.HTTP().GetWithResponse(fmt.Sprintf("%s/3bot/%d", client.rootEndpoint, id), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot record for ID %d from daemon: %v", id, err)
	}
	return &result.Record, nil
}

func (client *PluginClient) BotRecordForString(str string) (*tbtypes.BotRecord, error) {
	// try str as a BotID
	id, err := strconv.ParseUint(str, 10, 32)
	if err == nil {
		return client.GetRecordForID(tbtypes.BotID(id))
	}

	// try str as a BotName
	var name tbtypes.BotName
	err = name.LoadString(str)
	if err == nil {
		return client.GetRecordForName(name)
	}

	// should be a public key, last choice
	var publicKey types.PublicKey
	err = publicKey.LoadString(str)
	if err != nil {
		return nil, errors.New("argument should be a valid BotID, BotName or PublicKey")
	}
	return client.GetRecordForKey(publicKey)
}

func (client *PluginClient) GetBotTransactionIdentifiers(id tbtypes.BotID) ([]types.TransactionID, error) {
	var result tbapi.GetBotTransactions
	err := client.bc.HTTP().GetWithResponse(fmt.Sprintf("/consensus/3bot/%s/transactions", id.String()), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot transactions for ID %s from daemon: %v", id.String(), err)
	}
	return result.Identifiers, nil
}
