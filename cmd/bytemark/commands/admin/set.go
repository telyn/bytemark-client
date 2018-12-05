package admin

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {

	Commands = append(Commands, cli.Command{
		Name:   "set",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "iops limit",
				Usage:     "set the IOPS limit of a disc",
				UsageText: "--admin set disc iops limit <server> <disc> <limit>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "disc",
						Usage: "the name of the disc to alter the iops limit of",
					},
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server the disc belongs to",
						Value: new(flags.VirtualMachineName),
					},
					cli.IntFlag{
						Name:  "iops-limit",
						Usage: "the limit to set",
					},
				},
				Action: app.Action(args.Optional("server", "disc", "iops-limit"), with.RequiredFlags("server", "disc", "iops-limit"), with.Auth, func(c *app.Context) error {
					iopsLimit := c.Int("iops-limit")
					if iopsLimit < 1 {
						return fmt.Errorf("IOPS limit must be at least 1")
					}
					vmName := c.VirtualMachineName("server")

					return c.Client().SetDiscIopsLimit(vmName, c.String("disc"), iopsLimit)
				}),
			},
		},
	})
}
