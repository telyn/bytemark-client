package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// AccountNameFlag is used for all --account flags, including the global one.
type AccountNameFlag string

// Set runs lib.Client.ParseAccountName using the global.Client to make sure we get just the 'pure' account name; no cluster / endpoint details
func (name *AccountNameFlag) Set(value string) error {
	*name = AccountNameFlag(global.Client.ParseAccountName(value, global.Config.GetIgnoreErr("account")))
	return nil
}

// String returns the AccountNameFlag as a string.
func (name *AccountNameFlag) String() string {
	return string(*name)
}

// GroupNameFlag is used for all --account flags, including the global one.
type GroupNameFlag lib.GroupName

// Set runs lib.Client.ParseGroupName using the global.Client to make sure we have a valid group name
func (name *GroupNameFlag) Set(value string) error {
	gp := global.Client.ParseGroupName(value, global.Config.GetGroup())
	*name = GroupNameFlag(*gp)
	return nil
}

// String returns the GroupNameFlag as a string.
func (name GroupNameFlag) String() string {
	return lib.GroupName(name).String()
}

// VirtualMachineNameFlag is used for all --account flags, including the global one.
type VirtualMachineNameFlag lib.VirtualMachineName

// Set runs lib.Client.ParseVirtualMachineName using the global.Client to make sure we have a valid group name
func (name *VirtualMachineNameFlag) Set(value string) error {
	vm, err := global.Client.ParseVirtualMachineName(value, global.Config.GetVirtualMachine())
	if err != nil {
		return err
	}
	*name = VirtualMachineNameFlag(*vm)
	return nil
}

// String returns the VirtualMachineNameFlag as a string.
func (name VirtualMachineNameFlag) String() string {
	return lib.VirtualMachineName(name).String()
}
