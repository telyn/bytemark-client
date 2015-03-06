package cmd

import "strings"

func FirstNotEmpty(choices ...string) string {
	for _, choice := range choices {
		if choice != "" {
			return choice

		}
	}
	return ""
}

// Represents an account reference (either the account name or number)
// name should be in lowercase, if it is "DEFAULT" then the account from config will be used.
type AccountId struct {
	name string
}

// Represents an account reference (either the account name or number)
// name should be in lowercase, if it is "DEFAULT" then the account from config will be used.

type GroupId struct {
	account string
	name    string
}

type VirtualMachineId struct {
	account string
	group   string
	name    string
}

func ParseVirtualMachineName(name string) (vmId VirtualMachineId) {
	// 1, 2 or 3 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	vmId.group = "DEFAULT"
	vmId.account = "DEFAULT"
	vmId.name = "DEFAULT"

	// a for loop seems an odd choice here maybe but it means
	// I don't need to do lots of ifs to see if the next bit exists
Loop:
	for i, bit := range bits {
		switch i {
		case 0:
			vmId.name = strings.ToLower(bit)
			break
		case 1:
			// want to be able to do vm..account to get the default group
			if bit == "" {
				break
			}
			vmId.group = strings.ToLower(bit)
			break
		case 2:
			vmId.account = strings.ToLower(bit)
			break Loop
		}
	}
	return vmId
}

func ParseGroupName(name string) (groupId GroupId) {
	// 1 or 2 pieces with optional extra cruft for the fqdn
	bits := strings.Split(name, ".")
	groupId.name = "DEFAULT"
	groupId.account = "DEFAULT"

Loop:
	for i, bit := range bits {
		switch i {
		case 0:
			groupId.name = strings.ToLower(bit)
			break
		case 1:
			groupId.account = strings.ToLower(bit)
			break Loop
		}
	}
	return groupId

}

func ParseAccountName(name string) (accountId AccountId) {
	// 1 piece with optional extra cruft for the fqdn

	// there's a micro-optimisation to do here to not use Split,
	// but really, who can be bothered to read the documentation?
	accountId.name = "DEFAULT"

	bits := strings.Split(name, ".")
	accountId.name = bits[0]

	return accountId

}

