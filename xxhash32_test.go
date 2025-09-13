package xxhash

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSum32_AgainstDigest(t *testing.T) {
	cases := []struct {
		name string
		b    []byte
	}{
		{"empty", []byte{}},
		{"a", []byte("a")},
		{"hello world", []byte("hello world")},
		{"short-seq", func() []byte {
			b := make([]byte, 8)
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
			got := Sum32(c.b)
			var d Digest32
			d.Reset()
			if _, err := d.Write(c.b); err != nil {
				t.Fatalf("Write failed: %v", err)
			}
			want := d.Sum32()
			if got != want {
				t.Fatalf("Sum32 mismatch for %s: got 0x%x want 0x%x (len=%d)", c.name, got, want, len(c.b))
			}
		})
	}
}

func TestSum32_IncrementalEquivalence(t *testing.T) {
	// prepare a deterministic buffer
	buf := make([]byte, 1024)
	for i := range buf {
		buf[i] = byte(i)
	}
	base := Sum32(buf)

	splits := []int{1, 2, 3, 4, 7, 8, 15, 16, 17, 32, 64, 128}
	for _, split := range splits {
		split := split
		t.Run(fmt.Sprintf("split=%d", split), func(t *testing.T) {
			var d Digest32
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
			if got := d.Sum32(); got != base {
				t.Fatalf("incremental Sum32 mismatch for split %d: got 0x%x want 0x%x", split, got, base)
			}
		})
	}
}

func TestSum32_Consistency(t *testing.T) {
	b := []byte("consistency-test-bytes")
	v1 := Sum32(b)
	v2 := Sum32(b)
	if v1 != v2 {
		t.Fatalf("Sum32 inconsistent: first 0x%x second 0x%x", v1, v2)
	}
}

func TestSum32_KnownValues(t *testing.T) {
	// Test against known XXH32 values
	tests := []struct {
		input    string
		expected uint32
	}{
		{"", 0x02cc5d05},
		{"a", 0x550d7456},
		{"abc", 0x32d153ff},
		{"message digest", 0x7c948494},
		{"abcdefghijklmnopqrstuvwxyz", 0x63a14d5f},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 0x9c285e64},
		{"12345678901234567890123456789012345678901234567890123456789012345678901234567890", 0x9c05f475},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Sum32([]byte(tt.input))
			if got != tt.expected {
				t.Errorf("Sum32(%q) = 0x%08x, want 0x%08x", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintHash32(t *testing.T) {
	t.Log(Sum32([]byte("harsh")))
	t.Log(Sum32([]byte("")))
}
