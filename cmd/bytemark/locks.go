package main

import (
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "lock",
		Usage:     "lock hardware profiles to prevent upgrading",
		UsageText: "bytemark lock hwprofile <server>",
		Description: `lock hardware profiles to prevent upgrading
		
This command locks the given server's hardware profile in place, preventing it from being automatically upgraded if a new is released. 'compatibility' hardware profiles are never automatically upgraded.`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "hwprofile",
			Usage:       "lock hardware profiles to prevent upgrading",
			UsageText:   "bytemark lock hwprofile <server>",
			Description: `This command locks the given server's hardware profile in place, preventing it from being automatically upgraded if a new is released. 'compatibility' hardware profiles are never automatically upgraded.`,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to lock",
					Value: new(VirtualMachineNameFlag),
				},
			},
			Action: With(OptionalArgs("server"), RequiredFlags("server"), AuthProvider, func(c *Context) error {
				vmName := c.VirtualMachineName("server")
				return global.Client.SetVirtualMachineHardwareProfileLock(&vmName, true)
			}),
		}},
	}, cli.Command{
		Name:      "unlock",
		Usage:     "unlock hardware profiles to allow upgrading",
		UsageText: "bytemark unlock hwprofile <server>",
		Description: `unlock hardware profiles to allow upgrading

		This command unlocks the given server's hardware profile, allowing it to be automatically upgraded if a new is released. 'compatibility' hardware profiles are never automatically upgraded.`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "hwprofile",
			Usage:       "unlock hardware profiles to allow upgrading",
			UsageText:   "bytemark unlock hwprofile <server>",
			Description: `This command unlocks the given server's hardware profile, allowing it to be automatically upgraded if a new is released. 'compatibility' hardware profiles are never automatically upgraded.`,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to unlock",
					Value: new(VirtualMachineNameFlag),
				},
			},
			Action: With(OptionalArgs("server"), RequiredFlags("server"), AuthProvider, func(c *Context) error {
				vmName := c.VirtualMachineName("server")
				return global.Client.SetVirtualMachineHardwareProfileLock(&vmName, false)
			}),
		}},
	})
}
