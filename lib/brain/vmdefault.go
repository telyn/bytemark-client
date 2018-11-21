package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

type VirtualMachineDefault struct {
	Name           string             `json:"name"`
	Public         bool               `json:"public"`
	ServerSettings VirtualMachineSpec `json:"server_settings"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (vmd VirtualMachineDefault) DefaultFields(f output.Format) string {
	return "Name, Public, ServerSettings"
}

// PrettyPrint outputs a nice human-readable overview of the VM Default to the given writer.
func (vmd VirtualMachineDefault) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{ define "vmdspec_sgl" }} ▸ {{.Name }} with public => {{.Public }}{{ end }}`
	return prettyprint.Run(wr, template, "vmdspec"+string(detail), vmd)
}

// TotalDiscSize returns the sum of all disc capacities in the VM for the given storage grade.
// Provide the empty string to sum all discs regardless of storage grade.
func (vmd VirtualMachineDefault) TotalDiscSize(storageGrade string) (total int) {
	total = 0
	for _, disc := range vmd.ServerSettings.Discs {
		if storageGrade == "" || storageGrade == disc.StorageGrade {
			total += disc.Size
		}
	}
	return total
}
