package util

import (
	bigv "bigv.io/client/lib"
	"testing"
)

func TestParseDiscSpec(t *testing.T) {
	type test struct {
		Spec  string
		Discs []bigv.Disc
	}
	tests := []test{
		test{"25", []bigv.Disc{bigv.Disc{StorageGrade: "", Size: 25 * 1024}}},
		test{"10,32", []bigv.Disc{bigv.Disc{StorageGrade: "", Size: 10 * 1024}, bigv.Disc{StorageGrade: "", Size: 32 * 1024}}},
		test{"15,archive:50", []bigv.Disc{bigv.Disc{StorageGrade: "", Size: 15 * 1024}, bigv.Disc{StorageGrade: "archive", Size: 50 * 1024}}},
	}

	for n, x := range tests {
		r, err := ParseDiscSpec(x.Spec, true)
		if err != nil {
			t.Error(err)
			continue
		}
		if len(x.Discs) != len(r) {
			t.Errorf("Test %d: Expected %d discs, got %d", n, len(x.Discs), len(r))
		}
		for i, d := range x.Discs {
			if d.StorageGrade != r[i].StorageGrade {
				t.Errorf("TestParseDiscSpec %d: Disc %d: Expected storage grade '%s', got '%s'", n, i, d.StorageGrade, r[i].StorageGrade)
			}
			if d.Size != r[i].Size {
				t.Errorf("TestParseDiscSpec %d: Disc %d: Expected size '%d', got '%d'", n, i, d.Size, r[i].Size)
			}
		}
	}

}
