package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "stop",
		Usage:       "stop a server, as though pulling the power cable out",
		UsageText:   "stop server <server>",
		Description: "This command will instantly power down a server. Note that this may cause data loss, particularly on servers with unjournaled file systems (e.g. ext2)",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "server",
			Usage:       "stop a server, as though pulling the power cable out",
			UsageText:   "stop server <server>",
			Description: "This command will instantly power down a server. Note that this may cause data loss, particularly on servers with unjournaled file systems (e.g. ext2)",

			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to stop",
					Value: new(flags.VirtualMachineNameFlag),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
				vmName := flags.VirtualMachineName(c, "server")
				c.Log("Attempting to stop %s...", vmName)
				err = c.Client().StopVirtualMachine(vmName)
				if err != nil {
					return
				}

				c.Log("%s stopped successfully.", vmName)
				return
			}),
		}},
	})
}
