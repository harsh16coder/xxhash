package xxhash

import (
	"encoding/binary"

	"github.com/zeebo/xxh3"
)

type Uint128 struct {
	Hi, Lo uint64
}

// Digest128 implements hash 128
type Digest128 struct {
	buf  []byte
	seed uint64
}

// New128 creates a new Digest128 with seed = 0
func New128() *Digest128 {
	return &Digest128{buf: make([]byte, 0, 256), seed: 0}
}

// New128WithSeed creates a new Digest128 with a given seed
func New128WithSeed(seed uint64) *Digest128 {
	return &Digest128{buf: make([]byte, 0, 256), seed: seed}
}

func (d *Digest128) Write(p []byte) (int, error) {
	d.buf = append(d.buf, p...)
	return len(p), nil
}

func Sum128(b []byte) Uint128 {
	h := xxh3.Hash128(b)
	return Uint128{Hi: h.Hi, Lo: h.Lo}
}

func Sum128WithSeed(b []byte, seed uint64) Uint128 {
	h := xxh3.Hash128Seed(b, seed)
	return Uint128{Hi: h.Hi, Lo: h.Lo}
}

func Sum128String(s string) Uint128 {
	h := xxh3.HashString128(s)
	return Uint128{Hi: h.Hi, Lo: h.Lo}
}

func Sum128StringWithSeed(s string, seed uint64) Uint128 {
	h := xxh3.HashString128Seed(s, seed)
	return Uint128{Hi: h.Hi, Lo: h.Lo}
}

func (d *Digest128) Sum128() Uint128 {
	if d.seed == 0 {
		return Sum128(d.buf)
	}
	return Sum128WithSeed(d.buf, d.seed)
}

func (d *Digest128) Sum(b []byte) []byte {
	h := d.Sum128()
	var out [16]byte
	binary.BigEndian.PutUint64(out[0:8], h.Hi)
	binary.BigEndian.PutUint64(out[8:16], h.Lo)
	return append(b, out[:]...)
}

func (d *Digest128) Reset() {
	d.buf = d.buf[:0]
}
