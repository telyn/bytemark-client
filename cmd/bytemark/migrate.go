package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "migrate",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:        "disc",
				Usage:       "migrate a disc to a new storage pool",
				UsageText:   "bytemark --admin migrate disc <disc> [new_storage_pool]",
				Description: `This command migrates a disc to a new storage pool. If a new storage pool isn't supplied, a new one is picked automatically.`,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "disc",
						Usage: "the ID of the disc to migrate",
					},
					cli.StringFlag{
						Name:  "new_storage_pool",
						Usage: "the storage pool to move the disc to",
					},
				},
				Action: With(OptionalArgs("disc", "new_storage_pool"), RequiredFlags("disc"), AuthProvider, func(c *Context) (err error) {
					disc := c.Int("disc")
					storagePool := c.String("new_storage_pool")

					if err := global.Client.MigrateDisc(disc, storagePool); err != nil {
						return err
					}

					log.Outputf("%d migration initiated\n", disc)

					return nil
				}),
			},
			{
				Name:        "vm",
				Usage:       "migrate a VM to a new head",
				UsageText:   "bytemark --admin migrate vm <server> [new_head]",
				Description: `This command migrates a VM to a new head. If a new head isn't supplied, a new one is picked automatically.`,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "server",
						Usage: "the ID of the VM to migrate",
					},
					cli.StringFlag{
						Name:  "new_head",
						Usage: "the head to move the VM to",
					},
				},
				Action: With(OptionalArgs("server", "new_head"), RequiredFlags("server"), AuthProvider, func(c *Context) (err error) {
					vm := c.Int("server")
					head := c.String("new_head")

					if err := global.Client.MigrateVM(vm, head); err != nil {
						return err
					}

					log.Outputf("%d migration initiated\n", vm)

					return nil
				}),
			},
		},
	})
}
