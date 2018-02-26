package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "approve",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "server",
				Aliases:   []string{"vm"},
				Usage:     "approve a server, and optionally power it on",
				UsageText: "bytemark --admin approve server <name> [--power-on]",
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "The server to approve",
						Value: new(app.VirtualMachineNameFlag),
					},
					cli.BoolFlag{
						Name:  "power-on",
						Usage: "If set, powers on the server",
					},
				},
				Action: app.Action(args.Optional("server", "power-on"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) error {
					vm := c.VirtualMachineName("server")

					if err := c.Client().ApproveVM(vm, c.Bool("power-on")); err != nil {
						return err
					}

					log.Outputf("Server %s was successfully approved\n", vm.String())

					return nil
				}),
			},
		},
	})
}
