package main

import (
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
						Action: With(OptionalArgs("disc"), RequiredFlags("disc"), AuthProvider, func(c *Context) error {
							if err := global.Client.CancelDiscMigration(c.Int("disc")); err != nil {
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
						Action: With(OptionalArgs("server"), RequiredFlags("server"), AuthProvider, func(c *Context) error {
							if err := global.Client.CancelVMMigration(c.Int("server")); err != nil {
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
