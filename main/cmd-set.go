package main

import (
	//bigv "bigv.io/client/lib"
	"fmt"
	"os"
	"strconv"
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

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse VM name\r\n")
		return E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	e := cmds.bigv.SetVirtualMachineHardwareProfileLock(name, true)
	return processError(e)
}

// UnlockHWProfile implements the unlock-hwprofile command
func (cmds *CommandSet) UnlockHWProfile(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse VM name\r\n")
		return E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	e := cmds.bigv.SetVirtualMachineHardwareProfileLock(name, false)
	return processError(e)
}

// SetCores implements the set-cores command
func (cmds *CommandSet) SetCores(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	if len(args) != 2 {
		fmt.Println("must specify a VM name and a number of CPUs")
		cmds.HelpForSet()
		return E_PEBKAC
	}

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse VM name\r\n")
		return E_PEBKAC
	}

	// decide on the number of cores to set now
	cores, err := strconv.Atoi(args[1])
	if err != nil || cores < 1 {
		fmt.Fprintf(os.Stderr, "Invalid number of cores \"%s\"\r\n", args[1])
		return E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	// submit command to BigV
	err = cmds.bigv.SetVirtualMachineCores(name, cores)
	return processError(err)
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

	// name and hardware profile required
	if len(args) != 2 {
		fmt.Println("must specify a VM name and a hardware profile")
		cmds.HelpForSet()
		return E_PEBKAC
	}
	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse VM name\r\n")
		return E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	// if lock_hwp or unlock_hwp are specified, account this into the call
	if *lock_hwp {
		err = cmds.bigv.SetVirtualMachineHardwareProfile(name, args[1], true)
	} else if *unlock_hwp {
		err = cmds.bigv.SetVirtualMachineHardwareProfile(name, args[1], false)
		// otherwise omit lock
	} else {
		err = cmds.bigv.SetVirtualMachineHardwareProfile(name, args[1])
	}

	// return
	return processError(err)
}

// SetMemory implements the set-memory command
func (cmds *CommandSet) SetMemory(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	if len(args) != 2 {
		fmt.Println("must specify a VM name and an amount of memory")
		cmds.HelpForSet()
		return E_PEBKAC
	}

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse VM name\r\n")
		return E_PEBKAC
	}

	memory, err := ParseSize(args[1])
	if err != nil || memory < 1 {
		fmt.Fprintf(os.Stderr, "Invalid amount of memory \"%s\"\r\n", args[1])
		return E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	err = cmds.bigv.SetVirtualMachineMemory(name, memory)
	return processError(err)
}
