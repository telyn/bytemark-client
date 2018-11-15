package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/output"
)

// TODO(tom): add test

// VMDefault represents a VM Default, as passed around from the vm_defaults endpoint
type VMDefault struct {
	CdromURL              string `json:"cdrom_url,omitempty"`
	Cores                 int    `json:"cores,omitempty"`
	Memory                int    `json:"memory,omitempty"`
	Name                  string `json:"name,omitempty"`
	HardwareProfile       string `json:"hardware_profile,omitempty"`

	// zone name can be set during creation but not changed
	ZoneName string `json:"zone_name,omitempty"`

	// the following cannot be set
	Discs             Discs              `json:"discs,omitempty"`
	ID                int                `json:"id,omitempty"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (vmd VMDefault ) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "Memory, Cores, Name "
	}
	return "Memory, Cores, Discs, CdromURL, HardwareProfile"
}