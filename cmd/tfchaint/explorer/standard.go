package explorer

import (
	"github.com/threefoldfoundation/tfchain/pkg/config"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/modules"
)

// MainnetGroupedExplorer is a GroupedExplorer preconfigured for the official public testnet explorers
type MainnetGroupedExplorer struct {
	*GroupedExplorer
}

// NewMainnetGroupedExplorer creates a preconfigured grouped explorer for the public testnet nodes
func NewMainnetGroupedExplorer() *MainnetGroupedExplorer {
	mainnetUrls := []string{
		"https://explorer.threefoldtoken.com",
		"https://explorer2.threefoldtoken.com",
		"https://explorer3.threefoldtoken.com",
		"https://explorer4.threefoldtoken.com",
	}
	var explorers []*Explorer
	for _, url := range mainnetUrls {
		explorers = append(explorers, NewExplorer(url, "Rivine-Agent", ""))
	}
	explorer := &MainnetGroupedExplorer{NewGroupedExplorer(explorers...)}
	// This call doesn't return an error since it just loads hard coded constants
	cts, _ := explorer.GetChainConstants()
	tftypes.RegisterTransactionTypesForStandardNetwork(nil, tftypes.NopERC20TransactionValidator{}, cts.OneCoin, config.GetStandardDaemonNetworkConfig())
	return explorer
}

// GetChainConstants returns the hardcoded chain constants for mainnet. No call is made to the explorers
func (te *MainnetGroupedExplorer) GetChainConstants() (modules.DaemonConstants, error) {
	return modules.NewDaemonConstants(config.GetBlockchainInfo(), config.GetStandardnetGenesis()), nil
}

// Name of the backend
func (te *MainnetGroupedExplorer) Name() string {
	return "standard"
}
