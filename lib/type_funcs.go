package lib

import (
	"fmt"
	"net"
	"strings"
)

func (vm VirtualMachineName) String() string {
	if vm.Group == "" {
		vm.Group = "default"
	}
	if vm.Account == "" {
		return fmt.Sprintf("%s.%s", vm.VirtualMachine, vm.Group)
	}
	return fmt.Sprintf("%s.%s.%s", vm.VirtualMachine, vm.Group, vm.Account)
}

func (g GroupName) String() string {
	if g.Group == "" {
		g.Group = "default"
	}
	if g.Account == "" {
		return g.Group
	}
	return g.Group + "." + g.Account
}

func (bigv *bigvClient) validateVirtualMachineName(vm *VirtualMachineName) error {
	if vm.Account == "" {
		vm.Account = bigv.authSession.Username
	}
	if vm.Group == "" {
		vm.Group = "default"
	}
	if vm.VirtualMachine == "" {
		return BadNameError{Type: "virtual machine", ProblemField: "name", ProblemValue: vm.VirtualMachine}
	}
	return nil
}

func (bigv *bigvClient) validateGroupName(group *GroupName) error {
	if group.Account == "" {
		group.Account = bigv.authSession.Username
	}
	if group.Group == "" {
		group.Group = "default"
	}
	return nil
}

func (bigv *bigvClient) validateAccountName(account *string) error {
	if *account == "" {
		*account = bigv.authSession.Username
	}
	return nil
}

// ParseVirtualMachineName parses a VM name given in vm[.group[.account[.extrabits]]] format
func (bigv *bigvClient) ParseVirtualMachineName(name string, defaults ...VirtualMachineName) (vm VirtualMachineName, err error) {
	// 1, 2 or 3 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	if len(defaults) == 0 {
		vm.Group = ""
		vm.Account = ""
		vm.VirtualMachine = ""
	} else {
		vm.Group = defaults[0].Group
		vm.Account = defaults[0].Account
		vm.VirtualMachine = defaults[0].VirtualMachine
	}

	if len(bits) > 3 && bits[len(bits)-1] == "" {
		bits = bits[0 : len(bits)-1]
	}

	// a for loop seems an odd choice here maybe but it means
	// I don't need to do lots of ifs to see if the next bit exists
Loop:
	for i, bit := range bits {
		bit = strings.TrimSpace(strings.ToLower(bit))
		if bit != "" {
			switch i {
			case 0:
				vm.VirtualMachine = bit
				break
			case 1:
				// want to be able to do vm..account to get the default group
				vm.Group = bit
				break
			case 2:
				vm.Account = bit
				break Loop
			}
		}
	}
	if vm.VirtualMachine == "" {
		return vm, BadNameError{Type: "virtual machine", ProblemField: "name", ProblemValue: vm.VirtualMachine}
	}
	return vm, nil
}

// ParseGroupName parses a group name given in group[.account[.extrabits]] format.
func (bigv *bigvClient) ParseGroupName(name string, defaults ...GroupName) (group GroupName) {
	// 1 or 2 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	if len(defaults) == 0 {
		group.Group = ""
		group.Account = ""
	} else {
		group.Group = defaults[0].Group
		group.Account = defaults[0].Account
	}

Loop:
	for i, bit := range bits {
		bit = strings.TrimSpace(strings.ToLower(bit))
		if bit != "" {
			switch i {
			case 0:
				group.Group = bit
				break
			case 1:
				group.Account = bit
				break Loop
			}
		}
	}
	return group

}

// ParseAccountName parses a group name given in .account[.extrabits] format.
func (bigv *bigvClient) ParseAccountName(name string, defaults ...string) (account string) {
	// 1 piece with optional extra cruft for the fqdn

	if len(defaults) == 0 {
		account = ""
	} else {
		account = defaults[0]
	}

	// there's a micro-optimisation to do here to not use Split,
	// but really, who can be bothered to?
	bits := strings.Split(name, ".")
	if bits[0] != "" {
		account = bits[0]
	}

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
