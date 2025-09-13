package xxhash

import (
	"github.com/zeebo/xxh3"
)

type Uint128 struct {
	Hi, Lo uint64
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
