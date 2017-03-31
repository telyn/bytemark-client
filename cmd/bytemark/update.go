package main

import (
	"errors"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "update",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:    "server",
				Aliases: []string{"vm"},
				Action:  cli.ShowSubcommandHelp,
				Subcommands: []cli.Command{
					{
						Name:        "migration",
						Usage:       "update the settings of an in-progress migration",
						UsageText:   "bytemark --admin update server migration <name> [--migrate-speed] [--migrate-downtime]",
						Description: `This command migrates a server to a new head. If a new head isn't supplied, a new one is picked automatically.`,
						Flags: []cli.Flag{
							cli.GenericFlag{
								Name:  "server",
								Usage: "the server to migrate",
								Value: new(VirtualMachineNameFlag),
							},
							cli.Int64Flag{
								Name:  "migrate-speed",
								Usage: "the max speed to migrate the server at",
							},
							cli.IntFlag{
								Name:  "migrate-downtime",
								Usage: "the max allowed downtime",
							},
						},
						Action: With(OptionalArgs("server", "migrate-speed", "migrate-downtime"), RequiredFlags("server"), AuthProvider, func(c *Context) error {
							vm := c.VirtualMachineName("server")

							var speed *int64
							var downtime *int

							if c.Context.IsSet("migrate-speed") {
								s := c.Int64("migrate-speed")
								speed = &s
							}
							if c.Context.IsSet("migrate-downtime") {
								d := c.Int("migrate-downtime")
								downtime = &d
							}

							if speed == nil && downtime == nil {
								return errors.New("Nothing to update")
							}

							if err := global.Client.UpdateVMMigration(&vm, speed, downtime); err != nil {
								return err
							}

							log.Outputf("Migration for server %s updated\n", vm.String())

							return nil
						}),
					},
				},
			},
		},
	})
}
