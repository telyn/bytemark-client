package cmd

import (
	client "bigv.io/client/lib"
	"fmt"
	"strings"
)

// FORMAT_VM_WITH_* controls formatting of VMs.
// Add or or them together to get what you want
const (
	// include IP addresses in the output
	FORMAT_VM_WITH_ADDRS = 1 << iota
	// include individual disc sizes & storage grades
	FORMAT_VM_WITH_DISCS
	// Assuming the VM has a CD inserted, include the URL of the image being used as the CD
	FORMAT_VM_WITH_CDURL
)

const (
	FORMAT_VMLIST_NAMES = iota
	FORMAT_VMLIST_WITH_GROUPS
	FORMAT_VMLIST_FULL
)

const FORMAT_DEFAULT_WIDTH = 80

// FormatVirtualMachine pretty-prints a VM. The optional second argument is a bitmask of FORMAT_VM_WITH_* integers,
// and the optional third is the width you'd like to display.
func FormatVirtualMachine(vm *client.VirtualMachine, options ...int) string {
	width := FORMAT_DEFAULT_WIDTH
	format := FORMAT_VM_WITH_ADDRS | FORMAT_VM_WITH_DISCS

	if len(options) >= 1 {
		format = options[0]
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
