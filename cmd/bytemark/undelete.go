package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name: "undelete",
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
	})
}
