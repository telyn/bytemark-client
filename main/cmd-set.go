package main

import (
	//bigv "bigv.io/client/lib"
	"fmt"
)

// HelpForLocking provides usage information for locking and unlocking hardware
// profiles.
func (cmds *CommandSet) HelpForLocks() ExitCode {
	return E_SUCCESS
}

// HelpForSet provides usage information for the set command and its subcommands.
func (cmds *CommandSet) HelpForSet() ExitCode {
	return E_SUCCESS
}

// LockHWProfile implements the lock-hwprofile command
func (cmds *CommandSet) LockHWProfile(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	e := cmds.bigv.SetVirtualMachineHardwareProfileLock(name, true)
	return processError(e)
}

// UnlockHWProfile implements the unlock-hwprofile command
func (cmds *CommandSet) UnlockHWProfile(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name := cmds.bigv.ParseVirtualMachineName(args[0])
	cmds.EnsureAuth()

	e := cmds.bigv.SetVirtualMachineHardwareProfileLock(name, false)
	return processError(e)
}

// SetHWProfile implements the set-hwprofile command
func (cmds *CommandSet) SetHWProfile(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	lock_hwp := flags.Bool("lock", false, "")
	unlock_hwp := flags.Bool("unlock", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	// do nothing if --lock and --unlock are both specified
	if *lock_hwp && *unlock_hwp {
		fmt.Println("ambiguous command, both lock and unlock specified")
		cmds.HelpForSet()
		return E_PEBKAC
	}

	// identify vm
	var e error

	// name and hardware profile required
	if len(args) < 2 {
		fmt.Println("must specify a VM name and a hardware profile")
		cmds.HelpForSet()
		return E_PEBKAC
	}
	name := cmds.bigv.ParseVirtualMachineName(args[0])

	cmds.EnsureAuth()

	// if lock_hwp or unlock_hwp are specified, account this into the call
	if *lock_hwp {
		e = cmds.bigv.SetVirtualMachineHardwareProfile(name, args[1], true)
	} else if *unlock_hwp {
		e = cmds.bigv.SetVirtualMachineHardwareProfile(name, args[1], false)
		// otherwise omit lock
	} else {
		e = cmds.bigv.SetVirtualMachineHardwareProfile(name, args[1])
	}

	// return
	return processError(e)
}
