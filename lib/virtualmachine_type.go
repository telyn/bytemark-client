package lib

import (
	"net"
)

// VirtualMachine represents a VirtualMachine, as passed around from the virtual_machines endpoint
type VirtualMachine struct {
	Autoreboot            bool   `json:"autoreboot_on"`
	CdromURL              string `json:"cdrom_url"`
	Cores                 int    `json:"cores"`
	Memory                int    `json:"memory"`
	Name                  string `json:"name"`
	PowerOn               bool   `json:"power_on"`
	HardwareProfile       string `json:"hardware_profile"`
	HardwareProfileLocked bool   `json:"hardware_profile_locked"`
	GroupID               int    `json:"group_id"`

	// zone name can be set during creation but not changed
	ZoneName string `json:"zone_name"`

	// the following cannot be set
	Discs             []*Disc             `json:"discs"`
	ID                int                 `json:"id"`
	ManagementAddress *net.IP             `json:"management_address"`
	Deleted           bool                `json:"deleted"`
	Hostname          string              `json:"hostname"`
	Head              string              `json:"head"`
	NetworkInterfaces []*NetworkInterface `json:"network_interfaces"`

	// TODO(telyn): new fields (last_imaged_with and there is another but I forgot)
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

// AllIPv4Addresses flattens all the IPs for a VM into a single IPs (a []*net.IP with some convenience methods)
func (vm *VirtualMachine) AllIPv4Addresses() (ips IPs) {
	for _, nic := range vm.NetworkInterfaces {
		for _, ip := range nic.IPs {
			if ip != nil && ip.To4() != nil {
				ips = append(ips, ip)
			}
		}
		for _, ip := range nic.ExtraIPs {
			if ip != nil && ip.To4() != nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips
}

// AllIPv6Addresses flattens all the v6 IPs for a VM into a single IPs (a []*net.IP with some convenience methods)
func (vm *VirtualMachine) AllIPv6Addresses() (ips IPs) {
	for _, nic := range vm.NetworkInterfaces {
		for _, ip := range nic.IPs {
			if ip != nil && ip.To4() == nil {
				ips = append(ips, ip)
			}
		}
		for _, ip := range nic.ExtraIPs {
			if ip != nil && ip.To4() == nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips
}
