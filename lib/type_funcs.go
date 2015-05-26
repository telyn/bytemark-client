package lib

import (
	"fmt"
	"net"
)

func (vm VirtualMachineName) String() string {
	return fmt.Sprintf("%s.%s.%s", vm.VirtualMachine, vm.Group, vm.Account)
}

func (g GroupName) String() string {
	return g.Group + "." + g.Account
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
		for _, ip := range nic.Ips {
			if net.ParseIP(ip) != nil && net.ParseIP(ip).To4() != nil {
				ips = append(ips, ip)
			}
		}
		for ip := range nic.ExtraIps {
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
		for _, ip := range nic.Ips {
			if net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil {
				ips = append(ips, ip)
			}
		}
		for ip := range nic.ExtraIps {
			if net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips
}
