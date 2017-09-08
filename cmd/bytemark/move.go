package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "move",
		Usage:     "rename a server and move it across groups and accounts",
		UsageText: "bytemark move <old name> <new name>",
		Description: `This command renames a server and moves it between groups/accounts. You may only move servers between accounts you are an administrator of.
		
EXAMPLES

	bytemark move oxygen boron
		This will rename the server called oxygen in your default group to boron, still in your default group
		
	bytemark move sunglasses sunglasses.development
		This will move the server called sunglasses into the development group, keeping its name as sunglasses
		
	bytemark move charata.chaco.argentina rennes.bretagne.france
		This will move the server called charata in the chaco group in the argentina account, placing it in the bretagne group in the france account and rename it to rennes.`,
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "from",
				Usage: "the server to move",
				Value: new(VirtualMachineNameFlag),
			},
			cli.GenericFlag{
				Name:  "to",
				Usage: "the new name for the server",
				Value: new(VirtualMachineNameFlag),
			},
		},
		Action: With(OptionalArgs("from", "to"), RequiredFlags("from", "to"), AuthProvider, func(c *Context) (err error) {
			from := c.VirtualMachineName("from")
			to := c.VirtualMachineName("to")

			err = c.Client().MoveVirtualMachine(from, to)
			if err != nil {
				log.Output("Couldn't rename server.")
				return err
			}
			log.Outputf("Successfully moved %v to %v\r\n", from, to)
			return nil
		}),
	})
}
