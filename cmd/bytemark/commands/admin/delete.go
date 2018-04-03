package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "delete",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "vlan",
				Usage:     "delete a given VLAN",
				UsageText: "--admin delete vlan <id>",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "id",
						Usage: "the ID of the VLAN to delete",
					},
				},
				Action: app.Action(args.Optional("id"), with.RequiredFlags("id"), with.Auth, func(c *app.Context) error {
					if err := c.Client().DeleteVLAN(c.Int("id")); err != nil {
						return err
					}

					c.Log("VLAN deleted")

					return nil
				}),
			},
		},
	})
}
