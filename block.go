// Package chain implements a small, educational blockchain.
//
// Stage 1: a local chain of hash-linked blocks with integrity validation.
package chain

import (
	"crypto/sha256"
	"encoding/binary"
	"time"
)

// Hash is a SHA-256 digest. A fixed-size array (not a slice) can be compared
// with == and has no nil/empty ambiguity when it is fed into the next hash.
type Hash [sha256.Size]byte

// Block is one element of the chain. Every field except Hash is covered by
// the hash, so changing any of them changes Hash.
type Block struct {
	Index     uint64 // position in the chain, genesis is 0
	Timestamp int64  // Unix time in seconds
	Data      []byte // arbitrary payload for now
	PrevHash  Hash   // hash of the previous block, all zeros for genesis
	Hash      Hash   // ComputeHash() of this block
}

// genesisTime is fixed so that every participant derives exactly the same
// genesis block. If it depended on time.Now(), two chains would never share
// a common root.
var genesisTime = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

// Genesis returns the first block of the chain. It is deterministic: calling
// it twice yields identical blocks and hashes.
func Genesis() Block {
	return newBlock(0, genesisTime.Unix(), []byte("genesis"), Hash{})
}

// NewBlock creates the block that follows prev, carrying data, stamped with
// the current time, and with its hash already computed.
func NewBlock(prev Block, data []byte) Block {
	return newBlock(prev.Index+1, time.Now().Unix(), data, prev.Hash)
}

func newBlock(index uint64, timestamp int64, data []byte, prevHash Hash) Block {
	b := Block{Index: index, Timestamp: timestamp, Data: data, PrevHash: prevHash}
	b.Hash = b.ComputeHash()
	return b
}

// ComputeHash returns the SHA-256 of the block contents (everything except
// the Hash field itself).
//
// The fields are serialized in a fixed order with fixed-width encodings
// (8-byte big-endian integers, the 32-byte previous hash) followed by the
// variable-length Data. Because every field before Data has a known size, the
// encoding is unambiguous: two different blocks can never serialize to the
// same bytes. Plain string concatenation would not give that guarantee
// ("1"+"23" and "12"+"3" collide).
func (b Block) ComputeHash() Hash {
	h := sha256.New()

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], b.Index)
	h.Write(buf[:])
	binary.BigEndian.PutUint64(buf[:], uint64(b.Timestamp))
	h.Write(buf[:])
	h.Write(b.PrevHash[:])
	h.Write(b.Data)

	var out Hash
	h.Sum(out[:0])
	return out
}
