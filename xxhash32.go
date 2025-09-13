// Package xxhash implements the 32-bit variant of xxHash (XXH32) as described
// at https://xxhash.com/.
package xxhash

import (
	"encoding/binary"
	"errors"
	"math/bits"
)

const (
	prime32_1 uint32 = 2654435761
	prime32_2 uint32 = 2246822519
	prime32_3 uint32 = 3266489917
	prime32_4 uint32 = 668265263
	prime32_5 uint32 = 374761393
)

// Store the primes in an array as well.
//
// The consts are used when possible in Go code to avoid MOVs but we need a
// contiguous array for the assembly code.
var primes32 = [...]uint32{prime32_1, prime32_2, prime32_3, prime32_4, prime32_5}

// Digest32 implements hash.Hash32.
//
// Note that a zero-valued Digest32 is not ready to receive writes.
// Call Reset or create a Digest32 using New32 before calling other methods.
type Digest32 struct {
	v1    uint32
	v2    uint32
	v3    uint32
	v4    uint32
	total uint64
	mem   [16]byte
	n     int // how much of mem is used
}

// New32 creates a new Digest32 with a zero seed.
func New32() *Digest32 {
	return New32WithSeed(0)
}

// New32WithSeed creates a new Digest32 with the given seed.
func New32WithSeed(seed uint32) *Digest32 {
	var d Digest32
	d.ResetWithSeed(seed)
	return &d
}

// Reset clears the Digest32's state so that it can be reused.
// It uses a seed value of zero.
func (d *Digest32) Reset() {
	d.ResetWithSeed(0)
}

// ResetWithSeed clears the Digest32's state so that it can be reused.
// It uses the given seed to initialize the state.
func (d *Digest32) ResetWithSeed(seed uint32) {
	d.v1 = seed + prime32_1 + prime32_2
	d.v2 = seed + prime32_2
	d.v3 = seed
	d.v4 = seed - prime32_1
	d.total = 0
	d.n = 0
}

// Size always returns 4 bytes.
func (d *Digest32) Size() int { return 4 }

// BlockSize always returns 16 bytes.
func (d *Digest32) BlockSize() int { return 16 }

// Write adds more data to d. It always returns len(b), nil.
func (d *Digest32) Write(b []byte) (n int, err error) {
	n = len(b)
	d.total += uint64(n)

	memleft := d.mem[d.n:]

	if d.n+n < 16 {
		// This new data doesn't even fill the current block.
		copy(memleft, b)
		d.n += n
		return
	}

	if d.n > 0 {
		// Finish off the partial block.
		c := copy(memleft, b)
		d.v1 = round32(d.v1, u32(d.mem[0:4]))
		d.v2 = round32(d.v2, u32(d.mem[4:8]))
		d.v3 = round32(d.v3, u32(d.mem[8:12]))
		d.v4 = round32(d.v4, u32(d.mem[12:16]))
		b = b[c:]
		d.n = 0
	}

	if len(b) >= 16 {
		// One or more full blocks left.
		nw := writeBlocks32(d, b)
		b = b[nw:]
	}

	// Store any remaining partial block.
	copy(d.mem[:], b)
	d.n = len(b)

	return
}

// Sum appends the current hash to b and returns the resulting slice.
func (d *Digest32) Sum(b []byte) []byte {
	s := d.Sum32()
	var a [4]byte
	binary.BigEndian.PutUint32(a[:], s)
	return append(b, a[:]...)
}

// Sum32 returns the current hash.
func (d *Digest32) Sum32() uint32 {
	var h uint32

	if d.total >= 16 {
		h = rol32_1(d.v1) + rol32_7(d.v2) + rol32_12(d.v3) + rol32_18(d.v4)
	} else {
		h = d.v3 + prime32_5
	}

	h += uint32(d.total)

	b := d.mem[:d.n]
	for ; len(b) >= 4; b = b[4:] {
		h += u32(b[:4]) * prime32_3
		h = rol32_17(h) * prime32_4
	}
	for ; len(b) > 0; b = b[1:] {
		h += uint32(b[0]) * prime32_5
		h = rol32_11(h) * prime32_1
	}

	h ^= h >> 15
	h *= prime32_2
	h ^= h >> 13
	h *= prime32_3
	h ^= h >> 16

	return h
}

const (
	magic32         = "xxh\x03"
	marshaledSize32 = len(magic32) + 4*4 + 8 + 16
)

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d *Digest32) MarshalBinary() ([]byte, error) {
	b := make([]byte, marshaledSize32)
	off := 0
	copy(b[off:], magic32)
	off += len(magic32)
	binary.LittleEndian.PutUint32(b[off:], d.v1)
	off += 4
	binary.LittleEndian.PutUint32(b[off:], d.v2)
	off += 4
	binary.LittleEndian.PutUint32(b[off:], d.v3)
	off += 4
	binary.LittleEndian.PutUint32(b[off:], d.v4)
	off += 4
	binary.LittleEndian.PutUint64(b[off:], d.total)
	off += 8
	copy(b[off:], d.mem[:d.n])
	off += d.n
	// remaining bytes are zeroed by make
	return b, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (d *Digest32) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic32) || string(b[:len(magic32)]) != magic32 {
		return errors.New("xxhash: invalid hash state identifier")
	}
	if len(b) != marshaledSize32 {
		return errors.New("xxhash: invalid hash state size")
	}
	b = b[len(magic32):]
	b, d.v1 = consumeUint32(b)
	b, d.v2 = consumeUint32(b)
	b, d.v3 = consumeUint32(b)
	b, d.v4 = consumeUint32(b)
	b, d.total = consumeUint64(b)
	copy(d.mem[:], b)
	d.n = int(d.total % uint64(len(d.mem)))
	return nil
}

func consumeUint32(b []byte) ([]byte, uint32) {
	x := u32(b)
	return b[4:], x
}

//func u32(b []byte) uint32 { return binary.LittleEndian.Uint32(b) }

func round32(acc, input uint32) uint32 {
	acc += input * prime32_2
	acc = rol32_13(acc)
	acc *= prime32_1
	return acc
}

func rol32_1(x uint32) uint32  { return bits.RotateLeft32(x, 1) }
func rol32_7(x uint32) uint32  { return bits.RotateLeft32(x, 7) }
func rol32_11(x uint32) uint32 { return bits.RotateLeft32(x, 11) }
func rol32_12(x uint32) uint32 { return bits.RotateLeft32(x, 12) }
func rol32_13(x uint32) uint32 { return bits.RotateLeft32(x, 13) }
func rol32_17(x uint32) uint32 { return bits.RotateLeft32(x, 17) }
func rol32_18(x uint32) uint32 { return bits.RotateLeft32(x, 18) }
