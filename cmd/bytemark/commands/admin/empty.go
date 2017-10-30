package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "empty",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "storage pool",
				Usage:     "empty a storage pool",
				UsageText: "bytemark --admin empty storage pool <storage-pool>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "storage-pool",
						Usage: "the ID or label of the storage pool to be emptied",
					},
				},
				Action: app.Action(args.Optional("storage-pool"), with.RequiredFlags("storage-pool"), with.Auth, func(c *app.Context) error {
					if err := c.Client().EmptyStoragePool(c.String("storage-pool")); err != nil {
						return err
					}

					log.Output("Storage pool updated")

					return nil
				}),
			},
			{
				Name:      "head",
				Usage:     "empty a head",
				UsageText: "bytemark --admin empty head <head>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "head",
						Usage: "the ID or label of the head to be emptied",
					},
				},
				Action: app.Action(args.Optional("head"), with.RequiredFlags("head"), with.Auth, func(c *app.Context) error {
					if err := c.Client().EmptyHead(c.String("head")); err != nil {
						return err
					}

					log.Output("Head updated")

					return nil
				}),
			},
		},
	})
}
