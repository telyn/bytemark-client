package admin

import (
	"errors"
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	billingMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
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

	Commands = append(Commands, cli.Command{
		Name:   "update",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "billing-definition",
				Usage:     "update a bmbilling definition",
				UsageText: "--admin update billing-definition [flags] [name] [value]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name",
						Usage: "the name of the definition to set",
					},
					cli.StringFlag{
						Name:  "value",
						Usage: "the value of the definition to set",
					},
					cli.StringFlag{
						Name:  "group",
						Usage: "the group a user must be in to update the definition",
					},
				},
				Action: app.Action(args.Optional("name", "value"), with.RequiredFlags("name", "value"), with.Auth, func(ctx *app.Context) error {
					def := billing.Definition{
						Name:           ctx.String("name"),
						Value:          ctx.String("value"),
						UpdateGroupReq: ctx.String("group"),
					}
					if _, err := billingMethods.GetDefinition(ctx.Client(), def.Name); err != nil {
						if _, ok := err.(lib.NotFoundError); ok {
							ctx.LogErr("Couldn't find a definition called %s - aborting.", def.Name)
							return nil
						}
						return err
					}
					err := billingMethods.UpdateDefinition(ctx.Client(), def)
					if err == nil {
						ctx.LogErr("Updated %s to %s", def.Name, def.Value)
					}
					return err

				}),
			}, {
				Name:      "head",
				Usage:     "update the settings of a head",
				UsageText: "--admin update head <head> [--usage-strategy] [--overcommit-ratio] [--label]",
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
				Action: app.Action(args.Optional("head", "usage-strategy", "overcommit-ratio", "label"), with.RequiredFlags("head"), with.Auth, func(c *app.Context) error {
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
				UsageText: "--admin update tail <tail> [--usage-strategy] [--overcommit-ratio] [--label]",
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
				Action: app.Action(args.Optional("tail", "usage-strategy", "overcommit-ratio", "label"), with.RequiredFlags("tail"), with.Auth, func(c *app.Context) error {
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
				UsageText: "--admin update storage pool [--usage-strategy new-strategy] [--overcommit-ratio new-ratio] [--label new-label] [--migration-concurrency new-limit] <storage pool>",
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
					cli.IntFlag{
						Name:  "migration-concurrency",
						Usage: "the number of concurrent migrations the storage pool can handle",
					},
				},
				Action: app.Action(args.Optional("storage-pool", "usage-strategy", "overcommit-ratio"), with.RequiredFlags("storage-pool"), with.Auth, func(c *app.Context) error {

					options := brain.StoragePool{
						UsageStrategy:        c.String("usage-strategy"),
						OvercommitRatio:      c.Int("overcommit-ratio"),
						Label:                c.String("label"),
						MigrationConcurrency: c.Int("migration-concurrency"),
					}

					if err := c.Client().UpdateStoragePool(c.String("storage-pool"), options); err != nil {
						return err
					}

					log.Outputf("Storage pool %s updated\n", c.String("storage-pool"))

					return nil
				}),
			},
			{
				Name:        "server-migration",
				Usage:       "update the settings of an in-progress migration",
				UsageText:   "--admin update server-migration <name> [--migrate-speed] [--migrate-downtime]",
				Description: `This command migrates a server to a new head. If a new head isn't supplied, a new one is picked automatically.`,
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server to migrate",
						Value: new(flags.VirtualMachineName),
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
				Action: app.Action(args.Optional("server", "migrate-speed", "migrate-downtime"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) error {
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
			{
				Name:        "migration",
				Usage:       "update a migration job",
				UsageText:   "--admin update migration --id <id> --priority <priority> --cancel-disc <disc> --cancel-pool <pool> --cancel-tail <tail> | --cancel-all",
				Description: `This command allows you to update an ongoing migration job by altering its priority, cancelling migrating discs, pools, tails, or canceling everything for the current job`,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "id",
						Usage: "the id of the migration job",
					},
					cli.IntFlag{
						Name:  "priority",
						Usage: "the priority of the current job",
					},
					cli.StringSliceFlag{
						Name:  "cancel-disc",
						Usage: "the disc(s) to cancel migration of",
					},
					cli.StringSliceFlag{
						Name:  "cancel-pool",
						Usage: "the pool(s) to cancel migration of",
					},
					cli.StringSliceFlag{
						Name:  "cancel-tail",
						Usage: "the tail(s) to cancel migration of",
					},
					cli.BoolFlag{
						Name:  "cancel-all",
						Usage: "cancel the all migrations of the job",
					},
				},
				Action: app.Action(with.RequiredFlags("id"), with.Auth, func(c *app.Context) error {
					id := c.Context.Int("id")
					discs := c.Context.StringSlice("cancel-disc")
					pools := c.Context.StringSlice("cancel-pool")
					tails := c.Context.StringSlice("cancel-tail")

					allCancelled := append(discs, pools...)
					allCancelled = append(allCancelled, tails...)

					if len(allCancelled) == 0 && !c.Context.IsSet("priority") && !c.Context.IsSet("cancel-all") {
						return fmt.Errorf("No Flags have been set. Please specify a priority, ")
					}

					if c.Context.IsSet("cancel-all") {
						if len(allCancelled) > 0 || c.Context.IsSet("priority") {
							return fmt.Errorf("You have set additional flags as well as --cancel-all. Nothing else can be specified when --cancel-all has been set")
						}

						err := brainMethods.CancelMigrationJob(c.Client(), id)
						if err != nil {
							return err
						}
						c.LogErr("All migrations for job %d have been cancelled.", id)
						return err
					}

					modifications := brain.MigrationJobModification{
						Cancel: brain.MigrationJobLocations{
							Discs: stringsToNumberOrStrings(discs),
							Pools: stringsToNumberOrStrings(pools),
							Tails: stringsToNumberOrStrings(tails),
						},
						Options: brain.MigrationJobOptions{
							Priority: c.Context.Int("priority"),
						},
					}

					err := brainMethods.EditMigrationJob(c.Client(), id, modifications)
					if err != nil {
						return err
					}

					if c.Context.IsSet("priority") {
						c.Log("Priority updated for job %d", id)
					}

					for _, cancelled := range allCancelled {
						c.Log("Migration cancelled for %s on job %d", cancelled, id)
					}

					return err
				}),
			},
		},
	})
}
