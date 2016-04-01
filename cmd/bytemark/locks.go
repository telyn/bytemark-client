package main

import (
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name: "lock",
		Subcommands: []cli.Command{{
			Name: "hwprofile",
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
				return global.Client.SetVirtualMachineHardwareProfileLock(c.VirtualMachineName, true)
			}),
		}},
	}, cli.Command{
		Name: "unlock",
		Subcommands: []cli.Command{{
			Name: "hwprofile",
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
				return global.Client.SetVirtualMachineHardwareProfileLock(c.VirtualMachineName, false)
			}),
		}},
	})
}
