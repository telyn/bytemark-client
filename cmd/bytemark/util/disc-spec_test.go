package util

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"testing"
)

func TestParseDiscSpec(t *testing.T) {
	type test struct {
		Spec string
		Disc lib.Disc
	}
	tests := []test{
		test{"25", lib.Disc{StorageGrade: "", Size: 25 * 1024}},
		test{"archive:200", lib.Disc{StorageGrade: "archive", Size: 200 * 1024}},
		test{"archive:200M", lib.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200Mib", lib.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200mib", lib.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200mb", lib.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200MB", lib.Disc{StorageGrade: "archive", Size: 200}},
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
