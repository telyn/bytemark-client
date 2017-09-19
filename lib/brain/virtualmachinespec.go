package brain

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
)

// VirtualMachineSpec represents the specification for a VM that is passed to the create_vm endpoint
type VirtualMachineSpec struct {
	VirtualMachine VirtualMachine `json:"virtual_machine"`
	Discs          []Disc         `json:"discs,omitempty"`
	Reimage        *ImageInstall  `json:"reimage,omitempty"` // may want to be null, so is a pointer
	IPs            *IPSpec        `json:"ips,omitempty"`     // may want to be null, so is a pointer
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (spec VirtualMachineSpec) DefaultFields(f output.Format) string {
	// TODO: work on this?
	return "String"
}

func (spec VirtualMachineSpec) String() string {
	if spec.Reimage == nil {
		return "No image specified"
	}
	return fmt.Sprintf("image: %s\nroot password: %s", spec.Reimage.Distribution, spec.Reimage.RootPassword)
}

// PrettyPrint outputs a human-readable spec to the given writer.
// TODO(telyn): rewrite using templates
func (spec VirtualMachineSpec) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	output := make([]string, 0, 10)
	output = append(output, fmt.Sprintf("Name: '%s'", spec.VirtualMachine.Name))
	s := ""
	if spec.VirtualMachine.Cores > 1 {
		s = "s"
	}

	mems := fmt.Sprintf("%d", spec.VirtualMachine.Memory/1024)
	if 0 != math.Mod(float64(spec.VirtualMachine.Memory), 1024) {
		mem := float64(spec.VirtualMachine.Memory) / 1024.0
		mems = fmt.Sprintf("%.2f", mem)
	}
	output = append(output, fmt.Sprintf("Specs: %d core%s and %sGiB memory", spec.VirtualMachine.Cores, s, mems))

	locked := ""
	if spec.VirtualMachine.HardwareProfile != "" {
		if spec.VirtualMachine.HardwareProfileLocked {
			locked = " (locked)"
		}
		output = append(output, fmt.Sprintf("Hardware profile: %s%s", spec.VirtualMachine.HardwareProfile, locked))
	}

	if spec.IPs != nil {
		if spec.IPs.IPv4 != "" {
			output = append(output, fmt.Sprintf("IPv4 address: %s", spec.IPs.IPv4))
		}
		if spec.IPs.IPv6 != "" {
			output = append(output, fmt.Sprintf("IPv6 address: %s", spec.IPs.IPv6))
		}
	}

	if spec.Reimage != nil {
		if spec.Reimage.Distribution == "" {
			if spec.VirtualMachine.CdromURL == "" {
				output = append(output, "No image or CD URL specified")
			} else {
				output = append(output, fmt.Sprintf("CD URL: %s", spec.VirtualMachine.CdromURL))
			}
		} else {
			output = append(output, "Image: "+spec.Reimage.Distribution)
		}
		output = append(output, "Root/Administrator password: "+spec.Reimage.RootPassword)
	} else {

		if spec.VirtualMachine.CdromURL == "" {
			output = append(output, "No image or CD URL specified")
		} else {
			output = append(output, fmt.Sprintf("CD URL: %s", spec.VirtualMachine.CdromURL))
		}
	}

	s = ""
	if len(spec.Discs) > 1 {
		s = "s"
	}
	if len(spec.Discs) > 0 {
		output = append(output, fmt.Sprintf("%d disc%s: ", len(spec.Discs), s))
		for i, disc := range spec.Discs {
			desc := fmt.Sprintf("Disc %d", i)
			if i == 0 {
				desc = "Boot disc"
			}

			output = append(output, fmt.Sprintf("    %s %d GiB, %s grade", desc, disc.Size/1024, disc.StorageGrade))
			if len(disc.BackupSchedules) > 0 {
				bs := disc.BackupSchedules[0]
				output = append(output, fmt.Sprintf("      with automated backups %s starting at %s", bs.IntervalInWords(), bs.StartDate))
				output = append(output, fmt.Sprintf("      keeping up to %d backups (%d GiB backup usage at maximum)", bs.Capacity, bs.Capacity*disc.Size/1024))
			} else {
				output = append(output, fmt.Sprintf("      without automated backups"))
			}
		}
	} else {
		output = append(output, "No discs specified")
	}
	_, err := wr.Write([]byte(strings.Join(output, "\r\n") + "\r\n\r\n"))
	return err
}
