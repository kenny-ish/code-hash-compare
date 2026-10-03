# code-hash-compare

Checks whether an address has the same bytecode on several EVM chains. Protocols often deploy to the
same address everywhere (CREATE2, deterministic deployers), but the same address doesn't guarantee
the same code.

```bash
go run . 0xcA11bde05977b3631167028862bE2a173976CA11     # Multicall3
go run . -chains eth,base,bsc 0x000000000022D473030F116dDEE9F6B43aC78BA3
```

For each chain it prints the code size and the start of the code's SHA-256 hash, then says whether
all chains match. SHA-256 is used because Go's standard library has no keccak. Only equality matters
here, so it doesn't need to match the on-chain code hash.

```bash
go test ./...
```
