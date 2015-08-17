package main

import (
	"fmt"
	"os"
)

func (cmds *CommandSet) HelpForPower() ExitCode {
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
	return E_USAGE_DISPLAYED
}

func (cmds *CommandSet) Start(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to start %s...\r\n", name.VirtualMachine)
	}
	err = cmds.bigv.StartVirtualMachine(name)
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " started successfully.")
	}
	return E_SUCCESS
}

func (cmds *CommandSet) Shutdown(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	restart := flags.Bool("restart", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to shutdown %s...\r\n", name.VirtualMachine)
	}

	err = cmds.bigv.ShutdownVirtualMachine(name, !*restart)
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " was shutdown successfully.")
	}
	return E_SUCCESS
}
func (cmds *CommandSet) Stop(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to stop %s...\r\n", name.VirtualMachine)
	}
	err = cmds.bigv.StopVirtualMachine(name)
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " stopped successfully.")
	}
	return E_SUCCESS
}

func (cmds *CommandSet) Restart(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to restart %s...\r\n", name.VirtualMachine)
	}
	err = cmds.bigv.RestartVirtualMachine(name)
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " restart successfully.")
	}
	return E_SUCCESS
}

func (cmds *CommandSet) ResetVM(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Printf("Attempting to reset %s...\r\n", name.VirtualMachine)
	}
	err = cmds.bigv.ResetVirtualMachine(name)
	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {
		fmt.Println(name.VirtualMachine, " reset successfully.")
	}
	return E_SUCCESS

}
