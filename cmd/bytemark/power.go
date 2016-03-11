package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
)

func (cmds *CommandSet) HelpForPower() util.ExitCode {
	log.Log("bytemark power commands")
	log.Log()
	log.Log("usage: bytemark start <server>")
	log.Log("       bytemark shutdown <server>")
	log.Log("       bytemark restart <server>")
	log.Log("       bytemark reset <server>")
	log.Log()
	log.Log()
	log.Log("start: Starts a stopped server.")
	log.Log()
	log.Log("shutdown: Sends the ACPI shutdown signal, as if you had")
	log.Log("          pressed the power/standby button. Allows the")
	log.Log("          operating system to gracefully shut down.")
	log.Log("          Hardware changes will be applied after the")
	log.Log("          machine has been started again.")
	log.Log()
	log.Log("stop: Stops a running server, as if you had just pulled the")
	log.Log("      cord out. Hardware changes will be applied when the")
	log.Log("      machine has been started again.")
	log.Log()
	log.Log("restart: Stops and then starts a running server, as if you had")
	log.Log("         pulled the cord out, then plugged it in and")
	log.Log("         powered the machine on again.")
	log.Log()
	log.Log("reset: Instantly restarts a running server, as if you had")
	log.Log("       pressed the reset button. Doesn't apply hardware")
	log.Log("       changes.")
	log.Log()
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) Start(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)
	if len(args) > 1 {
		log.Log("Too many arguments to `bytemark start` - please specify only a single server.")
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("server name cannnot be blank")
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

	if len(args) > 1 {
		log.Log("Too many arguments to `bytemark shutdown` - please specify only a single server.")
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("server name cannnot be blank")
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

	if len(args) > 1 {
		log.Log("Too many arguments to `bytemark stop` - please specify only a single server.")
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("server name cannnot be blank")
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

	if len(args) > 1 {
		log.Log("Too many arguments to `bytemark restart` - please specify only a single server.")
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("server name cannnot be blank")
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

func (cmds *CommandSet) ResetServer(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	if len(args) > 1 {
		log.Log("Too many arguments to `bytemark reset` - please specify only a single server.")
		cmds.HelpForPower()
		return util.E_PEBKAC
	}
	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForPower()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("server name cannnot be blank")
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
