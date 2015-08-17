package main

import (
	"testing"
)

func TestParseSize(t *testing.T) {
	type test struct {
		Spec string
		Size int
	}
	tests := []test{
		test{"1", 1024},
		test{"2", 2048},
		test{"25", 25600},
		test{"2 GiB", 2048},
		test{"2 GB", 2048},
		test{"2 gb", 2048},
		test{"2GiB", 2048},
		test{"200MB", 200},
		test{"200 MB", 200},
		test{"200 mb", 200},
	}
	for n, x := range tests {
		r, err := ParseSize(x.Spec)
		if err != nil {
			t.Error(err)
			continue
		}
		if r != x.Size {
			t.Errorf("Test %d ('%s'): Expected %d, got %d", n, x.Spec, x.Size, r)
		}
	}
}
