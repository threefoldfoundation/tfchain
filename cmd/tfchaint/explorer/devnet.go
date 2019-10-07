package explorer

import (
	tfcli "github.com/threefoldfoundation/tfchain/extensions/tfchain/client"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/pkg/client"
)

// DevnetGroupedExplorer is a GroupedExplorer preconfigured for the official public testnet explorers
type DevnetGroupedExplorer struct {
	*GroupedExplorer
}

// NewDevnetGroupedExplorer creates a preconfigured grouped explorer for the public testnet nodes
func NewDevnetGroupedExplorer() *DevnetGroupedExplorer {
	devnetURL := []string{
		"http://localhost:2015",
	}
	var explorers []*Explorer
	for _, url := range devnetURL {
		explorers = append(explorers, NewExplorer(url, "Rivine-Agent", ""))
	}
	explorer := &DevnetGroupedExplorer{NewGroupedExplorer(explorers...)}

	// register transactions for development network of tfchain
	bc, err := client.NewBaseClient(explorer, nil)
	if err != nil {
		panic(err)
	}
	tfcli.RegisterDevnetTransactions(bc)

	return explorer
}

// GetChainConstants returns the hardcoded chain constants for devnet. No call is made to the explorers
func (te *DevnetGroupedExplorer) GetChainConstants() (modules.DaemonConstants, error) {
	return modules.NewDaemonConstants(config.GetBlockchainInfo(), config.GetDevnetGenesis()), nil
}

// Name of the backend
func (te *DevnetGroupedExplorer) Name() string {
	return "devnet"
}
