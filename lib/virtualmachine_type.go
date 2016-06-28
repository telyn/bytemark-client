package lib

import (
	"net"
)

// VirtualMachine represents a VirtualMachine, as passed around from the virtual_machines endpoint
type VirtualMachine struct {
	Autoreboot            bool   `json:"autoreboot_on,omitempty"`
	CdromURL              string `json:"cdrom_url,omitempty"`
	Cores                 int    `json:"cores,omitempty"`
	Memory                int    `json:"memory,omitempty"`
	Name                  string `json:"name,omitempty"`
	PowerOn               bool   `json:"power_on,omitempty"`
	HardwareProfile       string `json:"hardware_profile,omitempty"`
	HardwareProfileLocked bool   `json:"hardware_profile_locked,omitempty"`
	GroupID               int    `json:"group_id,omitempty"`

	// zone name can be set during creation but not changed
	ZoneName string `json:"zone_name,omitempty"`

	// the following cannot be set
	Discs             []*Disc             `json:"discs,omitempty"`
	ID                int                 `json:"id,omitempty"`
	ManagementAddress *net.IP             `json:"management_address,omitempty"`
	Deleted           bool                `json:"deleted,omitempty"`
	Hostname          string              `json:"hostname,omitempty"`
	Head              string              `json:"head,omitempty"`
	NetworkInterfaces []*NetworkInterface `json:"network_interfaces,omitempty"`

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
		for ip := range nic.ExtraIPs {
			netip := net.ParseIP(ip)
			if netip != nil && netip.To4() != nil {
				ips = append(ips, &netip)
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
		for ip := range nic.ExtraIPs {
			netip := net.ParseIP(ip)
			if netip != nil && netip.To4() == nil {
				ips = append(ips, &netip)
			}
		}
	}
	return ips
}
