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
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server that the disc is attached to",
					Value: new(VirtualMachineNameFlag),
				},
			},
			Action: With(OptionalArgs("server"), RequiredFlags("server"), VirtualMachineProvider("server"), func(c *Context) (err error) {
				vmName := c.VirtualMachineName("server")
				if !c.VirtualMachine.Deleted {
					log.Errorf("%s was already restored\r\n", c.VirtualMachine.Hostname)
					return
				}

				err = global.Client.UndeleteVirtualMachine(vmName)

				if err != nil {
					return
				}
				log.Logf("Successfully restored %s\r\n", c.VirtualMachine.Hostname)
				return
			}),
		}},
	})
}
