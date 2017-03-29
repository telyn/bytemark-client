package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "approve",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "server",
				Aliases:   []string{"vm"},
				Usage:     "approve a server, and optionally power it on",
				UsageText: "bytemark --admin approve server <name> [power_on]",
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "The server to approve",
						Value: new(VirtualMachineNameFlag),
					},
					cli.BoolFlag{
						Name:  "power_on",
						Usage: "If set, powers on the server.",
					},
				},
				Action: With(OptionalArgs("server", "power_on"), RequiredFlags("server"), AuthProvider, func(c *Context) (err error) {
					vm := c.VirtualMachineName("server")

					if err := global.Client.ApproveVM(&vm, c.Bool("power_on")); err != nil {
						return err
					}

					log.Outputf("Server %s was successfully approved\n", vm.String())

					return nil
				}),
			},
		},
	})
}
