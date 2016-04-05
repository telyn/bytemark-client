package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "undelete",
		Usage:       "Restores a previously deleted cloud server",
		UsageText:   "bytemark undelete server <name>",
		Action:      cli.ShowSubcommandHelp,
		Description: `This command restores a previously deleted cloud server. Note that it cannot be used to restore a server that has been permanently deleted (purged).`,
		Subcommands: []cli.Command{{
			Name:        "server",
			Usage:       "Restores a previously deleted cloud server",
			UsageText:   "bytemark undelete server <name>",
			Description: `This command restores a previously deleted cloud server. Note that it cannot be used to restore a server that has been permanently deleted (purged).`,
			Action: With(VirtualMachineProvider, func(c *Context) (err error) {
				if !c.VirtualMachine.Deleted {
					log.Errorf("%s was already undeleted\r\n", c.VirtualMachine.Hostname)
					return
				}

				err = global.Client.UndeleteVirtualMachine(c.VirtualMachineName)

				if err != nil {
					return
				}
				log.Logf("Successfully restored %s\r\n", c.VirtualMachine.Hostname)
				return
			}),
		}},
	})
}
