package explorer

import (
	"errors"
	"net"

	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/types"
)

var (
	// ErrNoHealthyExplorers is returned if all explorers in the group fail to respond in time
	ErrNoHealthyExplorers = errors.New("No explorer could statisfy the request")
)

// GroupedExplorer is a Backend which can call multiple explorers, calling another explorer if one is down
type GroupedExplorer struct {
	explorers []*Explorer
}

// NewGroupedExplorer creates a new GroupedExplorer from existing regular Explorers
func NewGroupedExplorer(explorers ...*Explorer) *GroupedExplorer {
	return &GroupedExplorer{explorers: explorers}
}

// NewTestnetGroupedExplorer creates a preconfigured grouped explorer for the public testnet nodes
func NewTestnetGroupedExplorer() *GroupedExplorer {
	testnetUrls := []string{"https://explorer.testnet.threefoldtoken.com", "https://explorer2.testnet.threefoldtoken.com"}
	var explorers []*Explorer
	for _, url := range testnetUrls {
		explorers = append(explorers, NewExplorer(url, "Rivine-Agent", ""))
	}
	return NewGroupedExplorer(explorers...)
}

// CheckAddress returns all interesting transactions and blocks related to a given unlockhash
func (e *GroupedExplorer) CheckAddress(addr types.UnlockHash) ([]api.ExplorerBlock, []api.ExplorerTransaction, error) {
	for _, explorer := range e.explorers {
		blocks, transactions, err := explorer.CheckAddress(addr)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			continue
		}
		return blocks, transactions, err
	}
	return nil, nil, ErrNoHealthyExplorers
}

// CurrentHeight returns the current chain height
func (e *GroupedExplorer) CurrentHeight() (types.BlockHeight, error) {
	for _, explorer := range e.explorers {
		height, err := explorer.CurrentHeight()
		if err, ok := err.(net.Error); ok && err.Timeout() {
			continue
		}
		return height, err
	}
	return 0, ErrNoHealthyExplorers
}

// SendTxn sends a txn to the backend to ultimately include it in the transactionpool
func (e *GroupedExplorer) SendTxn(tx types.Transaction) error {
	for _, explorer := range e.explorers {
		err := explorer.SendTxn(tx)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			continue
		}
		return err
	}
	return ErrNoHealthyExplorers
}

// GetChainConstants gets the currently active chain constants for this backend
func (e *GroupedExplorer) GetChainConstants() (modules.DaemonConstants, error) {
	for _, explorer := range e.explorers {
		cts, err := explorer.GetChainConstants()
		if err, ok := err.(net.Error); ok && err.Timeout() {
			continue
		}
		return cts, err
	}
	return modules.DaemonConstants{}, ErrNoHealthyExplorers
}
