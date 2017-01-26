package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "restore",
		Usage:       "restores a previously deleted cloud server",
		UsageText:   "bytemark restore server <name>",
		Description: `restores a previously deleted cloud server`,

		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:      "server",
			Usage:     "restores a previously deleted cloud server",
			UsageText: "bytemark restore server <name>",
			Description: `This command restores a previously deleted cloud server to its non-deleted state.
Note that it cannot be used to restore a server that has been permanently deleted (purged).`,
			Action: With(VirtualMachineProvider, func(c *Context) (err error) {
				if !c.VirtualMachine.Deleted {
					log.Errorf("%s was already restored\r\n", c.VirtualMachine.Hostname)
					return
				}

				err = global.Client.UndeleteVirtualMachine(c.VirtualMachineName)

				if err != nil {
					return
				}
				log.Logf("Successfully restored %s\r\n", c.VirtualMachine.Hostname)
				return
			}),
		}, {
			Name:        "backup",
			Usage:       "restore the given backup",
			UsageText:   `bytemark restore backup <server name> <disc label> <backup label>`,
			Description: "Restores the given backup. Before doing this, a new backup is made of the disc's current state.",
			Action: With(VirtualMachineNameProvider, DiscLabelProvider, func(c *Context) (err error) {
				backup, err := c.NextArg()
				if err != nil {
					return
				}
				// TODO(telyn): eventually RestoreBackup will return backups as the first argument. We should process that and output info :)
				_, err = global.Client.RestoreBackup(*c.VirtualMachineName, *c.DiscLabel, backup)
				if err != nil {
					return
				}
				log.Logf("Disc '%s' is now being restored from backup '%s'", *c.DiscLabel, backup)
				return
			}),
		}},
	})
}
