package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "console",
		Usage:     "connect to a server's serial or graphical console - as though physically plugging in",
		UsageText: "bytemark console [--serial | --vnc | --panel] [--no-connect] <server>",
		Description: `Out-of-band access to a server's serial or graphical (VNC) console.
Under systems with no GUI, sometimes errors will output to the graphical console and not the serial console.
Defaults to connecting to the serial console for the given server.`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "serial",
				Usage: "Connect to the serial console. Cannot be set at same time as  --vnc or --panel. You must have an ssh client on your computer.",
			},
			cli.BoolFlag{
				Name:  "vnc",
				Usage: "Connect to the graphical console. Cannot be set at the same time as --serial or --panel. You must have a openssh client and a VNC client set up on your computer.",
			},
			cli.BoolFlag{
				Name:  "panel",
				Usage: "Connect to the graphical console via the Bytemark panel. Cannot be set at the same time as --serial or --vnc. You must have a graphical browser installed on your computer.",
			},
			cli.BoolFlag{
				Name:  "no-connect",
				Usage: "Output connection instructions, rather than directly connecting.",
			},
			cli.StringFlag{
				Name:  "ssh-args",
				Usage: "Arguments that will be passed to SSH (only applies to --serial).",
			},
		},
		Action: With(VirtualMachineNameProvider, AuthProvider, func(ctx *Context) error {
			if ctx.Context.Bool("serial") && ctx.Context.Bool("panel") {
				return ctx.Help("You must only specify one of --serial and --panel!")
			}

			vm, err := global.Client.GetVirtualMachine(ctx.VirtualMachineName)
			if err != nil {
				return err
			}
			if ctx.Context.Bool("no_connect") {
				console_serial_instructions(vm)
				log.Log()
				console_vnc_instructions(vm)
				return nil
			} else {
				if ctx.Context.Bool("panel") {
					return console_panel(vm)
				} else {
					return console_serial(vm, ctx.String("ssh-args"))
				}
			}

		}),
	})
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
	log.Logf("Opening %s in a browser.\r\n", url)
	return util.CallBrowser(url)
}

func collect_args(args string) (slice []string) {
	in_dbl := false
	in_sgl := false
	slice = make([]string, 0)
	cur := make([]rune, 0)
	for _, ch := range args {
		if in_dbl && ch == '"' {
			in_dbl = false
			continue
		} else if in_sgl && ch == '\'' {
			in_sgl = false
			continue
		}

		if !in_sgl && ch == '"' {
			in_dbl = true
			continue
		} else if !in_dbl && ch == '\'' {
			in_sgl = true
			continue
		}

		if !in_dbl && !in_sgl && ch == ' ' {
			slice = append(slice, string(cur))
			cur = make([]rune, 0)
			continue
		}
		cur = append(cur, ch)

	}
	slice = append(slice, string(cur))
	return
}

func console_serial(vm *lib.VirtualMachine, sshargs string) error {
	host := fmt.Sprintf("%s@%s", global.Client.GetSessionUser(), vm.ManagementAddress)
	log.Logf("ssh %s\r\n", host)

	bin, err := exec.LookPath("ssh")
	if err != nil {
		return err
	}

	sshargsli := collect_args(sshargs)
	args := make([]string, len(sshargsli)+2)
	copy(args[1:], sshargsli[0:])
	args[0] = "ssh"
	args[len(args)-1] = host

	log.Debugf(5, "%+v\r\n", args)

	err = syscall.Exec(bin, args, os.Environ())
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			if errno != 0 {
				log.Log("Attempting to exec ssh failed. Please file a bug report.")
				return err
			}
		} else {
			log.Log("Couldn't connect to the management address.")
			return err
		}
	}
	return nil
}
