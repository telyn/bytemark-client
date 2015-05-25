package lib

import (
	"strings"
)

// by convention this function uses DEFAULT in all-caps to mean "the default group/account", as set in the config, rather than the "default" group in BigV itself.

// ParseVirtualMachineName parses a VM name given in vm[.group[.account[.extrabits]]] format
func (bigv *BigVClient) ParseVirtualMachineName(name string) (vm VirtualMachineName) {
	// 1, 2 or 3 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	vm.Group = ""
	vm.Account = ""
	vm.VirtualMachine = ""

	if len(bits) > 3 && bits[len(bits)-1] == "" {
		bits = bits[0 : len(bits)-1]
	}

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

// by convention this function uses DEFAULT in all-caps to mean "the default group/account", as set in the config, rather than the "default" group in BigV itself.

// ParseGroupName parses a group name given in group[.account[.extrabits]] format.
func (bigv *BigVClient) ParseGroupName(name string) (group GroupName) {
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

// by convention this function uses DEFAULT in all-caps to mean "the default account", as set in the config, rather than the "default" group in BigV itself.

// ParseAccountName parses a group name given in .account[.extrabits] format.
func (bigv *BigVClient) ParseAccountName(name string) (account string) {
	// 1 piece with optional extra cruft for the fqdn

	// there's a micro-optimisation to do here to not use Split,
	// but really, who can be bothered to?
	account = ""

	bits := strings.Split(name, ".")
	account = bits[0]

	return account

}
