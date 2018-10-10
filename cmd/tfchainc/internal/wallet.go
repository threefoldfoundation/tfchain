package internal

import (
	"encoding/json"
	"fmt"

	"github.com/threefoldfoundation/tfchain/pkg/api"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/rivine/rivine/pkg/client"
	"github.com/rivine/rivine/types"
)

// WalletClient is used to easily interact with the wallet through the HTTP REST API.
type WalletClient struct {
	client *client.CommandLineClient
}

// NewWalletClient creates a new WalletClient,
// that can be used for easy interaction with the Wallet API exposed via the HTTP REST API.
func NewWalletClient(cli *client.CommandLineClient) *WalletClient {
	if cli == nil {
		panic("no CommandLineClient given")
	}
	return &WalletClient{
		client: cli,
	}
}

// NewPublicKey creates a new public key (from an index and the wallet's primary seed), and returns it.
func (wallet *WalletClient) NewPublicKey() (tftypes.PublicKey, error) {
	var result api.WalletPublicKeyGET
	err := wallet.client.GetAPI("/wallet/publickey", &result)
	if err != nil {
		return tftypes.PublicKey{}, fmt.Errorf("failed to get (new) public key: %v", err)
	}
	return result.PublicKey, nil
}

// FundCoins collects coin inputs owned by this daemon's wallet,
// that are sufficient to fund the given amount, optionally returning a refund coin output as well.
func (wallet *WalletClient) FundCoins(amount types.Currency) ([]types.CoinInput, *types.CoinOutput, error) {
	var result api.WalletFundCoins
	err := wallet.client.GetAPI(fmt.Sprintf("/wallet/fund/coins?amount=%s", amount.String()), &result)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get fund coins: %v", err)
	}
	return result.CoinInputs, result.RefundCoinOutput, nil
}

// GreedySignTx signs the given transactions greedy,
// meaning that all fulfillments that can be signed, will be signed.
func (wallet *WalletClient) GreedySignTx(t *types.Transaction) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	err = wallet.client.PostResp("/wallet/sign", string(b), t)
	if err != nil {
		return fmt.Errorf("Failed to sign transaction: %v", err)
	}
	return nil
}
