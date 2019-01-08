package erc20

import (
	"fmt"
	"testing"

	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/types"
)

func Test_blockBuffer_new(t *testing.T) {
	size := 6
	buf := newBlockBuffer(uint(size))
	if buf.size != 6 {
		t.Errorf("BlockBuffer new: Invalid size, got %v, want %v", buf.size, size)
	}
	if len(buf.blocks) != size {
		t.Errorf("BlockBuffer new: Invalid internal slice length, got %v, want %v", len(buf.blocks), size)
	}
	if (buf.current) != 0 {
		t.Errorf("BlockBuffer new: Invalid initialization of pointer to current element, got %v, want 0", buf.current)
	}
}

func Test_blockBuffer_pushBlock(t *testing.T) {
	for size := 1; size < 11; size++ {
		t.Run(fmt.Sprintf("blockbuffer_add_%v", size), func(t *testing.T) {
			buf := newBlockBuffer(uint(size))
			for i := 0; i < size*2; i++ {
				res := buf.pushBlock(types.Block{}, modules.ConsensusChangeID{})
				if i < size && res != nil {
					t.Errorf("BlockBuffer pushBlock: expected nil response, got %v", res)
				}
				if i >= size && res == nil {
					t.Error("BlockBuffer pushBlock: expected response, got nil")
				}
			}
			for idx, el := range buf.blocks {
				if el == nil {
					t.Errorf("Expected to find element in buffer at index %v, got nil", idx)
				}
			}
			if buf.current != 0 {
				t.Errorf("Expected pointer to current element to be at 0, but is at %v instead", buf.current)
			}
		})
	}
}

func Test_blockBuffer_rewindBlock(t *testing.T) {
	for size := 1; size < 11; size++ {
		t.Run(fmt.Sprintf("blockbuffer_rewind_%v", size), func(t *testing.T) {
			buf := newBlockBuffer(uint(size))
			// fill buffer
			for range buf.blocks {
				buf.pushBlock(types.Block{}, modules.ConsensusChangeID{})
			}
			for i := 0; i < size*2; i++ {
				buf.rewindBlock()
			}
			for idx, el := range buf.blocks {
				if el != nil {
					t.Errorf("Expected buffer to be empty, but found something at index %v", idx)
				}
			}
			if buf.current != 0 {
				t.Errorf("Expected pointer to current element to be at 0, but is at %v instead", buf.current)
			}
			// fill buffer
			for range buf.blocks {
				buf.pushBlock(types.Block{}, modules.ConsensusChangeID{})
			}
			removals := size / 2
			for i := 0; i < removals; i++ {
				buf.rewindBlock()
			}
			for idx, el := range buf.blocks {
				if idx < size-removals && el == nil {
					t.Errorf("Expected to find something at index %v, but got nil", idx)
				}
				if idx >= size-removals && el != nil {
					t.Errorf("Expected element at idx %v to be nil, but found something", idx)
				}
			}
			// Current needs to have the index of the next slot that will be used.
			// EDGE CASE: single element buffers always point to 0
			if buf.size == 1 {
				if buf.current != 0 {
					t.Errorf("Expected pointer to current element to be at %v, but is at %v", 0, buf.current)
				}
				return
			}
			if int(buf.current) != size-removals {
				t.Log(size, removals, size-removals, buf.current)
				t.Errorf("Expected pointer to current element to be at %v, but is at %v", size-removals, buf.current)
			}
		})
	}
}
