package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {
	undelete_server := With(VirtualMachineProvider, func(c *Context) (err error) {
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
	})
	commands = append(commands, cli.Command{
		Name:  "undelete",
		Usage: "Restores a previously deleted cloud server",
		Subcommands: []cli.Command{{
			Name:   "server",
			Usage:  "Restores a previously deleted cloud server",
			Action: undelete_server,
		}},
		Action: undelete_server,
	})
}
