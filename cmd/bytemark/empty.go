package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
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
				Action: With(OptionalArgs("storage-pool"), RequiredFlags("storage-pool"), AuthProvider, func(c *Context) error {
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
				Action: With(OptionalArgs("head"), RequiredFlags("head"), AuthProvider, func(c *Context) error {
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
