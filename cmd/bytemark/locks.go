package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "lock",
		Usage:     "lock hardware profiles to prevent upgrading",
		UsageText: "lock hwprofile <server>",
		Description: `lock hardware profiles to prevent upgrading
		
This command locks the given server's hardware profile in place, preventing it from being automatically upgraded if a new is released. 'compatibility' hardware profiles are never automatically upgraded.`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "hwprofile",
			Usage:       "lock hardware profiles to prevent upgrading",
			UsageText:   "lock hwprofile <server>",
			Description: `This command locks the given server's hardware profile in place, preventing it from being automatically upgraded if a new is released. 'compatibility' hardware profiles are never automatically upgraded.`,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to lock",
					Value: new(app.VirtualMachineNameFlag),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) error {
				vmName := c.VirtualMachineName("server")
				return c.Client().SetVirtualMachineHardwareProfileLock(vmName, true)
			}),
		}},
	}, cli.Command{
		Name:      "unlock",
		Usage:     "unlock hardware profiles to allow upgrading",
		UsageText: "unlock hwprofile <server>",
		Description: `unlock hardware profiles to allow upgrading

		This command unlocks the given server's hardware profile, allowing it to be automatically upgraded if a new is released. 'compatibility' hardware profiles are never automatically upgraded.`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "hwprofile",
			Usage:       "unlock hardware profiles to allow upgrading",
			UsageText:   "unlock hwprofile <server>",
			Description: `This command unlocks the given server's hardware profile, allowing it to be automatically upgraded if a new is released. 'compatibility' hardware profiles are never automatically upgraded.`,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to unlock",
					Value: new(app.VirtualMachineNameFlag),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) error {
				vmName := c.VirtualMachineName("server")
				return c.Client().SetVirtualMachineHardwareProfileLock(vmName, false)
			}),
		}},
	})
}
