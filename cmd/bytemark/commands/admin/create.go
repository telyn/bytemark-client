package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "create",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "ip range",
				Usage:     "create a new IP range in a VLAN",
				UsageText: "bytemark --admin create ip range <ip-range> <vlan-num>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "ip-range",
						Usage: "the IP range to add",
					},
					cli.IntFlag{
						Name:  "vlan-num",
						Usage: "The VLAN number to add the IP range to",
					},
				},
				Action: app.Action(args.Optional("ip-range", "vlan-num"), with.RequiredFlags("ip-range", "vlan-num"), with.Auth, func(c *app.Context) error {
					if err := c.Client().CreateIPRange(c.String("ip-range"), c.Int("vlan-num")); err != nil {
						return err
					}
					log.Logf("IP range created\r\n")
					return nil
				}),
			}, {
				Name:      "migration",
				Usage:     "creates a new migration job",
				UsageText: "bytemark --admin create migration [--priority <number>] [--disc <disc>]... [--pool <pool>]... [--tail <tail>]... [--to-pool <pool>]...",
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
						Usage: "a disc to migrate to some other destination. Can be specified multiple times to migrate multiple discs",
					},
					cli.IntFlag{
						Name:  "priority",
						Usage: "an optional priority to set - bigger is higher priority",
					},
				},
				Action: app.Action(with.Auth, func(ctx *app.Context) err {
					jobRequest := brain.MigrationJob{
						Sources: brain.MigrationJobLocations{
							Discs: ctx.StringSlice("disc"),
							Pools: ctx.StringSlice("pool"),
							Tails: ctx.StringSlice("tail"),
						},
						Destinations: brain.MigrationJobDestinations{
							Pools: ctx.StringSlice("to-pool"),
						},
						Priority: ctx.Int("priority"),
					}
					job, err := adminRequests.CreateMigrationJob(jobRequest)
					if err != nil {
						ctx.Output("Migration job created")
						return ctx.OutputInDesiredForm(job)
					}
					ctx.Output("Couldn't create migration job")
					return err
				}),
			}, {
				Name:        "user",
				Usage:       "creates a new cluster admin or cluster superuser",
				UsageText:   "bytemark --admin create user <username> <privilege>",
				Description: `creates a new cluster admin or superuser. The privilege field must be either cluster_admin or cluster_su.`,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "username",
						Usage: "The username of the new user",
					},
					cli.StringFlag{
						Name:  "privilege",
						Usage: "The privilege to grant to the new user",
					},
				},
				Action: app.Action(args.Optional("username", "privilege"), with.RequiredFlags("username", "privilege"), with.Auth, func(c *app.Context) error {
					// Privilege is just a string and not a app.PrivilegeFlag, since it can only be "cluster_admin" or "cluster_su"
					if err := c.Client().CreateUser(c.String("username"), c.String("privilege")); err != nil {
						return err
					}
					log.Logf("User %s has been created with %s privileges\r\n", c.String("username"), c.String("privilege"))
					return nil
				}),
			}, {
				Name:      "vlan group",
				Aliases:   []string{"vlan-group"},
				Usage:     "creates groups for private VLANs",
				UsageText: "bytemark --admin create vlan group <group> [vlan-num]",
				Description: `Create a group in the specified account, with an optional VLAN specified.

Used when setting up a private VLAN for a customer.`,
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "group",
						Usage: "the name of the group to create",
						Value: new(app.GroupNameFlag),
					},
					cli.IntFlag{
						Name:  "vlan-num",
						Usage: "The VLAN number to add the group to",
					},
				},
				Action: app.Action(args.Optional("group", "vlan-num"), with.RequiredFlags("group"), with.Auth, func(c *app.Context) error {
					gp := c.GroupName("group")
					if err := c.Client().AdminCreateGroup(gp, c.Int("vlan-num")); err != nil {
						return err
					}
					log.Logf("Group %s was created under account %s\r\n", gp.Group, gp.Account)
					return nil
				}),
			},
		},
	})
}
