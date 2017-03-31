package main

import (
	"errors"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	readUpdateFlags := func(c *Context) (usageStrategy *string, overcommitRatio *int, label *string) {
		if c.Context.IsSet("usage-strategy") {
			v := c.String("usage-strategy")
			usageStrategy = &v
		}

		if c.Context.IsSet("overcommit-ratio") {
			v := c.Int("overcommit-ratio")
			overcommitRatio = &v
		}

		if c.Context.IsSet("label") {
			v := c.String("label")
			label = &v
		}

		return
	}

	adminCommands = append(adminCommands, cli.Command{
		Name:   "update",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "head",
				Usage:     "update the settings of a head",
				UsageText: "bytemark --admin update head <head> [--usage-strategy] [--overcommit-ratio] [--label]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "head",
						Usage: "the ID or label of the head to be updated",
					},
					cli.StringFlag{
						Name:  "usage-strategy",
						Usage: "the usage strategy of the head",
					},
					cli.IntFlag{
						Name:  "overcommit-ratio",
						Usage: "the overcommit ratio of the head",
					},
					cli.StringFlag{
						Name:  "label",
						Usage: "the label of the head",
					},
				},
				Action: With(OptionalArgs("head", "usage-strategy", "overcommit-ratio", "label"), RequiredFlags("head"), AuthProvider, func(c *Context) error {
					usageStrategy, overcommitRatio, label := readUpdateFlags(c)

					options := &lib.UpdateHead{
						UsageStrategy:   usageStrategy,
						OvercommitRatio: overcommitRatio,
						Label:           label,
					}

					if err := global.Client.UpdateHead(c.String("head"), options); err != nil {
						return err
					}

					log.Outputf("Head %s updated\n", c.String("head"))

					return nil
				}),
			},
			{
				Name:      "tail",
				Usage:     "update the settings of a tail",
				UsageText: "bytemark --admin update tail <tail> [--usage-strategy] [--overcommit-ratio] [--label]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "tail",
						Usage: "the ID or label of the tail to be updated",
					},
					cli.StringFlag{
						Name:  "usage-strategy",
						Usage: "the usage strategy of the tail",
					},
					cli.IntFlag{
						Name:  "overcommit-ratio",
						Usage: "the overcommit ratio of the tail",
					},
					cli.StringFlag{
						Name:  "label",
						Usage: "the label of the tail",
					},
				},
				Action: With(OptionalArgs("tail", "usage-strategy", "overcommit-ratio", "label"), RequiredFlags("tail"), AuthProvider, func(c *Context) error {
					usageStrategy, overcommitRatio, label := readUpdateFlags(c)

					options := &lib.UpdateTail{
						UsageStrategy:   usageStrategy,
						OvercommitRatio: overcommitRatio,
						Label:           label,
					}

					if err := global.Client.UpdateTail(c.String("tail"), options); err != nil {
						return err
					}

					log.Outputf("Tail %s updated\n", c.String("tail"))

					return nil
				}),
			},
			{
				Name:      "storage_pool",
				Usage:     "update the settings of a storage pool",
				UsageText: "bytemark --admin update storage_pool <storage_pool> [--usage-strategy] [--overcommit-ratio] [--label]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "storage_pool",
						Usage: "the ID or label of the storage pool to be updated",
					},
					cli.StringFlag{
						Name:  "usage-strategy",
						Usage: "the usage strategy of the storage pool",
					},
					cli.IntFlag{
						Name:  "overcommit-ratio",
						Usage: "the overcommit ratio of the storage pool",
					},
					cli.StringFlag{
						Name:  "label",
						Usage: "the label of the storage pool",
					},
				},
				Action: With(OptionalArgs("storage_pool", "usage-strategy", "overcommit-ratio", "label"), RequiredFlags("storage_pool"), AuthProvider, func(c *Context) error {
					usageStrategy, overcommitRatio, label := readUpdateFlags(c)

					options := &lib.UpdateStoragePool{
						UsageStrategy:   usageStrategy,
						OvercommitRatio: overcommitRatio,
						Label:           label,
					}

					if err := global.Client.UpdateStoragePool(c.String("storage_pool"), options); err != nil {
						return err
					}

					log.Outputf("Storage pool %s updated\n", c.String("storage_pool"))

					return nil
				}),
			},
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
