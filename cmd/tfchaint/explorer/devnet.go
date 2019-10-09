package explorer

import (
	tfcli "github.com/threefoldfoundation/tfchain/extensions/tfchain/client"
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

// Name of the backend
func (te *DevnetGroupedExplorer) Name() string {
	return "devnet"
}
