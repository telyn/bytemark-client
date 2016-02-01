package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
)

func (cmds *CommandSet) HelpForPower() util.ExitCode {
	log.Log("bytemark power commands")
	log.Log()
	log.Log("usage: bytemark start")
	log.Log("       bytemark shutdown")
	log.Log("       bytemark restart")
	log.Log("       bytemark reset")
	log.Log("       bytemark reset")
	log.Log()
	log.Log()
	log.Log("start: Starts a stopped VM.")
	log.Log()
	log.Log("shutdown: Sends the ACPI shutdown signal, as if you had")
	log.Log("          pressed the power/standby button. Allows the")
	log.Log("          operating system to gracefully shut down.")
	log.Log("          Hardware changes will be applied after the")
	log.Log("          machine has been started again.")
	log.Log()
	log.Log("stop: Stops a running VM, as if you had just pulled the")
	log.Log("      cord out. Hardware changes will be applied when the")
	log.Log("      machine has been started again.")
	log.Log()
	log.Log("restart: Stops and then starts a running VM, as if you had")
	log.Log("         pulled the cord out, then plugged it in and")
	log.Log("         powered the machine on again.")
	log.Log()
	log.Log("reset: Instantly restarts a running VM, as if you had")
	log.Log("       pressed the reset button. Doesn't apply hardware")
	log.Log("       changes.")
	log.Log()
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) Start(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Virtual machine name cannnot be blank")
		return util.E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		log.Logf("Attempting to start %s...\r\n", name.VirtualMachine)
	}
	err = cmds.client.StartVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}

	log.Logf("%s started successfully.\r\n", name.VirtualMachine)
	return util.E_SUCCESS
}

func (cmds *CommandSet) Shutdown(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	restart := flags.Bool("restart", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Virtual machine name cannnot be blank")
		return util.E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	log.Logf("Attempting to shutdown %s...\r\n", name.VirtualMachine)

	err = cmds.client.ShutdownVirtualMachine(name, !*restart)
	if err != nil {
		return util.ProcessError(err)
	}

	log.Logf("%s was shutdown successfully.\r\n", name.VirtualMachine)
	return util.E_SUCCESS
}
func (cmds *CommandSet) Stop(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Virtual machine name cannnot be blank")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	log.Logf("Attempting to stop %s...\r\n", name.VirtualMachine)
	err = cmds.client.StopVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}

	log.Logf("%s stopped successfully.\r\n", name.VirtualMachine)
	return util.E_SUCCESS
}

func (cmds *CommandSet) Restart(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Virtual machine name cannnot be blank")
		return util.E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	log.Logf("Attempting to restart %s...\r\n", name.VirtualMachine)
	err = cmds.client.RestartVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}

	log.Logf("%s restarted successfully.\r\n", name.VirtualMachine)
	return util.E_SUCCESS
}

func (cmds *CommandSet) ResetVM(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Virtual machine name cannnot be blank")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	log.Logf("Attempting to reset %s...\r\n", name.VirtualMachine)
	err = cmds.client.ResetVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}

	log.Errorf("%s reset successfully.\r\n", name.VirtualMachine)
	return util.E_SUCCESS

}
