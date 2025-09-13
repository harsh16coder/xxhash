package xxhash

// Sum32 computes the 32-bit xxHash digest of b with a zero seed.
func Sum32(b []byte) uint32 {
	total := uint64(len(b))
	var h uint32

	if len(b) >= 16 {
		var v1 uint32 = prime32_1
		v1 += prime32_2
		v2 := prime32_2
		v3 := uint32(0)
		var v4 uint32 = 0
		v4 -= prime32_1
		for len(b) >= 16 {
			v1 = round32(v1, u32(b[:4]))
			v2 = round32(v2, u32(b[4:8]))
			v3 = round32(v3, u32(b[8:12]))
			v4 = round32(v4, u32(b[12:16]))
			b = b[16:]
		}

		h = rol32_1(v1) + rol32_7(v2) + rol32_12(v3) + rol32_18(v4)
	} else {
		h = uint32(0) + prime32_5
	}

	h += uint32(total)

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

func writeBlocks32(d *Digest32, b []byte) int {
	n := 0
	for len(b) >= 16 {
		d.v1 = round32(d.v1, u32(b[:4]))
		d.v2 = round32(d.v2, u32(b[4:8]))
		d.v3 = round32(d.v3, u32(b[8:12]))
		d.v4 = round32(d.v4, u32(b[12:16]))
		b = b[16:]
		n += 16
	}
	return n
}
