package xxhash

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSum64_AgainstDigest(t *testing.T) {
	cases := []struct {
		name string
		b    []byte
	}{
		{"empty", []byte{}},
		{"a", []byte("a")},
		{"hello world", []byte("hello world")},
		{"short-seq", func() []byte {
			b := make([]byte, 16)
			for i := range b {
				b[i] = byte(i)
			}
			return b
		}()},
		{"long-repeat", bytes.Repeat([]byte{0xAA}, 1024)},
		{"0..255", func() []byte {
			b := make([]byte, 256)
			for i := 0; i < 256; i++ {
				b[i] = byte(i)
			}
			return b
		}()},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := Sum64(c.b)
			var d Digest64
			d.Reset()
			if _, err := d.Write(c.b); err != nil {
				t.Fatalf("Write failed: %v", err)
			}
			want := d.Sum64()
			if got != want {
				t.Fatalf("Sum64 mismatch for %s: got 0x%x want 0x%x (len=%d)", c.name, got, want, len(c.b))
			}
		})
	}
}

func TestSum64_IncrementalEquivalence(t *testing.T) {
	// prepare a deterministic buffer
	buf := make([]byte, 1024)
	for i := range buf {
		buf[i] = byte(i)
	}
	base := Sum64(buf)

	splits := []int{1, 2, 3, 4, 7, 8, 15, 16, 31, 32, 33, 64, 128}
	for _, split := range splits {
		split := split
		t.Run(fmt.Sprintf("split=%d", split), func(t *testing.T) {
			var d Digest64
			d.Reset()
			for i := 0; i < len(buf); i += split {
				end := i + split
				if end > len(buf) {
					end = len(buf)
				}
				if _, err := d.Write(buf[i:end]); err != nil {
					t.Fatalf("Write failed: %v", err)
				}
			}
			if got := d.Sum64(); got != base {
				t.Fatalf("incremental Sum64 mismatch for split %d: got 0x%x want 0x%x", split, got, base)
			}
		})
	}
}

func TestSum64_Consistency(t *testing.T) {
	b := []byte("consistency-test-bytes")
	v1 := Sum64(b)
	v2 := Sum64(b)
	if v1 != v2 {
		t.Fatalf("Sum64 inconsistent: first 0x%x second 0x%x", v1, v2)
	}
}

func TestPrintHash64(t *testing.T) {
	t.Log(Sum64([]byte("harsh")))
	t.Log(Sum64([]byte("")))
	// Sum64([]byte("harsh"))
}
