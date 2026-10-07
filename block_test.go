package chain

import (
	"bytes"
	"testing"
)

func TestGenesis(t *testing.T) {
	g := Genesis()

	if g.Index != 0 {
		t.Errorf("genesis index = %d, want 0", g.Index)
	}
	if g.PrevHash != (Hash{}) {
		t.Errorf("genesis prevHash = %x, want all zeros", g.PrevHash)
	}
	if g.Hash != g.ComputeHash() {
		t.Errorf("genesis hash = %x, want %x", g.Hash, g.ComputeHash())
	}
	if Genesis().Hash != g.Hash {
		t.Error("Genesis() is not deterministic: two calls produced different hashes")
	}
}

func TestComputeHash(t *testing.T) {
	a := Block{Index: 1, Timestamp: 1700000000, Data: []byte("hello"), PrevHash: Genesis().Hash}
	same := a
	same.Data = bytes.Clone(a.Data)

	if a.ComputeHash() != same.ComputeHash() {
		t.Error("identical contents produced different hashes")
	}

	changed := a
	changed.Data = []byte("hellp")
	if a.ComputeHash() == changed.ComputeHash() {
		t.Error("different data produced the same hash")
	}
}

func TestNewBlockLinksToPrevious(t *testing.T) {
	g := Genesis()
	b := NewBlock(g, []byte("payload"))

	if b.Index != g.Index+1 {
		t.Errorf("index = %d, want %d", b.Index, g.Index+1)
	}
	if b.PrevHash != g.Hash {
		t.Errorf("prevHash = %x, want %x", b.PrevHash, g.Hash)
	}
	if !bytes.Equal(b.Data, []byte("payload")) {
		t.Errorf("data = %q, want %q", b.Data, "payload")
	}
	if b.Hash != b.ComputeHash() {
		t.Errorf("hash = %x, want %x", b.Hash, b.ComputeHash())
	}
}
