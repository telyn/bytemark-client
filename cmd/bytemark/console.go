package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
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
			cli.GenericFlag{
				Name:  "server",
				Usage: "The server whose console will be connected to",
				Value: new(app.VirtualMachineNameFlag),
			},
		},
		Action: app.With(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(ctx *app.Context) (err error) {
			vmName := ctx.VirtualMachineName("server")
			if ctx.Bool("serial") && ctx.Bool("panel") {
				return ctx.Help("You must only specify one of --serial and --panel!")
			}

			vm, err := ctx.Client().GetVirtualMachine(vmName)
			if err != nil {
				return
			}
			if ctx.Bool("no-connect") {
				serialConsoleInstructions(ctx, vm)
				log.Log()
				vncConsoleInstructions(ctx, vm)
				return nil
			}
			if ctx.Bool("panel") {
				err = openPanelConsole(ctx, vm)
			} else {
				err = connectSerialConsole(ctx, vm)
			}
			return

		}),
	})
}

func shortEndpoint(endpoint string) string {
	return strings.Split(endpoint, ".")[0]
}

func vncConsoleInstructions(c *app.Context, vm brain.VirtualMachine) {
	mgmtAddress := vm.ManagementAddress.String()
	if vm.ManagementAddress.To4() == nil {
		mgmtAddress = "[" + mgmtAddress + "]"
	}

	log.Log("VNC graphical console connection information for", vm.Hostname)
	log.Log()
	log.Logf("Ensure that your public key (contained in %s/.ssh/id_rsa.pub or %s/.ssh/id_dsa.pub) is present in your Bytemark user's keys (see `bytemark show keys`, `bytemark add key`)", os.Getenv("HOME"), os.Getenv("HOME"))
	log.Log()
	log.Logf("Then set up a tunnel using SSH: ssh -L 9999:%s:5900 %s@%s\r\n", mgmtAddress, c.Client().GetSessionUser(), mgmtAddress)
	log.Log()
	log.Log("You will then be able to connect to vnc://localhost:9999/")
	log.Log("Any port may be substituted for 9999 as long as the same port is used in both commands")
}

func serialConsoleInstructions(c *app.Context, vm brain.VirtualMachine) {
	mgmtAddress := vm.ManagementAddress.String()
	if vm.ManagementAddress.To4() == nil {
		mgmtAddress = "[" + mgmtAddress + "]"
	}
	log.Log("Serial console connection information for", vm.Hostname)
	log.Log()
	log.Logf("Ensure that your public key (contained in %s/.ssh/id_rsa.pub or %s/.ssh/id_dsa.pub) is present in your Bytemark user's keys (see `bytemark show keys`, `bytemark add key`)", os.Getenv("HOME"), os.Getenv("HOME"))
	log.Log()
	log.Logf("Then connect to %s@%s\r\n", c.Client().GetSessionUser(), mgmtAddress)

}

func openPanelConsole(c *app.Context, vm brain.VirtualMachine) error {
	ep := c.Config().EndpointName()
	token := c.Config().GetIgnoreErr("token")
	url := fmt.Sprintf("%s/vnc/?auth_token=%s&endpoint=%s&management_ip=%s", c.Config().PanelURL(), token, shortEndpoint(ep), vm.ManagementAddress)
	log.Logf("Opening %s in a browser.\r\n", url)
	return util.CallBrowser(url)
}

func collectArgs(args string) (slice []string) {
	inDbl := false
	inSgl := false

	var cur []rune
	if args == "" {
		return
	}
	for _, ch := range args {
		if inDbl && ch == '"' {
			inDbl = false
			continue
		} else if inSgl && ch == '\'' {
			inSgl = false
			continue
		}

		if !inSgl && ch == '"' {
			inDbl = true
			continue
		} else if !inDbl && ch == '\'' {
			inSgl = true
			continue
		}

		if !inDbl && !inSgl && ch == ' ' {
			slice = append(slice, string(cur))
			cur = make([]rune, 0)
			continue
		}
		cur = append(cur, ch)

	}
	slice = append(slice, string(cur))
	return
}

func connectSerialConsole(c *app.Context, vm brain.VirtualMachine) error {
	sshargs := c.String("ssh-args")
	host := fmt.Sprintf("%s@%s", c.Client().GetSessionUser(), vm.ManagementAddress)
	log.Logf("ssh %s\r\n", host)

	bin, err := exec.LookPath("ssh")
	if err != nil {
		return err
	}

	sshargsli := collectArgs(sshargs)
	args := make([]string, len(sshargsli)+2)
	copy(args[1:], sshargsli[0:])
	args[0] = "ssh"
	args[len(args)-1] = host

	log.Debugf(5, "%+v\r\n", args)

	/* #nosec */
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
