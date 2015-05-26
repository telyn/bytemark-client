package main

import (
	client "bigv.io/client/lib"
	"fmt"
	"strings"
)

// VMFormatOptions controls formatting of VMs in FormatVirtualMachine
// Add or or them together to get what you want
type VMFormatOptions uint8

const (
	// FORMAT_VM_WITH_ADDRS causes IP addresses to be included in the output
	FORMAT_VM_WITH_ADDRS VMFormatOptions = 1 << iota
	// FORMAT_VM_WITH_DISCS causes individual disc sizes & storage grades to be included in the output
	FORMAT_VM_WITH_DISCS
	// FORMAT_VM_WITH_CDURL causes the URL of the image being used as the CD to be included in the output, if applicable
	FORMAT_VM_WITH_CDURL
)

// VMListFormatMode is the way that FormatVirtualMachineList will format the VMList
type VMListFormatMode uint8

const (
	// FORMAT_VMLIST_NAMES outputs only the names of the VMs
	FORMAT_VMLIST_NAMES VMListFormatMode = iota
	// FORMAT_VMLIST_WITH_GROUPS outputs the VMs in name.group format
	FORMAT_VMLIST_WITH_GROUPS
	// FORMAT_VMLIST_FULL outputs the full hostnames of the VMs.
	FORMAT_VMLIST_FULL
)

// FORMAT_DEFAULT_WIDTH is the default width to attempt to print to.
const FORMAT_DEFAULT_WIDTH = 80

// FormatVirtualMachine pretty-prints a VM. The optional second argument is a bitmask of VMFormatOptions,
// and the optional third is the width you'd like to display..oh.
func FormatVirtualMachine(vm *client.VirtualMachine, options ...int) string {
	width := FORMAT_DEFAULT_WIDTH
	format := FORMAT_VM_WITH_ADDRS | FORMAT_VM_WITH_DISCS

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

	output = append(output, "")
	if (format & FORMAT_VM_WITH_DISCS) != 0 {
		for _, disc := range vm.Discs {
			output = append(output, fmt.Sprintf("Disc %s: %d GiB, %s grade", disc.Label, disc.Size/1024, disc.StorageGrade))
		}
		output = append(output, "")
	}

	if (format & FORMAT_VM_WITH_ADDRS) != 0 {
		output = append(output, fmt.Sprintf("IPv4 Addresses: %s\r\n", strings.Join(vm.AllIpv4Addresses(), ",\r\n                ")))
		output = append(output, fmt.Sprintf("IPv6 Addresses: %s\r\n", strings.Join(vm.AllIpv6Addresses(), ",\r\n                ")))
	}

	return strings.Join(output, "\r\n")

}
