package explorer

import (
	"github.com/threefoldfoundation/tfchain/pkg/config"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/modules"
)

// TestnetGroupedExplorer is a GroupedExplorer preconfigured for the official public testnet explorers
type TestnetGroupedExplorer struct {
	*GroupedExplorer
}

// NewTestnetGroupedExplorer creates a preconfigured grouped explorer for the public testnet nodes
func NewTestnetGroupedExplorer() *TestnetGroupedExplorer {
	testnetUrls := []string{
		"https://explorer.testnet.threefoldtoken.com",
		"https://explorer2.testnet.threefoldtoken.com",
	}
	var explorers []*Explorer
	for _, url := range testnetUrls {
		explorers = append(explorers, NewExplorer(url, "Rivine-Agent", ""))
	}
	explorer := &TestnetGroupedExplorer{NewGroupedExplorer(explorers...)}
	// This call doesn't return an error since it just loads hard coded constants
	cts, _ := explorer.GetChainConstants()
	tftypes.RegisterTransactionTypesForTestNetwork(nil, tftypes.NopERC20TransactionValidator{}, cts.OneCoin, config.GetTestnetDaemonNetworkConfig())
	return explorer
}

// GetChainConstants returns the hardcoded chain constants for testnet. No call is made to the explorers
func (te *TestnetGroupedExplorer) GetChainConstants() (modules.DaemonConstants, error) {
	return modules.NewDaemonConstants(config.GetBlockchainInfo(), config.GetTestnetGenesis()), nil
}

// Name of the backend
func (te *TestnetGroupedExplorer) Name() string {
	return "testnet"
}
