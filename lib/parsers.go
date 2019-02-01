package lib

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// ParseAccountName parses a group name given in account[.extrabits] format.
// If there is a blank account name, tries to figure out the best possible account name to use.
// If authentication has already happened, this also involves asking bmbilling.
func ParseAccountName(name string, defaults ...string) (account string) {
	// 1 piece with optional extra cruft for the fqdn

	if len(defaults) != 0 {
		account = defaults[0]
	}

	bits := strings.Split(name, ".")
	if bits[0] != "" {
		account = bits[0]
	}

	return account
}

// ParseGroupName parses a group name given in group[.account[.extrabits]] format.
func ParseGroupName(name string, defaults ...pathers.GroupName) (group pathers.GroupName) {
	if len(defaults) == 0 {
		group.Group = ""
		group.Account = ""
	} else {
		group.Group = defaults[0].Group
		group.Account = defaults[0].Account
	}

	bits := strings.SplitN(name, ".", 2)
	if len(bits) >= 1 {
		group.Group = bits[0]

	}
	if len(bits) >= 2 {
		// gross. sorry.
		group.Account = pathers.AccountName(ParseAccountName(bits[1], string(group.Account)))
	}
	return group

}

// ParseVirtualMachineName parses a VM name given in vm[.group[.account[.extrabits]]] format
func ParseVirtualMachineName(name string, defaults ...VirtualMachineName) (vm VirtualMachineName, err error) {
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
		gp := ParseGroupName(bits[1], pathers.GroupName{Group: vm.Group, Account: vm.Account})
		vm.Group = gp.Group
		vm.Account = gp.Account
	}

	if vm.VirtualMachine == "" {
		return vm, BadNameError{Type: "virtual machine", ProblemField: "name", ProblemValue: vm.VirtualMachine}
	}
	return vm, nil
}
