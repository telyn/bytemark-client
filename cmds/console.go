package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func (cmds *CommandSet) HelpForConsole() util.ExitCode {
	log.Log("bytemark console commands")
	log.Log()
	log.Log("usage: bytemark console [--serial | --vnc] [--connect | --panel] <server>")
	log.Log("       bytemark serial [--connect] <server>")
	log.Log("       bytemark vnc [--connect | --panel] <cloud server>")
	log.Log()
	log.Log("Out-of-band access to a machine's serial or graphical (VNC) console.")
	log.Log()
	log.Log("serial: outputs connection information for the out-of-band")
	log.Log("        serial console specified. If the --connect flag is")
	log.Log("        given, will attempt to open a connection as well.")
	log.Log()
	log.Log("vnc: Outputs instructions for connecting to the VNC console.")
	log.Log("     If --connect is given, will also attempt to connect to")
	log.Log("     the VNC console using the Bytemark panel.")
	log.Log()
	log.Log()
	return util.E_USAGE_DISPLAYED
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

func (cmds *CommandSet) showVNCHowTo(vm *lib.VirtualMachine) {
	log.Log("VNC connection information for", vm.Hostname)
	log.Log()
	log.Logf("Ensure that your public key (contained in %s/.ssh/id_rsa.pub or %s/.ssh/id_dsa.pub) is present in your Bytemark user's keys (see `bytemark show keys`, `bytemark add key`)", os.Getenv("HOME"), os.Getenv("HOME"))
	log.Log()
	log.Logf("Then set up a tunnel using SSH: ssh -L 9999:%s:5900 %s@%s\r\n", vm.ManagementAddress, cmds.client.GetSessionUser(), vm.ManagementAddress)
	log.Log()
	log.Log("You will then be able to connect to vnc://localhost:9999/")
	log.Log("Any port may be substituted for 9999 as long as the same port is used in both commands")
}

func (cmds *CommandSet) showSSHHowTo(vm *lib.VirtualMachine) {
	log.Log("Serial console connection information for", vm.Hostname)
	log.Log()
	log.Logf("Ensure that your public key (contained in %s/.ssh/id_rsa.pub or %s/.ssh/id_dsa.pub) is present in your Bytemark user's keys (see `bytemark show keys`, `bytemark add key`)", os.Getenv("HOME"), os.Getenv("HOME"))
	log.Log()
	log.Logf("Then connect to %s@%s\r\n", cmds.client.GetSessionUser(), vm.ManagementAddress)

}

//TODO: this function is really horrible.
func (cmds *CommandSet) Console(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	connect := flags.Bool("connect", false, "")
	//panel := flags.Bool("panel", false, "")
	flags.Bool("serial", false, "") // because we default to serial, we don't care if it's set
	vnc := flags.Bool("vnc", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForConsole()
		return util.E_PEBKAC
	}
	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Logf("server name cannot be blank\r\n")
		return util.E_PEBKAC
	}
	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	vm, err := cmds.client.GetVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}
	if *vnc {

		if !*connect {
			cmds.showVNCHowTo(vm)
			return util.E_SUCCESS
		}
		ep := cmds.config.EndpointName()
		token := cmds.config.GetIgnoreErr("token")
		url := fmt.Sprintf("%s/vnc/?auth_token=%s&endpoint=%s&management_ip=%s", cmds.config.PanelURL(), token, shortEndpoint(ep), vm.ManagementAddress)
		err = util.CallBrowser(url)
		return util.ProcessError(err)

	} else { // default to serial
		if !*connect {
			cmds.showSSHHowTo(vm)
			return util.E_SUCCESS
		}
		host := fmt.Sprintf("%s@%s", cmds.client.GetSessionUser(), vm.ManagementAddress)
		log.Logf("ssh %s\r\n", host)
		bin, err := exec.LookPath("ssh")
		if err != nil {
			return util.ProcessError(err, "Unable to find an ssh executable")
		}
		err = syscall.Exec(bin, []string{"ssh", host}, os.Environ())
		if err != nil {
			if errno, ok := err.(syscall.Errno); ok {
				if errno != 0 {
					return util.ProcessError(err, "Attempting to exec ssh failed. Please file a bug report.")
				}
			} else {
				return util.ProcessError(err, "Couldn't connect to the management address. Please ensure you have an SSH client in your $PATH.")

			}
		}
	}
	return util.E_SUCCESS

}
