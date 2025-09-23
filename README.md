# xxhash (Go)

This package provides Go implementations / wrappers for the xxHash family:

- XXH32 — 32-bit variant (xxhash32.go, xxhash32_asm.go)
- XXH64 — 64-bit variant (xxhash64.go, xxhash64_asm.go)
- XXH128 (XXH3/128) — wrapped from zeebo/xxh3 (xxhash128.go)

All code lives under the `xxhash` package. Tests are provided for validating correctness and incremental behavior.

## Motivation

The main reason I created this repository is that the xxHash family is currently fragmented:

- `xxh32` and `xxh64` are not available together with `xxh128` in a single Go package.
- To use them today, developers often rely on multiple repositories, each with its own API style and conventions.

My goal with this project is to **unify all major variants (XXH32, XXH64, and XXH128)** into one consistent implementation, so they can be imported and used from a single package.

The implementation here is thin and intentionally simple — it doesn’t diverge from the upstream algorithms but instead exposes them together under one roof with a coherent API.

## Import

Use the package path for your project. Example (local repo):

```go
import "github.com/harsh16coder/xxhash"
```

## Provided APIs

### XXH32 (32-bit)

Convenience:
- Sum32(b []byte) uint32

Incremental (hash.Hash32-like):
- type Digest32
  - New32() *Digest32
  - New32WithSeed(seed uint32) *Digest32
  - (*Digest32) Reset()
  - (*Digest32) ResetWithSeed(seed uint32)
  - (*Digest32) Write(b []byte) (int, error)
  - (*Digest32) Sum(b []byte) []byte           // appends big-endian 4 bytes
  - (*Digest32) Sum32() uint32
  - MarshalBinary / UnmarshalBinary

Example:

```go
var h uint32 = xxhash.Sum32([]byte("hello"))
// Incremental:
d := xxhash.New32()
d.Write([]byte("hello "))
d.Write([]byte("world"))
sum := d.Sum32()
```

### XXH64 (64-bit)

Convenience:
- Sum64(b []byte) uint64

Incremental (hash.Hash64-like):
- type Digest64
  - New64() *Digest64
  - NewWithSeed64(seed uint64) *Digest64
  - (*Digest64) Reset()
  - (*Digest64) ResetWithSeed(seed uint64)
  - (*Digest64) Write(b []byte) (int, error)
  - (*Digest64) Sum(b []byte) []byte           // appends big-endian 8 bytes
  - (*Digest64) Sum64() uint64
  - MarshalBinary / UnmarshalBinary

Example:

```go
// One-shot
sum64 := xxhash.Sum64([]byte("some data"))

// Incremental
d := xxhash.NewWithSeed64(0x1234)
d.Write([]byte("part1"))
d.Write([]byte("part2"))
final := d.Sum64()
```

Notes:
- Digest64 stores internal block state and can be marshaled/unmarshaled to persist state.
- Use `Sum(b []byte)` when you want the 8-byte big-endian representation appended to a slice.

### XXH128 / XXH3 (128-bit)

This package wraps the `github.com/zeebo/xxh3` implementation to provide convenient functions:

- type Uint128 { Hi, Lo uint64 }
- Sum128(b []byte) Uint128
- Sum128WithSeed(b []byte, seed uint64) Uint128
- Sum128String(s string) Uint128
- Sum128StringWithSeed(s string, seed uint64) Uint128

`Uint128` contains two uint64 words; to format as canonical 128-bit hex use:

```go
h := xxhash.Sum128([]byte("hello"))
hex := fmt.Sprintf("%016x%016x", h.Hi, h.Lo)
```

Use the `WithSeed` variants to provide a non-zero seed.

## Testing

From the repository root (Windows Terminal / PowerShell):

```powershell
cd "C:\username\projectpath\xxhash\"
go test -v
```

Run a single package or test with `go test ./xxhash -run TestName -v`.

The test suite validates:
- One-shot vs incremental equivalence
- Deterministic outputs
- Known test vectors (where provided)

## Thread-safety

- The one-shot convenience functions (`Sum32`, `Sum64`, `Sum128`, etc.) are safe to call concurrently.
- Digest types (`Digest32`, `Digest64`) are stateful and not goroutine-safe — protect with your own synchronization if accessed by multiple goroutines.

## Notes / Caveats

- The XXH128 code in this repo wraps `zeebo/xxh3` for correctness and performance. The high/low word order is exposed as `Hi` and `Lo`. Format accordingly when producing canonical hex.
- The package implements marshaling for Digest32/Digest64 to persist intermediate state (useful for streaming workflows).
- For production use, prefer the convenience functions for simple hashing and the Digest types for incremental hashing across blocks.

# Contribution

This project welcomes your PR and issues. For example, refactoring, adding features, correcting English, etc.

If you need any help, you can contact me on [X](https://X.com/harsh1614).

Thanks to all the people who already contributed!

<a href="https://github.com/harsh16coder/xxhash/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=harsh16coder/xxhash" />
</a>

## License
[MIT](./LICENSE)