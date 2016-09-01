package util

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// EnsureAccountName used to be how we defaulted an account. It's no longer used and will be removed when I do a deadcode delinting.
func EnsureAccountName(account string, config ConfigManager) string {
	if account == "" {
		return config.GetIgnoreErr("account")
	}
	return account
}

// EnsureGroupName used to be how we defaulted a group name. It's no longer used and will be removed when I do a deadcode delinting.
func EnsureGroupName(group lib.GroupName, config ConfigManager) lib.GroupName {
	group.Account = EnsureAccountName(group.Account, config)
	if group.Group == "" {
		group.Group = config.GetIgnoreErr("group")
	}
	return group
}

// EnsureVMName used to be how we defaulted a vm name. It's no longer used and will be removed when I do a deadcode delinting.
func EnsureVMName(vm lib.VirtualMachineName, config ConfigManager) lib.VirtualMachineName {
	vm.Account = EnsureAccountName(vm.Account, config)
	if vm.Group == "" {
		vm.Group = config.GetIgnoreErr("group")
	}
	// nowt we can actually do about a blank vm name.
	return vm
}
