package util

import (
	bigv "bytemark.co.uk/client/lib"
	"testing"
)

func TestParseDiscSpec(t *testing.T) {
	type test struct {
		Spec string
		Disc bigv.Disc
	}
	tests := []test{
		test{"25", bigv.Disc{StorageGrade: "", Size: 25 * 1024}},
		test{"archive:200", bigv.Disc{StorageGrade: "archive", Size: 200 * 1024}},
		test{"archive:200M", bigv.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200Mib", bigv.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200mib", bigv.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200mb", bigv.Disc{StorageGrade: "archive", Size: 200}},
		test{"archive:200MB", bigv.Disc{StorageGrade: "archive", Size: 200}},
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
