package explorer

import (
	tfcli "github.com/threefoldfoundation/tfchain/extensions/tfchain/client"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/pkg/client"
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

	// register transactions for testnet network of tfchain
	bc, err := client.NewBaseClient(explorer, nil)
	if err != nil {
		panic(err)
	}
	tfcli.RegisterTestnetTransactions(bc)

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
