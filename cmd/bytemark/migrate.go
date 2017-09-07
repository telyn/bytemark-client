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
				UsageText:   "bytemark --admin migrate disc <disc> [new-storage-pool]",
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
				Action: With(OptionalArgs("disc", "new-storage-pool"), RequiredFlags("disc"), AuthProvider, func(c *Context) (err error) {
					disc := c.Int("disc")
					storagePool := c.String("new-storage-pool")

					if err := global.Client.MigrateDisc(disc, storagePool); err != nil {
						return err
					}

					log.Outputf("Migration for disc %d initiated\n", disc)

					return nil
				}),
			},
			{
				Name:        "server",
				Aliases:     []string{"vm"},
				Usage:       "migrate a server to a new head",
				UsageText:   "bytemark --admin migrate server <name> [new-head]",
				Description: `This command migrates a server to a new head. If a new head isn't supplied, a new one is picked automatically.`,
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server to migrate",
						Value: new(VirtualMachineNameFlag),
					},
					cli.StringFlag{
						Name:  "new-head",
						Usage: "the head to move the server to",
					},
				},
				Action: With(OptionalArgs("server", "new-head"), RequiredFlags("server"), AuthProvider, func(c *Context) (err error) {
					vm := c.VirtualMachineName("server")
					head := c.String("new-head")

					if err := global.Client.MigrateVirtualMachine(vm, head); err != nil {
						return err
					}

					log.Outputf("Migration for server %s initiated\n", vm.String())

					return nil
				}),
			},
		},
	})
}
