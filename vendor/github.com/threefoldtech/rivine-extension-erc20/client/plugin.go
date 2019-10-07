package client

import (
	"fmt"

	erc20api "github.com/threefoldtech/rivine-extension-erc20/http"
	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"
	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/types"
)

// PluginClient is used to be able to get auth information from
// a daemon that has the authcointx extension enabled and running.
type PluginClient struct {
	bc           *client.BaseClient
	rootEndpoint string
}

var (
	_ erc20types.ERC20Registry = (*PluginClient)(nil)
)

// NewPluginConsensusClient creates a new PluginClient,
// that can be used for easy interaction with the API exposed via the Consensus endpoints
func NewPluginConsensusClient(bc *client.BaseClient) *PluginClient {
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
func NewPluginExplorerClient(bc *client.BaseClient) *PluginClient {
	if bc == nil {
		panic("no BaseClient given")
	}
	return &PluginClient{
		bc:           bc,
		rootEndpoint: "/explorer",
	}
}

func (client *PluginClient) GetERC20AddressForTFTAddress(uh types.UnlockHash) (erc20types.ERC20Address, bool, error) {
	var resp erc20api.GetERC20RelatedAddress
	err := client.bc.HTTP().GetWithResponse(fmt.Sprintf("%s/erc20/addresses/%s", client.rootEndpoint, uh.String()), &resp)
	if err != nil {
		if err == api.ErrStatusNotFound {
			return erc20types.ERC20Address{}, false, nil
		}
		return erc20types.ERC20Address{}, false, fmt.Errorf("failed to get related ERC20 address for unlockhash %s from daemon: %v", uh.String(), err)
	}
	return resp.ERC20Address, true, nil

}

func (client *PluginClient) GetTFTAddressForERC20Address(addr erc20types.ERC20Address) (types.UnlockHash, bool, error) {
	var resp erc20api.GetERC20RelatedAddress
	err := client.bc.HTTP().GetWithResponse(fmt.Sprintf("%s/erc20/addresses/%s", client.rootEndpoint, addr.String()), &resp)
	if err != nil {
		if err == api.ErrStatusNotFound {
			return types.UnlockHash{}, false, nil
		}
		return types.UnlockHash{}, false, fmt.Errorf("failed to get related unlockhash for ERC20 address %s from daemon: %v", addr.String(), err)
	}
	return resp.TFTAddress, true, nil
}

func (client *PluginClient) GetTFTTransactionIDForERC20TransactionID(hash erc20types.ERC20Hash) (types.TransactionID, bool, error) {
	var resp erc20api.GetERC20TransactionID
	err := client.bc.HTTP().GetWithResponse(fmt.Sprintf("%s/erc20/transactions/%s", client.rootEndpoint, hash.String()), &resp)
	if err != nil {
		if err == api.ErrStatusNotFound {
			return types.TransactionID{}, false, nil
		}
		return types.TransactionID{}, false, fmt.Errorf("failed to get related TFT transaction for ERC20 transaction %s from daemon: %v", hash.String(), err)
	}
	return resp.TfchainTransactionID, true, nil
}
