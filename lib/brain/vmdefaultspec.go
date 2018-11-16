package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/output"
)

// TODO (tom): Add pretty print & test file for this function
// VMDefaultSpec represents a VM Default specification.
type VMDefaultSpec struct {
	VMDefault VMDefault     `json:"vm_default,omitempty"`
	Discs     []Disc        `json:"disc,omitempty"`
	Reimage   *ImageInstall `json:"reimage,omitempty"` // may want to be null, so is a pointer
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (spec VMDefaultSpec) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "VirtualMachine, Discs, Reimage"
	}
	return "VirtualMachine, Discs, Reimage"
}