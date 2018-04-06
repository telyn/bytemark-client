package main

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/urfave/cli"
)

func init() {

	commands = append(commands, cli.Command{
		Name:      "set",
		Usage:     "change hardware properties of Bytemark servers",
		UsageText: "set cores <server>",
		Description: `change hardware properties of Bytemark servers
		
These commands set various hardware properties of Bytemark servers. Note that for cores to take effect you will need to restart the server.`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "cdrom",
				Usage:     "attach a cdrom to your Bytemark Cloud Server",
				UsageText: "attach cdrom <server> <cdurl>",
				Description: `attach a cdrom to your Bytemark Cloud Server

This command allows you to add a cdrom to your Bytemark server. The CD must be publicly available over HTTP in order to be attached.`,
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server to attach the CD to",
						Value: new(app.VirtualMachineNameFlag),
					},
					cli.StringFlag{
						Name:  "cd-url",
						Usage: "an HTTP(S) URL for an ISO image file to attach. If not set or set to the empty string, will 'eject' the current CD",
					},
				},
				Action: app.Action(args.Optional("server", "cd-url"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) error {
					vmName := c.VirtualMachineName("server")
					err := c.Client().SetVirtualMachineCDROM(vmName, c.String("cd-url"))
					if _, ok := err.(lib.InternalServerError); ok {
						return c.Help("Couldn't set the server's cdrom - check that you have provided a valid public HTTP url")
					}
					return err
				}),
			},
			{
				Name:        "cores",
				Usage:       "set the number of CPU cores on a Bytemark cloud server",
				UsageText:   "set cores <server name> <cores>",
				Description: "This command sets the number of CPU cores used by the cloud server. This will usually require a restart of the server to take effect.",
				Flags: []cli.Flag{
					forceFlag,
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server to alter",
						Value: new(app.VirtualMachineNameFlag),
					},
					cli.IntFlag{
						Name:  "cores",
						Usage: "the number of cores that should be available to the VM",
					},
				},
				Action: app.Action(args.Optional("server", "cores"), with.RequiredFlags("server", "cores"), with.VirtualMachine("server"), func(c *app.Context) error {
					// cores should be a flag
					vmName := c.VirtualMachineName("server")
					cores := c.Int("cores")

					if c.VirtualMachine.Cores < cores {
						if !c.Bool("force") && !util.PromptYesNo(c.Prompter(), fmt.Sprintf("You are increasing the number of cores from %d to %d. This may cause your VM to cost more, are you sure?", c.VirtualMachine.Cores, cores)) {
							return util.UserRequestedExit{}
						}
					}
					return c.Client().SetVirtualMachineCores(vmName, cores)
				}),
			},
		},
	})
}
