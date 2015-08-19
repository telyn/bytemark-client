package cmds

import (
	"bigv.io/client/cmds/util"
	"fmt"
	"os"
)

func (cmds *CommandSet) HelpForPower() util.ExitCode {
	fmt.Println("go-bigv power commands")
	fmt.Println()
	fmt.Println("usage: go-bigv start")
	fmt.Println("       go-bigv shutdown")
	fmt.Println("       go-bigv restart")
	fmt.Println("       go-bigv reset")
	fmt.Println("       go-bigv reset")
	fmt.Println()
	fmt.Println()
	fmt.Println("start: Starts a stopped VM.")
	fmt.Println()
	fmt.Println("shutdown: Sends the ACPI shutdown signal, as if you had")
	fmt.Println("          pressed the power/standby button. Allows the")
	fmt.Println("          operating system to gracefully shut down.")
	fmt.Println("          Hardware changes will be applied after the")
	fmt.Println("          machine has been started again.")
	fmt.Println()
	fmt.Println("stop: Stops a running VM, as if you had just pulled the")
	fmt.Println("      cord out. Hardware changes will be applied when the")
	fmt.Println("      machine has been started again.")
	fmt.Println()
	fmt.Println("restart: Stops and then starts a running VM, as if you had")
	fmt.Println("         pulled the cord out, then plugged it in and")
	fmt.Println("         powered the machine on again.")
	fmt.Println()
	fmt.Println("reset: Instantly restarts a running VM, as if you had")
	fmt.Println("       pressed the reset button. Doesn't apply hardware")
	fmt.Println("       changes.")
	fmt.Println()
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

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return util.E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to start %s...\r\n", name.VirtualMachine)
	}
	err = cmds.bigv.StartVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " started successfully.")
	}
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

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return util.E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to shutdown %s...\r\n", name.VirtualMachine)
	}

	err = cmds.bigv.ShutdownVirtualMachine(name, !*restart)
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " was shutdown successfully.")
	}
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

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to stop %s...\r\n", name.VirtualMachine)
	}
	err = cmds.bigv.StopVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " stopped successfully.")
	}
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

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return util.E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to restart %s...\r\n", name.VirtualMachine)
	}
	err = cmds.bigv.RestartVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " restart successfully.")
	}
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

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to reset %s...\r\n", name.VirtualMachine)
	}
	err = cmds.bigv.ResetVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " reset successfully.")
	}
	return util.E_SUCCESS

}
