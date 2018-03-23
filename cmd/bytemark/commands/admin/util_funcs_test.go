package admin

import "testing"

func TestStringsToJsonNumbers(t *testing.T) {
	in := []string{"hello", "world", "8092"}
	out := stringsToJsonNumbers(in)
	if len(out) != 3 {
		t.Fatalf("len(out) = %d, expected 3", len(out))
	}
	for i := range in {
		if string(out[i]) != in[i] {
			t.Fatalf("out[%d] = %q, expected %q", i, out[i], in[i])
		}
	}
}
