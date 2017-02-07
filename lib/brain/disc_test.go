package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
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

func TestDiscPrettyPrint(t *testing.T) {
	prettyprint.RunTests(t, []prettyprint.Test{
		{
			Object: Disc{
				Label:          "important-stuff",
				Size:           25500,
				BackupCount:    4,
				BackupsEnabled: true,
				StorageGrade:   "sata",
			},
			Detail:   prettyprint.SingleLine,
			Expected: "important-stuff - 24GiB, sata grade (has 4 backups)",
		},
		{
			Object: Disc{
				Label:          "important-stuff",
				Size:           22500,
				BackupCount:    4,
				BackupsEnabled: true,
				StorageGrade:   "sata",
			},
			Detail: prettyprint.Medium,
			Expected: `important-stuff - 21GiB, sata grade (has 4 backups)
No backups scheduled`,
		},
		{
			Object: Disc{
				Label:          "important-stuff",
				Size:           25500,
				BackupCount:    4,
				BackupsEnabled: true,
				StorageGrade:   "iceberg",
			},
			Detail: prettyprint.Full,
			Expected: `important-stuff - 24GiB, iceberg grade (has 4 backups) (restore in progress)
No backups scheduled`,
		},
	})
}
