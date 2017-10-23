package main

import (
	"errors"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	readUpdateFlags := func(c *app.Context) (usageStrategy *string, overcommitRatio *int, label *string) {
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
				Name:      "bmbilling",
				Usage:     "update bmbilling's definitions",
				UsageText: "bytemark --admin update bmbilling [--trial-days <days>] [--trial-pence <pence>]",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "trial-days",
						Usage: "the number of days in future trials",
					},
					cli.IntFlag{
						Name:  "trial-pence",
						Usage: "the maximum monthly cost, in pence, of future trials",
					},
				},
				Action: app.Action(with.Auth, func(ctx *app.Context) error {
					billingDefinitions := billing.Definitions{
						TrialDays:  ctx.Int("trial-days"),
						TrialPence: ctx.Int("trial-pence"),
					}
					return billingRequests.UpdateDefinitions(ctx.Client(), billingDefinitions)
				}),
			}, {
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
				Action: app.With(args.Optional("head", "usage-strategy", "overcommit-ratio", "label"), with.RequiredFlags("head"), with.Auth, func(c *app.Context) error {
					usageStrategy, overcommitRatio, label := readUpdateFlags(c)

					options := lib.UpdateHead{
						UsageStrategy:   usageStrategy,
						OvercommitRatio: overcommitRatio,
						Label:           label,
					}

					if err := c.Client().UpdateHead(c.String("head"), options); err != nil {
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
				Action: app.With(args.Optional("tail", "usage-strategy", "overcommit-ratio", "label"), with.RequiredFlags("tail"), with.Auth, func(c *app.Context) error {
					usageStrategy, overcommitRatio, label := readUpdateFlags(c)

					options := lib.UpdateTail{
						UsageStrategy:   usageStrategy,
						OvercommitRatio: overcommitRatio,
						Label:           label,
					}

					if err := c.Client().UpdateTail(c.String("tail"), options); err != nil {
						return err
					}

					log.Outputf("Tail %s updated\n", c.String("tail"))

					return nil
				}),
			},
			{
				Name:      "storage pool",
				Usage:     "update the settings of a storage pool",
				UsageText: "bytemark --admin update storage pool <storage-pool> [--usage-strategy] [--overcommit-ratio] [--label]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "storage-pool",
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
				Action: app.With(args.Optional("storage-pool", "usage-strategy", "overcommit-ratio", "label"), with.RequiredFlags("storage-pool"), with.Auth, func(c *app.Context) error {
					usageStrategy, overcommitRatio, label := readUpdateFlags(c)

					options := lib.UpdateStoragePool{
						UsageStrategy:   usageStrategy,
						OvercommitRatio: overcommitRatio,
						Label:           label,
					}

					if err := c.Client().UpdateStoragePool(c.String("storage-pool"), options); err != nil {
						return err
					}

					log.Outputf("Storage pool %s updated\n", c.String("storage-pool"))

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
								Value: new(app.VirtualMachineNameFlag),
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
						Action: app.With(args.Optional("server", "migrate-speed", "migrate-downtime"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) error {
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

							if err := c.Client().UpdateVMMigration(vm, speed, downtime); err != nil {
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
