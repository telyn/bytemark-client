package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/wait"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "shutdown",
		Usage:       "cleanly shut down a server",
		UsageText:   "shutdown server <server>",
		Description: "This command sends the ACPI shutdown signal to the server, causing a clean shut down. This is like pressing the power button on a computer you have physical access to.",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "server",
			Usage:       "cleanly shut down a server",
			UsageText:   "shutdown server <server>",
			Description: "This command sends the ACPI shutdown signal to the server, causing a clean shut down. This is like pressing the power button on a computer you have physical access to.",
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to shutdown",
					Value: new(flags.VirtualMachineNameFlag),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
				vmName := flags.VirtualMachineName(c, "server")
				c.Log("Shutting down %v...", vmName)
				err = c.Client().ShutdownVirtualMachine(vmName, true)
				if err != nil {
					return
				}

				err = wait.VMPowerOff(c, vmName)
				if err != nil {
					return
				}

				c.Log("Done!", vmName)
				return
			}),
		}},
	})
}
