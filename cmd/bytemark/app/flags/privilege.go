package flags

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

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
	GroupName          *pathers.GroupName
	VirtualMachineName *lib.VirtualMachineName
	Username           string
	Level              brain.PrivilegeLevel
	Value              string
}

// TargetType returns the prefix of the Privilege's Level. Use the brain.PrivilegeTargetType* constants for comparison
func (pf PrivilegeFlag) TargetType() string {
	return strings.SplitN(string(pf.Level), "_", 2)[0]
}

// fillPrivilegeTarget adds the object to the privilege, trying to use it as a VM, Group or Account name depending on what PrivilegeLevel the Privilege is for. The target is expected to be the NextArg at this point in the Context
func (pf *PrivilegeFlag) fillPrivilegeTarget(c *app.Context, args *privArgs) (err error) {
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
		var groupName pathers.GroupName
		switch pf.TargetType() {
		case brain.PrivilegeTargetTypeVM:
			vmName, err = lib.ParseVirtualMachineName(target, c.Config().GetVirtualMachine())
			pf.VirtualMachineName = &vmName
		case brain.PrivilegeTargetTypeGroup:
			groupName = lib.ParseGroupName(target, c.Config().GetGroup())
			pf.GroupName = &groupName
		case brain.PrivilegeTargetTypeAccount:
			pf.AccountName = lib.ParseAccountName(target, c.Config().GetIgnoreErr("account"))
		}
	}
	return
}

// Set sets the privilege given some string (should be in the form "<level> [[on] <target>] [to|from|for] <user>"
func (pf *PrivilegeFlag) Set(value string) (err error) {
	pf.Value = value
	return nil
}

// Preprocess parses the Privilege and looks up the target's ID so it can
// be made into a brain.Privilege
// This is an implementation of `app.Preprocessor`, which is detected and
// called automatically by actions created with `app.Action`
func (pf *PrivilegeFlag) Preprocess(c *app.Context) (err error) {
	args := privArgs(strings.Split(pf.Value, " "))

	level, err := args.shift()
	if err != nil {
		return
	}
	pf.Level = brain.PrivilegeLevel(level)

	err = pf.fillPrivilegeTarget(c, &args)
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
