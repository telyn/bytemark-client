package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util/sizespec"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"strings"
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

// ResizeMode represents whether to increment a size or just to set it.
type ResizeMode int

const (
	// ResizeModeSet will cause resize disk to set the disc size to the one specified
	ResizeModeSet = iota
	// ResizeModeIncrease will cause resize disk to increase the disc size by the one specified
	ResizeModeIncrease
)

// ResizeFlag is effectively an extension of SizeSpecFlag which has a ResizeMode. The Size stored in the flag is the size to set to or increase by depending on the Mode
type ResizeFlag struct {
	Mode ResizeMode
	Size int
}

// Set parses the string into a ResizeFlag. If it starts with +, Mode is set to ResizeModeIncrease. Otherwise, it's set to ResizeModeSet. The rest of the string is parsed as a sizespec using sizespec.Parse
func (rf *ResizeFlag) Set(value string) (err error) {
	rf.Mode = ResizeModeSet
	if strings.HasPrefix(value, "+") {
		rf.Mode = ResizeModeIncrease
		value = value[1:]
	}

	sz, err := sizespec.Parse(value)
	if err != nil {
		return
	}
	rf.Size = sz
	return
}

// String returns the size, in GiB or TiB (if the size is > 1TIB) with the unit used as a suffix. If Mode is ResizeModeIncrease, the string is prefixed with '+'
func (rf ResizeFlag) String() string {
	plus := ""
	if rf.Mode == ResizeModeIncrease {
		plus += "+"
	}
	sz := rf.Size
	units := "GiB"
	sz /= 1024
	if sz > 1024 {
		sz /= 1024
		units = "TiB"
	}
	return fmt.Sprintf("%s%d%s", plus, sz, units)
}
