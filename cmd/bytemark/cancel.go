package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "cancel",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:   "migration",
				Action: cli.ShowSubcommandHelp,
				Subcommands: []cli.Command{
					{
						Name:      "disc",
						Usage:     "cancel a disc migration",
						UsageText: "bytemark --admin cancel migration disc <disc>",
						Flags: []cli.Flag{
							cli.IntFlag{
								Name:  "disc",
								Usage: "the ID of the disc that is migrating",
							},
						},
						Action: app.With(args.Optional("disc"), with.RequiredFlags("disc"), with.Auth, func(c *app.Context) error {
							if err := c.Client().CancelDiscMigration(c.Int("disc")); err != nil {
								return err
							}

							log.Output("Migration cancelled")

							return nil
						}),
					},
					{
						Name:      "server",
						Aliases:   []string{"vm"},
						Usage:     "cancel a server migration",
						UsageText: "bytemark --admin cancel migration server <disc>",
						Flags: []cli.Flag{
							cli.IntFlag{
								Name:  "server",
								Usage: "the ID of the server that is migrating",
							},
						},
						Action: app.With(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) error {
							if err := c.Client().CancelVMMigration(c.Int("server")); err != nil {
								return err
							}

							log.Output("Migration cancelled")

							return nil
						}),
					},
				},
			},
		},
	})
}
