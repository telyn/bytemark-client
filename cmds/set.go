package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
	"strconv"
)

// HelpForLocking provides usage information for locking and unlocking hardware
// profiles.
func (cmds *CommandSet) HelpForLocks() util.ExitCode {
	return util.E_SUCCESS
}

// HelpForSet provides usage information for the set command and its subcommands.
func (cmds *CommandSet) HelpForSet() util.ExitCode {
	return util.E_SUCCESS
}

// LockHWProfile implements the lock-hwprofile command
func (cmds *CommandSet) LockHWProfile(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		log.Error("Failed to parse VM name")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	e := cmds.bigv.SetVirtualMachineHardwareProfileLock(name, true)
	return util.ProcessError(e)
}

// UnlockHWProfile implements the unlock-hwprofile command
func (cmds *CommandSet) UnlockHWProfile(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		log.Error("Failed to parse VM name")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	e := cmds.bigv.SetVirtualMachineHardwareProfileLock(name, false)
	return util.ProcessError(e)
}

// SetCores implements the set-cores command
func (cmds *CommandSet) SetCores(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	if len(args) != 2 {
		log.Log("must specify a VM name and a number of CPUs")
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		log.Errorf("Failed to parse VM name\r\n")
		return util.E_PEBKAC
	}

	// decide on the number of cores to set now
	cores, err := strconv.Atoi(args[1])
	if err != nil || cores < 1 {
		log.Errorf("Invalid number of cores \"%s\"\r\n", args[1])
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	// submit command to BigV
	err = cmds.bigv.SetVirtualMachineCores(name, cores)
	return util.ProcessError(err)
}

// SetHWProfile implements the set-hwprofile command
func (cmds *CommandSet) SetHWProfile(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	lock_hwp := flags.Bool("lock", false, "")
	unlock_hwp := flags.Bool("unlock", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	// do nothing if --lock and --unlock are both specified
	if *lock_hwp && *unlock_hwp {
		log.Log("Ambiguous command, both lock and unlock specified")
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	// name and hardware profile required
	if len(args) != 2 {
		log.Log("Must specify a VM name and a hardware profile")
		cmds.HelpForSet()
		return util.E_PEBKAC
	}
	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		log.Error("Failed to parse VM name")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
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
	return util.ProcessError(err)
}

// SetMemory implements the set-memory command
func (cmds *CommandSet) SetMemory(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	if len(args) != 2 {
		log.Log("Must specify a VM name and an amount of memory")
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		log.Error("Failed to parse VM name")
		return util.E_PEBKAC
	}

	memory, err := util.ParseSize(args[1])
	if err != nil || memory < 1 {
		log.Errorf("Invalid amount of memory \"%s\"\r\n", args[1])
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.bigv.SetVirtualMachineMemory(name, memory)
	return util.ProcessError(err)
}
