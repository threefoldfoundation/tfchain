package main

import (
	"fmt"

	"github.com/threefoldfoundation/tfchain/pkg/persist"
	"github.com/threefoldtech/rivine/modules"
	rivinetypes "github.com/threefoldtech/rivine/types"
)

func NewBridged(cs modules.ConsensusSet, txdb *persist.TransactionDB, bcInfo rivinetypes.BlockchainInfo, chainCts rivinetypes.ChainConstants, cancel <-chan struct{}) (*Bridged, error) {

	bridged := &Bridged{
		cs:       cs,
		txdb:     txdb,
		bcInfo:   bcInfo,
		chainCts: chainCts,
	}
	err := cs.ConsensusSetSubscribe(bridged, txdb.GetLastConsensusChangeID(), cancel)
	if err != nil {
		return nil, fmt.Errorf("bridged: failed to subscribe to consensus set: %v", err)
	}
	return bridged, nil
}

// Close bridged.
func (bridged *Bridged) Close() {
	bridged.mut.Lock()
	defer bridged.mut.Unlock()
	bridged.cs.Unsubscribe(bridged)
}

// ProcessConsensusChange implements modules.ConsensusSetSubscriber,
// used to apply/revert blocks.
func (bridged *Bridged) ProcessConsensusChange(css modules.ConsensusChange) {
	bridged.mut.Lock()
	defer bridged.mut.Unlock()
	// FIXME: how to get the current height? is this the correct way?
	currentHeight := bridged.cs.Height()

	blocks := css.RevertedBlocks
	blocks = append(blocks, css.AppliedBlocks...)

	for _, block := range blocks {
		fmt.Println("BLOCK : ", block)
		height, _ := bridged.cs.BlockHeightOfBlock(block)
		fmt.Println("HEIGHT: ", height)
		fmt.Println("Current height: ", currentHeight)
		// the block we're interested in shouldn't exist.
		if height-6 == currentHeight {
			// CODE HERE FOR to create the erc20 tokens or register a withdrawal address
			// And should return afterwards.
			fmt.Println("Differs by 6.")
		}
	}
}
