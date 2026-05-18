package buffer

import (
	"testing"
)

var buffertest = []struct {
	buffersize int
	in         []string
	out        []string
}{
	{5, []string{"a"}, []string{"a"}},
	{5, []string{"a", "b", "c"}, []string{"a", "b", "c"}},
	{5, []string{"a", "b", "c", "d", "e"}, []string{"a", "b", "c", "d", "e"}},
	{5, []string{"a", "b", "c", "d", "e", "f", "g"}, []string{"c", "d", "e", "f", "g"}},
}

func TestBuffer(t *testing.T) {
	for _, tc := range buffertest {
		var b Buffer = NewRingBuffer(tc.buffersize)
		for _, text := range tc.in {
			b.Add(text)
		}

		// Verify output
		out := b.Snapshot()
		if len(out) != len(tc.out) {
			t.Errorf("length mismatch: got %d, want %d", len(out), len(tc.out))
			continue
		}

		for i := range out {
			if out[i] != tc.out[i] {
				t.Errorf("got %q, want %q", out[i], tc.out[i])
			}
		}
	}
}
