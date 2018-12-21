package erc20

import (
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/types"
)

const (
	// blockDelay is the amount of blocks to wait before
	// pushing tft transactions to the ethereum contract.
	blockDelay = 6
)

type (
	// blockBuffer holds an amount of blocks
	blockBuffer struct {
		blocks  []*bufferedBlock
		current uint
		size    uint
	}

	// bufferedBlock holds a block and the consensus change ID at which
	// it was received
	bufferedBlock struct {
		types.Block
		modules.ConsensusChangeID
	}
)

func newBlockBuffer(size uint) *blockBuffer {
	return &blockBuffer{
		blocks:  make([]*bufferedBlock, size),
		current: 0,
		size:    size,
	}
}

// pushBlock adds a new block to the buffer, and returns the block previously there
func (buf *blockBuffer) pushBlock(block types.Block, ccid modules.ConsensusChangeID) *bufferedBlock {
	// coppy out the pointer to the old block
	oldBlock := buf.blocks[buf.current]
	// insert new block
	buf.blocks[buf.current] = &bufferedBlock{block, ccid}
	// move pointer to current block to the next slot
	buf.current = (buf.current + 1) % buf.size

	return oldBlock
}

// rewindBlock removes the latest block from the buffer
func (buf *blockBuffer) rewindBlock() {
	// push back pointer to point to latest inserted block
	buf.current = (buf.current + buf.size - 1) % buf.size
	// remove latest block
	buf.blocks[buf.current] = nil
}
