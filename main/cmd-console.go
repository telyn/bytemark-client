package main

import (
	"bigv.io/client/lib"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func (cmds *CommandSet) HelpForConsole() {
	fmt.Println("go-bigv console commands")
	fmt.Println()
	fmt.Println("usage: go-bigv console [--serial | --vnc] [--connect | --panel] <virtual machine>")
	fmt.Println("       go-bigv serial [--connect] <virtual machine>")
	fmt.Println("       go-bigv vnc [--connect | --panel] <virtual machine>")
	fmt.Println()
	fmt.Println("Out-of-band access to a machine's serial or graphical (VNC) console.")
	fmt.Println()
	fmt.Println("serial: outputs connection information for the out-of-band")
	fmt.Println("        serial console specified. If the --connect flag is")
	fmt.Println("        given, will attempt to open a connection as well.")
	fmt.Println()
	fmt.Println("vnc: Outputs instructions for connecting to the VNC console.")
	fmt.Println("     If --connect is given, will also attempt to connect to")
	fmt.Println("     the VNC console, falling back to the BigV panel if no ssh")
	fmt.Println("     or vnc clients can be found, or if --panel is specified.")
	fmt.Println()
	fmt.Println()
	//TODO: stop kidding around
	fmt.Println("haha just kidding, at the moment vnc always uses panel")
}

func shortEndpoint(endpoint string) string {
	return strings.Split(endpoint, ".")[0]
}

func getExitCode(cmd *exec.Cmd) (exitCode int, err error) {
	err = cmd.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitErr.Sys().(*syscall.WaitStatus); ok {
			return waitStatus.ExitStatus(), err

		}
	}
	return 0, err
}

func showVNCHowTo(vm *lib.VirtualMachine) {
	fmt.Println("VNC connection information for", vm.Hostname)
	fmt.Println()
	fmt.Printf("Ensure that your public key (contained in %s/.ssh/id_rsa.pub or %s/.ssh/id_dsa.pub) is present in your bigv user's keys (see `bigv show keys`, `bigv add key`)", os.Getenv("HOME"), os.Getenv("HOME"))
	fmt.Println()
	fmt.Println("Then set up a tunnel using SSH: ssh -L <some number>:%s:5900 %s@%s\r\n")
	fmt.Println()
	fmt.Println("You can now connect to VNC on localhost, port <some number>")
}

func (cmds *CommandSet) showSSHHowTo(vm *lib.VirtualMachine) {
	fmt.Println("Serial console connection information for", vm.Hostname)
	fmt.Println()
	fmt.Printf("Ensure that your public key (contained in %s/.ssh/id_rsa.pub or %s/.ssh/id_dsa.pub) is present in your bigv user's keys (see `bigv show keys`, `bigv add key`)", os.Getenv("HOME"), os.Getenv("HOME"))
	fmt.Println()
	fmt.Printf("Then connect to %s@%s\r\n", cmds.bigv.GetSessionUser(), vm.ManagementAddress)

}

//TODO: this function is really horrible.
func (cmds *CommandSet) Console(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	connect := flags.Bool("connect", false, "")
	//panel := flags.Bool("panel", false, "")
	flags.Bool("serial", false, "") // because we default to serial, we don't care if it's set
	vnc := flags.Bool("vnc", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := ShiftArgument(args, "virtual machine")
	if !ok {
		cmds.HelpForConsole()
		return E_PEBKAC
	}
	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannot be blank\r\n")
		return E_PEBKAC
	}
	cmds.EnsureAuth()

	vm, err := cmds.bigv.GetVirtualMachine(name)
	if err != nil {
		return processError(err)
	}
	if *vnc {

		if !*connect {
			showVNCHowTo(vm)
			return E_SUCCESS
		}
		ep := cmds.config.EndpointName()
		token := cmds.config.GetIgnoreErr("token")
		url := fmt.Sprintf("%s/vnc/?auth_token=%s&endpoint=%s&management_ip=%s", cmds.config.PanelURL(), token, shortEndpoint(ep), vm.ManagementAddress)

		return processError(callBrowser(url))

	} else { // default to serial
		if !*connect {
			cmds.showSSHHowTo(vm)
			return E_SUCCESS
		}
		host := fmt.Sprintf("%s@%s", cmds.bigv.GetSessionUser(), vm.ManagementAddress)
		fmt.Fprintf(os.Stderr, "ssh %s\r\n", host)
		bin, err := exec.LookPath("ssh")
		if err != nil {
			return processError(err, "Unable to find an ssh executable")
		}
		err = syscall.Exec(bin, []string{"ssh", host}, os.Environ())
		if err != nil {
			if errno, ok := err.(syscall.Errno); ok {
				if errno != 0 {
					return processError(err, "Attempting to exec ssh failed. Please file a bug report.")
				}
			} else {
				return processError(err, "Couldn't connect to the management address. Please ensure you have an SSH client in your $PATH.")

			}
		}
	}
	return E_SUCCESS

}
