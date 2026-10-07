package chain

import (
	"errors"
	"fmt"
)

// Validation errors. Validate wraps them with the index of the offending
// block, so callers should test them with errors.Is.
var (
	ErrEmptyChain     = errors.New("chain has no blocks")
	ErrInvalidGenesis = errors.New("first block is not the genesis block")
	ErrInvalidIndex   = errors.New("block index does not match its position")
	ErrInvalidHash    = errors.New("block hash does not match its contents")
	ErrBrokenLink     = errors.New("previous hash does not match the previous block")
)

// Blockchain is an append-only sequence of blocks starting at Genesis.
type Blockchain struct {
	Blocks []Block
}

// New returns a chain containing only the genesis block.
func New() *Blockchain {
	return &Blockchain{Blocks: []Block{Genesis()}}
}

// Last returns the most recent block. The chain always has at least one.
func (bc *Blockchain) Last() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

// AddBlock appends a new block carrying data and returns it.
func (bc *Blockchain) AddBlock(data []byte) Block {
	b := NewBlock(bc.Last(), data)
	bc.Blocks = append(bc.Blocks, b)
	return b
}

// Validate checks the integrity of the whole chain and returns the first
// problem found:
//
//   - the chain starts with the known genesis block;
//   - each block's Index equals its position;
//   - each block's Hash equals ComputeHash() of its contents;
//   - each block's PrevHash equals the Hash of the block before it.
//
// Recomputing every hash is what makes tampering detectable: altering a
// block's contents invalidates its Hash, and fixing that Hash breaks the
// PrevHash link of the next block, all the way to the tip.
func (bc *Blockchain) Validate() error {
	if len(bc.Blocks) == 0 {
		return ErrEmptyChain
	}
	if bc.Blocks[0].Hash != Genesis().Hash {
		return fmt.Errorf("block 0: %w", ErrInvalidGenesis)
	}
	for i, b := range bc.Blocks {
		if b.Index != uint64(i) {
			return fmt.Errorf("block %d: %w (index %d)", i, ErrInvalidIndex, b.Index)
		}
		if b.Hash != b.ComputeHash() {
			return fmt.Errorf("block %d: %w", i, ErrInvalidHash)
		}
		if i > 0 && b.PrevHash != bc.Blocks[i-1].Hash {
			return fmt.Errorf("block %d: %w", i, ErrBrokenLink)
		}
	}
	return nil
}
