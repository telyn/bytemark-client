package lib

import (
	"fmt"
	"net"
	"strings"
)

func (vm VirtualMachineName) String() string {
	return fmt.Sprintf("%s.%s.%s", vm.VirtualMachine, vm.Group, vm.Account)
}

func (g GroupName) String() string {
	return g.Group + "." + g.Account
}

// ParseVirtualMachineName parses a VM name given in vm[.group[.account[.extrabits]]] format
func (bigv *bigvClient) ParseVirtualMachineName(name string) (vm VirtualMachineName) {
	// 1, 2 or 3 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	vm.Group = ""
	vm.Account = ""
	vm.VirtualMachine = ""

	if len(bits) > 3 && bits[len(bits)-1] == "" {
		bits = bits[0 : len(bits)-1]
	}

	// TODO(telyn): ParseVirtualMachine isn't smart enough yet

	// a for loop seems an odd choice here maybe but it means
	// I don't need to do lots of ifs to see if the next bit exists
Loop:
	for i, bit := range bits {
		switch i {
		case 0:
			vm.VirtualMachine = strings.TrimSpace(strings.ToLower(bit))
			break
		case 1:
			// want to be able to do vm..account to get the default group
			if bit == "" {
				break
			}
			vm.Group = strings.TrimSpace(strings.ToLower(bit))
			break
		case 2:
			vm.Account = strings.TrimSpace(strings.ToLower(bit))
			break Loop
		}
	}
	return vm
}

// ParseGroupName parses a group name given in group[.account[.extrabits]] format.
func (bigv *bigvClient) ParseGroupName(name string) (group GroupName) {
	// 1 or 2 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	group.Group = ""
	group.Account = ""

Loop:
	for i, bit := range bits {
		switch i {
		case 0:
			group.Group = strings.TrimSpace(strings.ToLower(bit))
			break
		case 1:
			group.Account = strings.TrimSpace(strings.ToLower(bit))
			break Loop
		}
	}
	return group

}

// ParseAccountName parses a group name given in .account[.extrabits] format.
func (bigv *bigvClient) ParseAccountName(name string) (account string) {
	// 1 piece with optional extra cruft for the fqdn

	// there's a micro-optimisation to do here to not use Split,
	// but really, who can be bothered to?
	account = ""

	bits := strings.Split(name, ".")
	account = bits[0]

	return account

}

// TotalDiscSize returns the sum of all disc capacities in the VM for the given storage grade.
// Provide the empty string to sum all discs regardless of storage grade.
func (vm *VirtualMachine) TotalDiscSize(storageGrade string) (total int) {
	total = 0
	for _, disc := range vm.Discs {
		if storageGrade == "" || storageGrade == disc.StorageGrade {
			total += disc.Size
		}
	}
	return total
}

// AllIpv4Addresses flattens all the IPs for a VM into a single []string
func (vm *VirtualMachine) AllIpv4Addresses() (ips []string) {
	for _, nic := range vm.NetworkInterfaces {
		for _, ip := range nic.IPs {
			if net.ParseIP(ip) != nil && net.ParseIP(ip).To4() != nil {
				ips = append(ips, ip)
			}
		}
		for ip := range nic.ExtraIPs {
			if net.ParseIP(ip) != nil && net.ParseIP(ip).To4() != nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips
}

// AllIpv6Addresses flattens all the v6 IPs for a VM into a single []string
func (vm *VirtualMachine) AllIpv6Addresses() (ips []string) {
	for _, nic := range vm.NetworkInterfaces {
		for _, ip := range nic.IPs {
			if net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil {
				ips = append(ips, ip)
			}
		}
		for ip := range nic.ExtraIPs {
			if net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips
}
