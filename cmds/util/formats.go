package util

import (
	client "bytemark.co.uk/client/lib"
	"fmt"
	"math"
	"strings"
)

// VMFormatOptions controls formatting of VMs in FormatVirtualMachine
// Add or or them together to get what you want
type VMFormatOptions uint8

const (
	// _FormatVMWithAddrs causes IP addresses to be included in the output
	_FormatVMWithAddrs VMFormatOptions = 1 << iota
	// _FormatVMWithDiscs causes individual disc sizes & storage grades to be included in the output
	_FormatVMWithDiscs
	// _FormatVMWithCDURL causes the URL of the image being used as the CD to be included in the output, if applicable
	_FormatVMWithCDURL
)

// VMListFormatMode is the way that FormatVirtualMachineList will format the VMList
type VMListFormatMode uint8

const (
	// _FormatVMListName outputs only the names of the VMs
	_FormatVMListName VMListFormatMode = iota
	// _FormatVMListNameDotGroup outputs the VMs in name.group format
	_FormatVMListNameDotGroup
	// _FormatVMListFQDN outputs the full hostnames of the VMs.
	_FormatVMListFQDN
)

// FORMAT_DEFAULT_WIDTH is the default width to attempt to print to.
const _FormatDefaultWidth = 80

// FormatVirtualMachines loops through a bunch of VMs, formatting each one as it goes, and returns each formatted VM as a string.
// The options are the same as FormatVirtualMachine
func FormatVirtualMachines(vms []*client.VirtualMachine, options ...int) []string {
	output := make([]string, len(vms), len(vms))
	for i, vm := range vms {
		output[i] = FormatVirtualMachine(vm, options...)
	}
	return output
}

// FormatVirtualMachine pretty-prints a VM. The optional second argument is a bitmask of VMFormatOptions,
// and the optional third is the width you'd like to display..oh.
func FormatVirtualMachine(vm *client.VirtualMachine, options ...int) string {
	width := _FormatDefaultWidth
	format := _FormatVMWithAddrs | _FormatVMWithDiscs

	if len(options) >= 1 {
		format = VMFormatOptions(options[0])
	}

	if len(options) >= 2 {
		width = options[1]
	}

	output := make([]string, 0, 10)

	title := fmt.Sprintf(" VM %s, %d cores, %d GiB RAM, %d GiB on %d discs =", vm.Name, vm.Cores, vm.Memory/1024, vm.TotalDiscSize("")/1024, len(vm.Discs))
	padding := ""
	for i := 0; i < width-len(title); i++ {
		padding += "="
	}

	output = append(output, padding+title)

	output = append(output, fmt.Sprintf("Hostname: %s", vm.Hostname))
	if (format&_FormatVMWithCDURL) != 0 && vm.CdromURL != "" {
		output = append(output, fmt.Sprintf("CD-ROM: %s", vm.CdromURL))
	}

	output = append(output, "")
	if (format & _FormatVMWithDiscs) != 0 {
		for _, disc := range vm.Discs {
			output = append(output, fmt.Sprintf("Disc %s: %d GiB, %s grade", disc.Label, disc.Size/1024, disc.StorageGrade))
		}
		output = append(output, "")
	}

	if (format & _FormatVMWithAddrs) != 0 {
		output = append(output, fmt.Sprintf("IPv4 Addresses: %s\r\n", strings.Join(vm.AllIpv4Addresses(), ",\r\n                ")))
		output = append(output, fmt.Sprintf("IPv6 Addresses: %s\r\n", strings.Join(vm.AllIpv6Addresses(), ",\r\n                ")))
	}

	return strings.Join(output, "\r\n")
}

func FormatVirtualMachineSpec(group *client.GroupName, spec *client.VirtualMachineSpec) string {
	output := make([]string, 0, 10)
	output = append(output, fmt.Sprintf("Name: '%s'", spec.VirtualMachine.Name))
	output = append(output, fmt.Sprintf("Group: '%s'", group.Group))
	output = append(output, fmt.Sprintf("Account: '%s'", group.Account))
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

	if spec.Reimage.Distribution == "" {
		output = append(output, "No image specified")
	} else {
		output = append(output, "Image: "+spec.Reimage.Distribution)
	}

	s = ""
	if len(spec.Discs) > 1 {
		s = "s"
	}
	output = append(output, fmt.Sprintf("%d disc%s: ", len(spec.Discs), s))
	for i, disc := range spec.Discs {
		desc := fmt.Sprintf("Disc %d", i)
		if i == 0 {
			desc = "Boot disc"
		}

		output = append(output, fmt.Sprintf("    %s %d GiB, %s grade", desc, disc.Size/1024, disc.StorageGrade))
	}
	return strings.Join(output, "\r\n")

}
