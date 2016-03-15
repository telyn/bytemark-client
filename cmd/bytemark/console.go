package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func init() {
	commands = append(commands, cli.Command{
		Name:  "console",
		Usage: "bytemark console [--serial | --vnc | --panel] [--no-connect] <server>",
		UsageText: `Out-of-band access to a server's serial or graphical (VNC) console.
Under systems with no GUI, sometimes errors will output to the graphical console and not the serial console.
Defaults to connecting to the serial console for the given server.`,
		/*		Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "serial",
					Usage: "Connect to the serial console. Cannot be set at same time as  --vnc or --panel. You must have an ssh client on your computer.",
				},
				cli.BoolFlag{
					Name:  "vnc",
					Usage: "Connect to the graphical console. Cannot be set at the same time as --serial or --panel. You must have a openssh client and a VNC client set up on your computer.",
				},
				cli.BoolFlag{
				    Name: "panel",
				    Usage: "Connect to the graphical console via the Bytemark panel. Cannot be set at the same time as --serial or --vnc. You must have a graphical browser installed on your computer."
				},
				cli.BoolFlag{
				    Name: "no-connect",
				    Usage: "Output connection instructions, rather than directly connecting.",
				},
			},*/
		Action: func(ctx *cli.Context) {
			flags := util.MakeCommonFlagSet()
			panel := flags.Bool("panel", false, "")
			serial := flags.Bool("serial", false, "") // because we default to serial, we don't care if it's set
			no_connect := flags.Bool("no-connect", false, "")
			flags.Parse(ctx.Args())
			args := global.Config.ImportFlags(flags)

			if *serial && *panel {
				log.Logf("You must only specify one of --serial and --panel!")
				global.Error = util.PEBKACError{}
				return
			}

			nameStr, ok := util.ShiftArgument(&args, "server")
			if !ok {
				global.Error = util.PEBKACError{}
				return
			}
			name, err := global.Client.ParseVirtualMachineName(nameStr, global.Config.GetVirtualMachine())
			if err != nil {
				log.Logf("server name cannot be blank\r\n")
				global.Error = util.PEBKACError{}
				return
			}
			err = EnsureAuth()
			if err != nil {
				global.Error = err
				return
			}

			vm, err := global.Client.GetVirtualMachine(name)
			if err != nil {
				global.Error = err
				return
			}
			if *no_connect {
				console_serial_instructions(vm)
				log.Log()
				console_vnc_instructions(vm)
				return
			} else {
				if *panel {
					global.Error = console_panel(vm)
					return
				} else {
					global.Error = console_serial(vm)
					return
				}
			}

		},
	})
	// TODO(telyn): decide whether to keep serial and panel commands
	/*commands = append(commands, cli.Command{
			Name:      "serial",
			Usage:     "bytemark serial <server>",
			UsageText: "Out-of-band access to a server's serial console.",
			Description: `outputs connection information for the out-of-band serial console
	specified. If the --connect flag is given, will attempt to open a connection as well.`,
			Action: console_serial,
		})*/
	/*commands = append(commands, cli.Command{
			Name:      "vnc",
			Usage:     "bytemark vnc <server>", // [--connect | --panel]
			UsageText: "Out-of-band access to a server's graphical console.",
			Description: `Outputs instructions for connecting to the VNC console.
	If --connect is given, will also attempt to connect to the VNC console using the Bytemark panel.`,
			Action: console_serial,
		})*/
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

func console_vnc_instructions(vm *lib.VirtualMachine) {
	mgmtAddress := vm.ManagementAddress.String()
	if vm.ManagementAddress.To4() == nil {
		mgmtAddress = "[" + mgmtAddress + "]"
	}

	log.Log("VNC graphical console connection information for", vm.Hostname)
	log.Log()
	log.Logf("Ensure that your public key (contained in %s/.ssh/id_rsa.pub or %s/.ssh/id_dsa.pub) is present in your Bytemark user's keys (see `bytemark show keys`, `bytemark add key`)", os.Getenv("HOME"), os.Getenv("HOME"))
	log.Log()
	log.Logf("Then set up a tunnel using SSH: ssh -L 9999:%s:5900 %s@%s\r\n", mgmtAddress, global.Client.GetSessionUser(), mgmtAddress)
	log.Log()
	log.Log("You will then be able to connect to vnc://localhost:9999/")
	log.Log("Any port may be substituted for 9999 as long as the same port is used in both commands")
}

func console_serial_instructions(vm *lib.VirtualMachine) {
	mgmtAddress := vm.ManagementAddress.String()
	if vm.ManagementAddress.To4() == nil {
		mgmtAddress = "[" + mgmtAddress + "]"
	}
	log.Log("Serial console connection information for", vm.Hostname)
	log.Log()
	log.Logf("Ensure that your public key (contained in %s/.ssh/id_rsa.pub or %s/.ssh/id_dsa.pub) is present in your Bytemark user's keys (see `bytemark show keys`, `bytemark add key`)", os.Getenv("HOME"), os.Getenv("HOME"))
	log.Log()
	log.Logf("Then connect to %s@%s\r\n", global.Client.GetSessionUser(), mgmtAddress)

}

func console_panel(vm *lib.VirtualMachine) error {
	ep := global.Config.EndpointName()
	token := global.Config.GetIgnoreErr("token")
	url := fmt.Sprintf("%s/vnc/?auth_token=%s&endpoint=%s&management_ip=%s", global.Config.PanelURL(), token, shortEndpoint(ep), vm.ManagementAddress)
	return util.CallBrowser(url)
}

func console_serial(vm *lib.VirtualMachine) error {
	host := fmt.Sprintf("%s@%s", global.Client.GetSessionUser(), vm.ManagementAddress)
	log.Logf("ssh %s\r\n", host)
	bin, err := exec.LookPath("ssh")
	if err != nil {
		log.Log("Unable to find an ssh executable")
		return err
	}
	err = syscall.Exec(bin, []string{"ssh", host}, os.Environ())
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			if errno != 0 {
				log.Log("Attempting to exec ssh failed. Please file a bug report.")
				return err
			}
		} else {
			log.Log("Couldn't connect to the management address. Please ensure you have an SSH client in your $PATH.")
			return err
		}
	}
	return nil
}
