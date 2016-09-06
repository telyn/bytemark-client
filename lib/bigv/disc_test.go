package bigv

import (
	"github.com/cheekybits/is"
	"testing"
)

func getFixtureDisc() Disc {
	return Disc{
		Label:            "",
		StorageGrade:     "sata",
		Size:             26400,
		ID:               1,
		VirtualMachineID: 1,
		StoragePool:      "fakepool",
	}
}

func getFixtureDiscSet() []Disc {
	return []Disc{
		getFixtureDisc(),
		Disc{
			ID:           2,
			StorageGrade: "archive",
			Label:        "arch",
			Size:         1024000,
		},
		Disc{
			ID:           3,
			StorageGrade: "",
			Size:         2048,
		},
	}
}

func TestValidateDisc(t *testing.T) {
	is := is.New(t)
	discs := getFixtureDiscSet()
	for _, d := range discs {
		d2, err := d.Validate()
		is.Nil(err)

		is.Equal(d.ID, d2.ID)
		is.Equal(d.Label, d2.Label)
		is.Equal(d.Size, d2.Size)
		is.Equal(d.StoragePool, d2.StoragePool)
		is.Equal(d.VirtualMachineID, d2.VirtualMachineID)
		switch d.ID {
		case 1, 3:
			is.Equal("sata", d2.StorageGrade)
		case 2:
			is.Equal("archive", d2.StorageGrade)
		}
	}
}
