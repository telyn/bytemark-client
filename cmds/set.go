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

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Failed to parse VM name")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	e := cmds.client.SetVirtualMachineHardwareProfileLock(name, true)
	return util.ProcessError(e)
}

// UnlockHWProfile implements the unlock-hwprofile command
func (cmds *CommandSet) UnlockHWProfile(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Failed to parse VM name")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	e := cmds.client.SetVirtualMachineHardwareProfileLock(name, false)
	return util.ProcessError(e)
}

// SetCores implements the set-cores command
func (cmds *CommandSet) SetCores(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	coresStr, ok := util.ShiftArgument(&args, "number of CPU cores")
	if !ok {
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Errorf("Failed to parse VM name\r\n")
		return util.E_PEBKAC
	}

	// decide on the number of cores to set now
	cores, err := strconv.Atoi(coresStr)
	if err != nil || cores < 1 {
		log.Errorf("Invalid number of cores \"%s\"\r\n", coresStr)
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	// submit command to API
	err = cmds.client.SetVirtualMachineCores(name, cores)
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

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForSet()
		return util.E_PEBKAC
	}
	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())

	profileStr, ok := util.ShiftArgument(&args, "hardware profile")
	if !ok {
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

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
		err = cmds.client.SetVirtualMachineHardwareProfile(name, profileStr, true)
	} else if *unlock_hwp {
		err = cmds.client.SetVirtualMachineHardwareProfile(name, profileStr, false)
		// otherwise omit lock
	} else {
		err = cmds.client.SetVirtualMachineHardwareProfile(name, profileStr)
	}

	// return
	return util.ProcessError(err)
}

// SetMemory implements the set-memory command
func (cmds *CommandSet) SetMemory(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Failed to parse VM name")
		return util.E_PEBKAC
	}

	memoryStr, ok := util.ShiftArgument(&args, "memory size")
	if !ok {
		cmds.HelpForSet()
		return util.E_PEBKAC
	}

	memory, err := util.ParseSize(memoryStr)
	if err != nil || memory < 1 {
		log.Errorf("Invalid amount of memory \"%s\"\r\n", memoryStr)
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.client.SetVirtualMachineMemory(name, memory)
	return util.ProcessError(err)
}
