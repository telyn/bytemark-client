package main

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util/sizespec"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// AccountNameFlag is used for all --account flags, excluding the global one.
type AccountNameFlag string

// Set runs lib.ParseAccountName to make sure we get just the 'pure' account name; no cluster / endpoint details
func (name *AccountNameFlag) Set(value string) error {
	*name = AccountNameFlag(lib.ParseAccountName(value, global.Config.GetIgnoreErr("account")))
	return nil
}

// String returns the AccountNameFlag as a string.
func (name *AccountNameFlag) String() string {
	return string(*name)
}

// GroupNameFlag is used for all --account flags, including the global one.
type GroupNameFlag lib.GroupName

// Set runs lib.ParseGroupName to make sure we have a valid group name
func (name *GroupNameFlag) Set(value string) error {
	gp := lib.ParseGroupName(value, global.Config.GetGroup())
	*name = GroupNameFlag(gp)
	return nil
}

// String returns the GroupNameFlag as a string.
func (name GroupNameFlag) String() string {
	return lib.GroupName(name).String()
}

// VirtualMachineNameFlag is used for all --account flags, including the global one.
type VirtualMachineNameFlag lib.VirtualMachineName

// Set runs lib.ParseVirtualMachineName using the global.Client to make sure we have a valid group name
func (name *VirtualMachineNameFlag) Set(value string) error {
	vm, err := lib.ParseVirtualMachineName(value, global.Config.GetVirtualMachine())
	if err != nil {
		return err
	}
	*name = VirtualMachineNameFlag(vm)
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

// privArgs is an array of strings which can have an argument shifted off the front.
// there are other (better?) ways to do it - just a shift function which takes a *[]string?)
// there's nothing specific about privileges to it either except the error.. perhaps it should be ShiftableStringSlice?
// but again, a func shift(*[]string) (arg string, ok bool) is probably better.
type privArgs []string

func (args *privArgs) shift() (arg string, err error) {
	if len(*args) > 0 {
		arg = (*args)[0]
		*args = (*args)[1:]
	} else {
		err = fmt.Errorf("privileges require 3 parts - level, target and user")
	}
	return
}

// PrivilegeFlag is an un-realised brain.Privilege - where the target name has been parsed but hasn't been turned into IDs yet
type PrivilegeFlag struct {
	AccountName        string
	GroupName          *lib.GroupName
	VirtualMachineName *lib.VirtualMachineName
	Username           string
	Level              brain.PrivilegeLevel
}

// TargetType returns the prefix of the PrivilegeFlag's Level. Use the brain.PrivilegeTargetType* constants for comparison
func (pf PrivilegeFlag) TargetType() string {
	return strings.SplitN(string(pf.Level), "_", 2)[0]
}

// fillPrivilegeTarget adds the object to the privilege, trying to use it as a VM, Group or Account name depending on what PrivilegeLevel the Privilege is for. The target is expected to be the NextArg at this point in the Context
func (pf *PrivilegeFlag) fillPrivilegeTarget(args *privArgs) (err error) {
	if pf.TargetType() != brain.PrivilegeTargetTypeCluster {
		var target string
		target, err = args.shift()
		if err != nil {
			return
		}
		if target == "on" {
			target, err = args.shift()
			if err != nil {
				return
			}
		}
		var vmName lib.VirtualMachineName
		var groupName lib.GroupName
		switch pf.TargetType() {
		case brain.PrivilegeTargetTypeVM:
			vmName, err = lib.ParseVirtualMachineName(target, global.Config.GetVirtualMachine())
			pf.VirtualMachineName = &vmName
		case brain.PrivilegeTargetTypeGroup:
			groupName = lib.ParseGroupName(target, global.Config.GetGroup())
			pf.GroupName = &groupName
		case brain.PrivilegeTargetTypeAccount:
			pf.AccountName = lib.ParseAccountName(target, global.Config.GetIgnoreErr("account"))
		}
	}
	return
}

// Set sets the privilege given some string (should be in the form "<level> [[on] <target>] [to|from|for] <user>"
func (pf *PrivilegeFlag) Set(value string) (err error) {
	args := privArgs(strings.Split(value, " "))

	level, err := args.shift()
	if err != nil {
		return
	}
	pf.Level = brain.PrivilegeLevel(level)

	err = pf.fillPrivilegeTarget(&args)
	if err != nil {
		return
	}

	user, err := args.shift()
	if err != nil {
		return
	}
	if user == "to" || user == "from" || user == "for" {
		user, err = args.shift()
	}
	pf.Username = user

	if arg, err := args.shift(); err == nil {
		return fmt.Errorf("Unexpected '%s' after username '%s'", arg, pf.Username)
	}

	return
}

func (pf PrivilegeFlag) String() string {
	switch pf.TargetType() {
	case brain.PrivilegeTargetTypeVM:
		return fmt.Sprintf("%s on %s for %s", pf.Level, pf.VirtualMachineName, pf.Username)
	case brain.PrivilegeTargetTypeGroup:
		return fmt.Sprintf("%s on %s for %s", pf.Level, pf.GroupName, pf.Username)
	case brain.PrivilegeTargetTypeAccount:
		return fmt.Sprintf("%s on %s for %s", pf.Level, pf.AccountName, pf.Username)
	}
	return fmt.Sprintf("%s for %s", pf.Level, pf.Username)
}
