package xxhash

import (
	"fmt"
	"testing"
)

func TestSum128ActualValues(t *testing.T) {
	tests := []struct {
		input []byte
		seed  uint64
		want  string
	}{
		{
			input: []byte(""),
			seed:  0,
			want:  "99aa06d3014798d86001c324468d497f",
		},
		{
			input: []byte("a"),
			seed:  0,
			want:  "a96faf705af16834e6c632b61e964e1f",
		},
		{
			input: []byte("abc"),
			seed:  0,
			want:  "06b05ab6733a618578af5f94892f3950",
		},
		{
			input: []byte("harsh"),
			seed:  0,
			want:  "e55203a45bd698a0e91b4075969396a9",
		},
		{
			input: []byte("harsh"),
			seed:  12345,
			want:  "c99e74189a329b6a5adde6da35134269",
		},
		{
			input: []byte("The quick brown fox jumps over the lazy dog"),
			seed:  0,
			want:  "ddd650205ca3e7fa24a1cc2e3a8a7651",
		},
	}

	for _, tc := range tests {
		h := Sum128WithSeed(tc.input, tc.seed)
		got := fmt.Sprintf("%016x%016x", h.Hi, h.Lo)
		if got != tc.want {
			t.Errorf("Sum128StringWithSeed(%q, %d) = %v; want: %v",
				tc.input, tc.seed, got, tc.want)
		}
	}
}

func TestNew128struct(t *testing.T) {
	h := New128()
	h.Write([]byte("harsh"))
	hash := h.Sum128()
	got := fmt.Sprintf("%016x%016x", hash.Hi, hash.Lo)
	t.Log(got)
	want := "e55203a45bd698a0e91b4075969396a9"
	if got != want {
		t.Errorf("%v; want: %v", got, want)
	}
}
