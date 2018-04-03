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
		Name:   "reject",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "server",
				Aliases:   []string{"vm"},
				Usage:     "reject a server, and specify the reason for the rejection",
				UsageText: "--admin reject server <name> <reason>",
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "The server to reject",
						Value: new(app.VirtualMachineNameFlag),
					},
					cli.StringFlag{
						Name:  "reason",
						Usage: "The reason why the server is being rejected.",
					},
				},
				Action: app.Action(args.Optional("server"), args.Join("reason"), with.RequiredFlags("server", "reason"), with.Auth, func(c *app.Context) error {
					vm := c.VirtualMachineName("server")

					if err := c.Client().RejectVM(vm, c.String("reason")); err != nil {
						return err
					}

					log.Outputf("Server %s was successfully rejected\n", vm.String())

					return nil
				}),
			},
		},
	})
}
