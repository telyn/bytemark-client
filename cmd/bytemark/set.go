package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/urfave/cli"
)

func init() {

	commands = append(commands, cli.Command{
		Name:      "set",
		Usage:     "change hardware properties of Bytemark servers",
		UsageText: "set cdrom <server>",
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
		},
	})
}
