package lib

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func TestLabelDisc(t *testing.T) {
	is := is.New(t)
	discs := []brain.Disc{
		{
			Label:            "",
			StorageGrade:     "sata",
			Size:             26400,
			ID:               1,
			VirtualMachineID: 1,
			StoragePool:      "fakepool",
		},
		{
			ID:           2,
			StorageGrade: "archive",
			Label:        "arch",
			Size:         1024000,
		},
		{
			ID:           3,
			StorageGrade: "",
			Size:         2048,
		},
	}
	labelDiscs(discs)
	for _, d := range discs {
		switch d.ID {
		case 1:
			is.Equal("disc-1", d.Label)
		case 2:
			is.Equal("arch", d.Label)
		case 3:
			is.Equal("disc-3", d.Label)
		default:
			fmt.Printf("Unexpected disc ID %d\r\n", d.ID)
			t.Fail()
		}
	}
}
