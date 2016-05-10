package util

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
)

func EnsureAccountName(account string, config ConfigManager) string {
	if account == "" {
		return config.GetIgnoreErr("account")
	}
	return account
}

func EnsureGroupName(group lib.GroupName, config ConfigManager) lib.GroupName {
	group.Account = EnsureAccountName(group.Account, config)
	if group.Group == "" {
		group.Group = config.GetIgnoreErr("group")
	}
	return group
}

func EnsureVMName(vm lib.VirtualMachineName, config ConfigManager) lib.VirtualMachineName {
	vm.Account = EnsureAccountName(vm.Account, config)
	if vm.Group == "" {
		vm.Group = config.GetIgnoreErr("group")
	}
	// nowt we can actually do about a blank vm name.
	return vm
}
