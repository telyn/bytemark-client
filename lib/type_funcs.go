package lib

import (
	"fmt"
	"strings"
)

// StringWithSep combines all the IPs into a single string with the given seperator
func (ips IPs) StringSep(sep string) string {
	return strings.Join(ips.Strings(), sep)
}

func (ips IPs) Strings() (strings []string) {
	strings = make([]string, len(ips))
	for i, ip := range ips {
		strings[i] = ip.String()
	}
	return
}

func (ips IPs) String() string {
	return ips.StringSep(", ")
}

func (vm VirtualMachineName) String() string {
	if vm.Group == "" {
		vm.Group = "default"
	}
	if vm.Account == "" {
		return fmt.Sprintf("%s.%s", vm.VirtualMachine, vm.Group)
	}
	return fmt.Sprintf("%s.%s.%s", vm.VirtualMachine, vm.Group, vm.Account)
}

func (vm VirtualMachineName) GroupName() *GroupName {
	return &GroupName{
		Group:   vm.Group,
		Account: vm.Account,
	}
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

func (c *bytemarkClient) validateVirtualMachineName(vm *VirtualMachineName) error {
	if vm.Account == "" {
		vm.Account = c.authSession.Username
	}
	if vm.Group == "" {
		vm.Group = "default"
	}
	if vm.VirtualMachine == "" {
		return BadNameError{Type: "virtual machine", ProblemField: "name", ProblemValue: vm.VirtualMachine}
	}
	return nil
}

func (c *bytemarkClient) validateGroupName(group *GroupName) error {
	if group.Account == "" {
		group.Account = c.authSession.Username
	}
	if group.Group == "" {
		group.Group = "default"
	}
	return nil
}

func (c *bytemarkClient) validateAccountName(account *string) error {
	if *account == "" {
		*account = c.authSession.Username
	}
	return nil
}

// ParseVirtualMachineName parses a VM name given in vm[.group[.account[.extrabits]]] format
func (c *bytemarkClient) ParseVirtualMachineName(name string, defaults ...*VirtualMachineName) (vm *VirtualMachineName, err error) {
	vm = new(VirtualMachineName)

	if len(defaults) == 0 {
		vm.Group = ""
		vm.Account = ""
		vm.VirtualMachine = ""
	} else {
		vm.Group = defaults[0].Group
		vm.Account = defaults[0].Account
		vm.VirtualMachine = defaults[0].VirtualMachine
	}

	bits := strings.SplitN(name, ".", 2)
	vm.VirtualMachine = bits[0]
	if len(bits) > 1 {
		gp := c.ParseGroupName(bits[1], &GroupName{Group: vm.Group, Account: vm.Account})
		vm.Group = gp.Group
		vm.Account = gp.Account
	}

	if vm.VirtualMachine == "" {
		return vm, BadNameError{Type: "virtual machine", ProblemField: "name", ProblemValue: vm.VirtualMachine}
	}
	return vm, nil
}

// ParseGroupName parses a group name given in group[.account[.extrabits]] format.
func (c *bytemarkClient) ParseGroupName(name string, defaults ...*GroupName) (group *GroupName) {
	group = new(GroupName)
	// 1 or 2 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	if len(defaults) == 0 {
		group.Group = ""
		group.Account = ""
	} else {
		group.Group = defaults[0].Group
		group.Account = defaults[0].Account
	}

	bits = strings.SplitN(name, ".", 2)
	if len(bits) >= 1 {
		group.Group = bits[0]

	}
	if len(bits) >= 2 {
		group.Account = c.ParseAccountName(bits[1], group.Account)
	}
	return group

}

// ParseAccountName parses a group name given in .account[.extrabits] format.
// If there is a blank account name, tries to figure out the best possible account name to use.
// If authentication has already happened, this also involves asking bmbilling.
func (c *bytemarkClient) ParseAccountName(name string, defaults ...string) (account string) {
	// 1 piece with optional extra cruft for the fqdn

	if len(defaults) != 0 {
		account = defaults[0]
	}
	// if the default was blank, then use billing (if auth has already happened)
	if account == "" && c.authSession != nil {
		billAcc, err := c.getDefaultBillingAccount()
		if err == nil {
			account = billAcc.Name
		}
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

func (a *Account) FillBrain(b *brainAccount) {
	if b != nil {
		a.BrainID = b.ID
		a.Groups = b.Groups
		a.Suspended = b.Suspended
		a.Name = b.Name
	}
}
func (a *Account) FillBilling(b *billingAccount) {
	if b != nil {
		a.BillingID = b.ID
		a.Owner = b.Owner
		a.TechnicalContact = b.TechnicalContact
		a.CardReference = b.CardReference
		a.Name = b.Name
	}
}

func (a *Account) CountVirtualMachines() (servers int) {
	for _, g := range a.Groups {
		servers += len(g.VirtualMachines)
	}
	return
}

/*
func (a *Account) ToBillingAccount() *billingAccount {

}
func (a *Account) ToBrainAccount() *brainAccount {

}
*/
