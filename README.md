# chain

A small, educational blockchain written in Go with the standard library only.
The goal is to understand each building block of a blockchain by implementing
it from scratch, one stage at a time. It is not a cryptocurrency.

## Current stage: 1 — local chain

- `Block` with index, timestamp, payload, previous hash and its own SHA-256 hash.
- Deterministic genesis block.
- `Blockchain` that appends blocks and validates the integrity of the whole chain
  (content tampering, wrong hashes, broken links, missing blocks, foreign genesis).

```bash
go test ./...
```

## Roadmap

Proof-of-work, transactions, wallets and signatures, mempool, persistence,
nodes and P2P propagation, chain conflict resolution, CLI.
