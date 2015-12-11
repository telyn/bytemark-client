package util

import (
	bigv "bytemark.co.uk/client/lib"
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
		test{"archive:200", []bigv.Disc{bigv.Disc{StorageGrade: "archive", Size: 200 * 1024}}},
		test{"10M,11MiB,12MB,13G,14GB,15GiB", []bigv.Disc{bigv.Disc{Size: 10}, bigv.Disc{Size: 11}, bigv.Disc{Size: 12}, bigv.Disc{Size: 13 * 1024}, bigv.Disc{Size: 14 * 1024}, bigv.Disc{Size: 15 * 1024}}},
		test{"archive:200M", []bigv.Disc{bigv.Disc{StorageGrade: "archive", Size: 200}}},
		test{"archive:200Mib", []bigv.Disc{bigv.Disc{StorageGrade: "archive", Size: 200}}},
		test{"archive:200mib", []bigv.Disc{bigv.Disc{StorageGrade: "archive", Size: 200}}},
		test{"archive:200mb", []bigv.Disc{bigv.Disc{StorageGrade: "archive", Size: 200}}},
		test{"archive:200MB", []bigv.Disc{bigv.Disc{StorageGrade: "archive", Size: 200}}},
	}

	for n, x := range tests {
		r, err := ParseDiscSpec(x.Spec, true)
		if err != nil {
			t.Errorf("Test %d: %s", n, err.Error())
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
