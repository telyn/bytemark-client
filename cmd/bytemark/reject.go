package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "reject",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "server",
				Aliases:   []string{"vm"},
				Usage:     "reject a server, and specify the reason for the rejection.",
				UsageText: "bytemark --admin reject server <name> <reason>",
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "The server to reject",
						Value: new(VirtualMachineNameFlag),
					},
					cli.StringFlag{
						Name:  "reason",
						Usage: "The reason why the server is being rejected.",
					},
				},
				Action: With(OptionalArgs("server"), JoinArgs("reason"), RequiredFlags("server", "reason"), AuthProvider, func(c *Context) error {
					vm := c.VirtualMachineName("server")

					if err := global.Client.RejectVM(&vm, c.String("reason")); err != nil {
						return err
					}

					log.Outputf("Server %s was successfully rejected\n", vm.String())

					return nil
				}),
			},
		},
	})
}
