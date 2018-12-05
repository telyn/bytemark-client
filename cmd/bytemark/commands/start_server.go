package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "start",
		Usage:       "start a stopped server",
		UsageText:   "start server <server>",
		Description: "This command will start a server that is not currently running.",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "server",
			Usage:       "start a stopped server",
			UsageText:   "start server <server>",
			Description: "This command will start a server that is not currently running.",
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to start",
					Value: new(flags.VirtualMachineName),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
				vmName := c.VirtualMachineName("server")
				c.Log("Attempting to start %s...", vmName)
				err = c.Client().StartVirtualMachine(vmName)
				if err != nil {
					return
				}

				c.Log("%s started successfully.", vmName)
				return
			}),
		}},
	})
}
