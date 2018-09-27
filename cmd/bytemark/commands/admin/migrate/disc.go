package migrate

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "disc",
		Usage:       "migrate a disc to a new storage pool",
		UsageText:   "--admin migrate disc <disc> [new-storage-pool]",
		Description: `This command migrates a disc to a new storage pool. If a new storage pool isn't supplied, a new one is picked automatically.`,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "disc",
				Usage: "the ID of the disc to migrate",
			},
			cli.StringFlag{
				Name:  "new-storage-pool",
				Usage: "the storage pool to move the disc to",
			},
		},
		Action: app.Action(args.Optional("disc", "new-storage-pool"), with.RequiredFlags("disc"), with.Auth, func(c *app.Context) (err error) {
			disc := c.Int("disc")
			storagePool := c.String("new-storage-pool")

			if err := c.Client().MigrateDisc(disc, storagePool); err != nil {
				return err
			}

			log.Outputf("Migration for disc %d initiated\n", disc)

			return nil
		}),
	})
}
