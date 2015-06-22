package main

import (
	"fmt"
)

func (cmds *CommandSet) HelpForPower() {
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
}

func (cmds *CommandSet) Start(args []string) {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	fmt.Printf("Attempting to start %s...\r\n", name.VirtualMachine)
	err := cmds.bigv.StartVirtualMachine(name)
	if err != nil {
		exit(err)
	}

	if !cmds.config.GetBool("silent") {
		fmt.Println(name.VirtualMachine, " started successfully.")
	}
}

func (cmds *CommandSet) Shutdown(args []string) {
	flags := MakeCommonFlagSet()
	restart := flags.Bool("restart", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	fmt.Printf("Attempting to shutdown %s...\r\n", name.VirtualMachine)
	err := cmds.bigv.ShutdownVirtualMachine(name, !*restart)
	if err != nil {
		exit(err)
	}

	if !cmds.config.GetBool("silent") {
		fmt.Println(name.VirtualMachine, " was shutdown successfully.")
	}
}
func (cmds *CommandSet) Stop(args []string) {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	fmt.Printf("Attempting to stop %s...\r\n", name.VirtualMachine)
	err := cmds.bigv.StopVirtualMachine(name)
	if err != nil {
		exit(err)
	}

	if !cmds.config.GetBool("silent") {
		fmt.Println(name.VirtualMachine, " stopped successfully.")
	}
}

func (cmds *CommandSet) Restart(args []string) {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	fmt.Printf("Attempting to restart %s...\r\n", name.VirtualMachine)
	err := cmds.bigv.RestartVirtualMachine(name)
	if err != nil {
		exit(err)
	}

	if !cmds.config.GetBool("silent") {
		fmt.Println(name.VirtualMachine, " restart successfully.")
	}
}

func (cmds *CommandSet) ResetVM(args []string) {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	fmt.Printf("Attempting to reset %s...\r\n", name.VirtualMachine)
	err := cmds.bigv.ResetVirtualMachine(name)
	if err != nil {
		exit(err)
	}

	if !cmds.config.GetBool("silent") {
		fmt.Println(name.VirtualMachine, " reset successfully.")
	}

}
