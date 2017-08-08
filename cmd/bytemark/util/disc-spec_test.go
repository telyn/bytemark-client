package util

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func TestParseDiscSpec(t *testing.T) {
	type test struct {
		Spec string
		Disc brain.Disc
	}
	tests := []test{
		test{"25", brain.Disc{StorageGrade: "", Size: 25 * 1024}},
		test{"archive:200", brain.Disc{StorageGrade: "archive", Size: 200 * 1024}},
		test{"archive:200M", brain.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200Mib", brain.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200mib", brain.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200mb", brain.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200MB", brain.Disc{StorageGrade: "archive", Size: 200}},
	}

	for n, x := range tests {
		d := x.Disc
		r, err := ParseDiscSpec(x.Spec)
		if err != nil {
			t.Errorf("Test %d: %s", n, err.Error())
			continue
		}
		if d.StorageGrade != r.StorageGrade {
			t.Errorf("TestParseDiscSpec %d: Expected storage grade '%s', got '%s'", n, d.StorageGrade, r.StorageGrade)
		}
		if d.Size != r.Size {
			t.Errorf("TestParseDiscSpec %d: Expected size '%d', got '%d'", n, d.Size, r.Size)
		}
	}

}
