package brain

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
	"math"
	"strings"
)

// VirtualMachineSpec represents the specification for a VM that is passed to the create_vm endpoint
type VirtualMachineSpec struct {
	VirtualMachine *VirtualMachine `json:"virtual_machine"`
	Discs          []Disc          `json:"discs,omitempty"`
	Reimage        *ImageInstall   `json:"reimage,omitempty"`
	IPs            *IPSpec         `json:"ips,omitempty"`
}

// PrettyPrint outputs a human-readable spec to the given writer.
// TODO(telyn): rewrite using templates
func (spec VirtualMachineSpec) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	output := make([]string, 0, 10)
	output = append(output, fmt.Sprintf("Name: '%s'", pp.VirtualMachine.Name))
	s := ""
	if pp.VirtualMachine.Cores > 1 {
		s = "s"
	}

	mems := fmt.Sprintf("%d", pp.VirtualMachine.Memory/1024)
	if 0 != math.Mod(float64(pp.VirtualMachine.Memory), 1024) {
		mem := float64(pp.VirtualMachine.Memory) / 1024.0
		mems = fmt.Sprintf("%.2f", mem)
	}
	output = append(output, fmt.Sprintf("Specs: %d core%s and %sGiB memory", pp.VirtualMachine.Cores, s, mems))

	locked := ""
	if pp.VirtualMachine.HardwareProfile != "" {
		if pp.VirtualMachine.HardwareProfileLocked {
			locked = " (locked)"
		}
		output = append(output, fmt.Sprintf("Hardware profile: %s%s", pp.VirtualMachine.HardwareProfile, locked))
	}

	if pp.IPs != nil {
		if pp.IPs.IPv4 != "" {
			output = append(output, fmt.Sprintf("IPv4 address: %s", pp.IPs.IPv4))
		}
		if pp.IPs.IPv6 != "" {
			output = append(output, fmt.Sprintf("IPv6 address: %s", pp.IPs.IPv6))
		}
	}

	if pp.Reimage != nil {
		if pp.Reimage.Distribution == "" {
			if pp.VirtualMachine.CdromURL == "" {
				output = append(output, "No image or CD URL ppified")
			} else {
				output = append(output, fmt.Sprintf("CD URL: %s", pp.VirtualMachine.CdromURL))
			}
		} else {
			output = append(output, "Image: "+pp.Reimage.Distribution)
		}
		output = append(output, "Root/Administrator password: "+pp.Reimage.RootPassword)
	} else {

		if pp.VirtualMachine.CdromURL == "" {
			output = append(output, "No image or CD URL ppified")
		} else {
			output = append(output, fmt.Sprintf("CD URL: %s", pp.VirtualMachine.CdromURL))
		}
	}

	s = ""
	if len(pp.Discs) > 1 {
		s = "s"
	}
	if len(pp.Discs) > 0 {
		output = append(output, fmt.Sprintf("%d disc%s: ", len(pp.Discs), s))
		for i, disc := range pp.Discs {
			desc := fmt.Sprintf("Disc %d", i)
			if i == 0 {
				desc = "Boot disc"
			}

			output = append(output, fmt.Sprintf("    %s %d GiB, %s grade", desc, disc.Size/1024, disc.StorageGrade))
		}
	} else {
		output = append(output, "No discs ppified")
	}
	_, err := wr.Write([]byte(strings.Join(output, "\r\n") + "\r\n"))
	return err
}
