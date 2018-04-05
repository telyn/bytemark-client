package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "create",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "migration",
				Usage:     "creates a new migration job",
				UsageText: "--admin create migration [--priority <number>] [--disc <disc>]... [--pool <pool>]... [--tail <tail>]... [--to-pool <pool>]...",
				Description: `Requests that the brain starts a migration from the various discs, pools and tails specified to the various pools specified.
				
   If no --to-pools are specified, the brain will decide destinations on its own. To migrate multiple discs, pools, or tails, or to distribute migrating discs across multiple storage pools, simply specify those flags multiple times.
 
   Source discs, tails and pools and destination pools can all be specified as their numeric IDs, their label, or their IP address.
When there are multiple migration jobs at once and migration concurrency limits are being reached, jobs with highest priority will have their migrations scheduled ahead of other jobs.

EXAMPLES:
   # migrate one disc whose IP is fe80::1 to tail2-sata1
   bytemark --admin create migration --disc fe80::1 --to-pool tail2-sata1

   # empty a tail to wherever is appropriate with a high priority
   bytemark --admin create migration --priority 10 --tail tail1

   # migrate two discs, one called very-unique-label and one whose ID is 45859 to tail2-sata1
   bytemark --admin create migration --disc very-unique-label --disc 45859 --to-pool tail2-sata1

   # migrate one pool to another
   bytemark --admin create migration --pool tail1-sata1 --to-pool tail2-sata1

   # migrate two pools to one pool
   bytemark --admin create migration --pool tail1-sata1 --pool tail1-sata2 --to-pool tail2-sata1

   # migrate a pool and a specific disc away from their current tail
   bytemark --admin create migration --pool tail1-sata1 --disc 45859`,

				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "disc",
						Usage: "a disc to migrate to some other destination. Can be specified multiple times to migrate multiple discs",
					},
					cli.StringSliceFlag{
						Name:  "pool",
						Usage: "a storage pool to migrate to some other destination. Can be specified multiple times to migrate multiple storage pools",
					},
					cli.StringSliceFlag{
						Name:  "tail",
						Usage: "a tail to migrate to some other destination. Can be specified multiple times to migrate multiple tails",
					},
					cli.StringSliceFlag{
						Name:  "to-pool",
						Usage: "a target pool for migrations. Can be specified multiple times to migrate multiple discs",
					},
					cli.StringSliceFlag{
						Name:  "to-tail",
						Usage: "a target tail for migrations. Can be specified multiple times to migrate multiple discs",
					},
					cli.IntFlag{
						Name:  "priority",
						Usage: "an optional priority to set - bigger is higher priority",
					},
				},
				Action: app.Action(with.Auth, func(ctx *app.Context) error {
					jobRequest := brain.MigrationJobSpec{
						Sources: brain.MigrationJobLocations{
							Discs: StringsToNumberOrStrings(ctx.StringSlice("disc")),
							Pools: StringsToNumberOrStrings(ctx.StringSlice("pool")),
							Tails: StringsToNumberOrStrings(ctx.StringSlice("tail")),
						},
						Destinations: brain.MigrationJobLocations{
							Pools: StringsToNumberOrStrings(ctx.StringSlice("to-pool")),
							Tails: StringsToNumberOrStrings(ctx.StringSlice("to-tail")),
						},
						Options: brain.MigrationJobOptions{
							Priority: ctx.Int("priority"),
						},
					}
					job, err := brainRequests.CreateMigrationJob(ctx.Client(), jobRequest)
					if err == nil {
						return ctx.OutputInDesiredForm(job)
					}
					return err
				}),
			},
		},
	})
}
