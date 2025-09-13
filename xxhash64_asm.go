package xxhash

// Sum64 computes the 64-bit xxHash digest of b with a zero seed.
func Sum64(b []byte) uint64 {
	total := uint64(len(b))
	var h uint64

	if len(b) >= 32 {
		var v1 uint64 = prime1
		v1 += prime2
		v2 := prime2
		v3 := uint64(0)
		var v4 uint64 = 0
		v4 -= prime1

		for len(b) >= 32 {
			v1 = round(v1, u64(b[:8]))
			v2 = round(v2, u64(b[8:16]))
			v3 = round(v3, u64(b[16:24]))
			v4 = round(v4, u64(b[24:32]))
			b = b[32:]
		}

		h = rol1(v1) + rol7(v2) + rol12(v3) + rol18(v4)
		h = mergeRound(h, v1)
		h = mergeRound(h, v2)
		h = mergeRound(h, v3)
		h = mergeRound(h, v4)
	} else {
		h = uint64(0) + prime5
	}

	h += total

	for ; len(b) >= 8; b = b[8:] {
		k1 := round(0, u64(b[:8]))
		h ^= k1
		h = rol27(h)*prime1 + prime4
	}
	if len(b) >= 4 {
		h ^= uint64(u32(b[:4])) * prime1
		h = rol23(h)*prime2 + prime3
		b = b[4:]
	}
	for ; len(b) > 0; b = b[1:] {
		h ^= uint64(b[0]) * prime5
		h = rol11(h) * prime1
	}

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return h
}

func writeBlocks(d *Digest, b []byte) int {
	n := 0
	for len(b) >= 32 {
		d.v1 = round(d.v1, u64(b[:8]))
		d.v2 = round(d.v2, u64(b[8:16]))
		d.v3 = round(d.v3, u64(b[16:24]))
		d.v4 = round(d.v4, u64(b[24:32]))
		b = b[32:]
		n += 32
	}
	return n
}

// ...existing code...
