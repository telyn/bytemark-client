package main

import (
	"bigv.io/client/lib"
	"strings"
)

// by convention this function uses DEFAULT in all-caps to mean "the default group/account", as set in the config, rather than the "default" group in BigV itself.

// ParseVirtualMachineName parses a VM name given in vm[.group[.account[.extrabits]]] format
// I'm not sure if this should be in bigv.io/client/lib or not yet.
func ParseVirtualMachineName(name string) (vm lib.VirtualMachineName) {
	// 1, 2 or 3 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	vm.Group = "DEFAULT"
	vm.Account = "DEFAULT"
	vm.VirtualMachine = "DEFAULT"

	// a for loop seems an odd choice here maybe but it means
	// I don't need to do lots of ifs to see if the next bit exists
Loop:
	for i, bit := range bits {
		switch i {
		case 0:
			vm.VirtualMachine = strings.ToLower(bit)
			break
		case 1:
			// want to be able to do vm..account to get the default group
			if bit == "" {
				break
			}
			vm.Group = strings.ToLower(bit)
			break
		case 2:
			vm.Account = strings.ToLower(bit)
			break Loop
		}
	}
	return vm
}

// by convention this function uses DEFAULT in all-caps to mean "the default group/account", as set in the config, rather than the "default" group in BigV itself.

// ParseGroupName parses a group name given in group[.account[.extrabits]] format.
// I'm not sure if this should be in bigv.io/client/lib or not yet.
func ParseGroupName(name string) (group lib.GroupName) {
	// 1 or 2 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	group.Group = "DEFAULT"
	group.Account = "DEFAULT"

Loop:
	for i, bit := range bits {
		switch i {
		case 0:
			group.Group = strings.ToLower(bit)
			break
		case 1:
			group.Account = strings.ToLower(bit)
			break Loop
		}
	}
	return group

}

// by convention this function uses DEFAULT in all-caps to mean "the default account", as set in the config, rather than the "default" group in BigV itself.

// ParseAccountName parses a group name given in .account[.extrabits] format.
// I'm not sure if this should be in bigv.io/client/lib or not yet.
func ParseAccountName(name string) (account string) {
	// 1 piece with optional extra cruft for the fqdn

	// there's a micro-optimisation to do here to not use Split,
	// but really, who can be bothered to?
	account = "DEFAULT"

	bits := strings.Split(name, ".")
	account = bits[0]

	return account

}
