package admin

import "testing"

func TestStringsToNumberOrStrings(t *testing.T) {
	in := []string{"hello", "world", "8092"}
	out := stringsToNumberOrStrings(in)
	if len(out) != 3 {
		t.Fatalf("len(out) = %d, expected 3", len(out))
	}
	for i := range in {
		if string(out[i]) != in[i] {
			t.Fatalf("out[%d] = %q, expected %q", i, out[i], in[i])
		}
	}
}
