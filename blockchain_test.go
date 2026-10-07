package chain

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	bc := New()

	if len(bc.Blocks) != 1 {
		t.Fatalf("new chain has %d blocks, want 1", len(bc.Blocks))
	}
	if bc.Blocks[0].Hash != Genesis().Hash {
		t.Error("first block is not the genesis block")
	}
	if err := bc.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestAddBlock(t *testing.T) {
	bc := New()
	b := bc.AddBlock([]byte("first"))

	if len(bc.Blocks) != 2 {
		t.Fatalf("chain has %d blocks, want 2", len(bc.Blocks))
	}
	if bc.Last().Hash != b.Hash {
		t.Error("Last() is not the block just added")
	}
	if b.Index != 1 {
		t.Errorf("index = %d, want 1", b.Index)
	}
	if b.PrevHash != Genesis().Hash {
		t.Error("new block does not link to genesis")
	}
	if err := bc.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestAddSeveralBlocks(t *testing.T) {
	bc := New()
	for i := 1; i <= 5; i++ {
		bc.AddBlock(fmt.Appendf(nil, "block %d", i))
	}

	if len(bc.Blocks) != 6 {
		t.Fatalf("chain has %d blocks, want 6", len(bc.Blocks))
	}
	for i := 1; i < len(bc.Blocks); i++ {
		prev, cur := bc.Blocks[i-1], bc.Blocks[i]
		if cur.Index != uint64(i) {
			t.Errorf("block %d has index %d", i, cur.Index)
		}
		if cur.PrevHash != prev.Hash {
			t.Errorf("block %d does not link to block %d", i, i-1)
		}
	}
	if err := bc.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

// newTestChain returns a valid chain with genesis plus n blocks.
func newTestChain(t *testing.T, n int) *Blockchain {
	t.Helper()
	bc := New()
	for i := 1; i <= n; i++ {
		bc.AddBlock(fmt.Appendf(nil, "block %d", i))
	}
	if err := bc.Validate(); err != nil {
		t.Fatalf("test chain is not valid: %v", err)
	}
	return bc
}

func TestValidateDetectsModifiedData(t *testing.T) {
	bc := newTestChain(t, 3)
	bc.Blocks[1].Data = []byte("tampered")

	err := bc.Validate()
	if !errors.Is(err, ErrInvalidHash) {
		t.Errorf("Validate() = %v, want ErrInvalidHash", err)
	}
}

func TestValidateDetectsModifiedHash(t *testing.T) {
	bc := newTestChain(t, 3)
	bc.Blocks[2].Hash[0] ^= 0xff

	err := bc.Validate()
	if !errors.Is(err, ErrInvalidHash) {
		t.Errorf("Validate() = %v, want ErrInvalidHash", err)
	}
}

// Rewriting a block and recomputing its own hash makes the block internally
// consistent, but the next block still points to the old hash.
func TestValidateDetectsBrokenLink(t *testing.T) {
	bc := newTestChain(t, 3)
	bc.Blocks[1].Data = []byte("rewritten history")
	bc.Blocks[1].Hash = bc.Blocks[1].ComputeHash()

	err := bc.Validate()
	if !errors.Is(err, ErrBrokenLink) {
		t.Errorf("Validate() = %v, want ErrBrokenLink", err)
	}
}

func TestValidateDetectsRemovedBlock(t *testing.T) {
	bc := newTestChain(t, 3)
	bc.Blocks = append(bc.Blocks[:2], bc.Blocks[3:]...)

	err := bc.Validate()
	if !errors.Is(err, ErrInvalidIndex) {
		t.Errorf("Validate() = %v, want ErrInvalidIndex", err)
	}
}

func TestValidateRejectsForeignGenesis(t *testing.T) {
	other := newBlock(0, 0, []byte("another network"), Hash{})
	bc := &Blockchain{Blocks: []Block{other}}
	bc.AddBlock([]byte("looks fine on its own"))

	err := bc.Validate()
	if !errors.Is(err, ErrInvalidGenesis) {
		t.Errorf("Validate() = %v, want ErrInvalidGenesis", err)
	}
}

func TestValidateRejectsEmptyChain(t *testing.T) {
	bc := &Blockchain{}
	if err := bc.Validate(); !errors.Is(err, ErrEmptyChain) {
		t.Errorf("Validate() = %v, want ErrEmptyChain", err)
	}
}
