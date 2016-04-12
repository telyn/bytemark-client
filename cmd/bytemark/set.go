package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"fmt"
	"github.com/codegangsta/cli"
	"strconv"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "set",
		Usage:       "Change hardware properties of Bytemark servers",
		UsageText:   "bytemark set cores|memory|hwprofile <server>",
		Description: `These commands set various hardware properties of Bytemark servers. Note that for memory increases, cores and hwprofile to take effect you will need to restart the server.`,
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "cores",
			Usage:       "Set the number of CPU cores on a Bytemark cloud server",
			UsageText:   "bytemark set cores <server name> <cores>",
			Description: "This command sets the number of CPU cores used by the cloud server. This will usually require a restart of the server to take effect.",
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
				coresStr, err := c.NextArg()
				if err != nil {
					return err
				}
				cores, err := strconv.Atoi(coresStr)
				if err != nil || cores < 1 {
					c.Help(fmt.Sprintf("Invalid number of cores \"%s\"\r\n", coresStr))
				}
				return global.Client.SetVirtualMachineCores(c.VirtualMachineName, cores)

			}),
		}, {
			Name:        "hwprofile",
			Usage:       "Set the hardware profile used by the cloud server",
			UsageText:   "bytemark set hwprofile <server> <profile>",
			Description: "This sets the hardware profile used. Hardware profiles can be simply thought of as what virtual motherboard you're using - generally you want a pretty recent one for maximum speed, but if you're running a very old or experimental OS (e.g. DOS or OS/2 or something) you may require the compatibility one. See `bytemark hwprofiles` for which ones are currently available.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "lock",
					Usage: "Locks the hardware profile (prevents it from being automatically upgraded when we release a newer version)",
				},
				cli.BoolFlag{
					Name:  "unlock",
					Usage: "Unlocks the hardware profile (allows it to be automatically upgraded when we release a newer version)",
				},
			},
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
				if c.Bool("lock") && c.Bool("unlock") {
					return c.Help("Ambiguous command, both lock and unlock specified")
				}

				profileStr, err := c.NextArg()
				if err != nil {
					return c.Help("No hardware profile name was specified")
				}
				if c.Bool("lock") {
					return global.Client.SetVirtualMachineHardwareProfile(c.VirtualMachineName, profileStr, true)
				} else if c.Bool("unlock") {
					return global.Client.SetVirtualMachineHardwareProfile(c.VirtualMachineName, profileStr, false)
				} else {
					return global.Client.SetVirtualMachineHardwareProfile(c.VirtualMachineName, profileStr)
				}
			}),
		}, {
			Name:        "memory",
			Usage:       "Sets the amount of memory the server has",
			UsageText:   "bytemark set memory <server> <memory size>",
			Description: "Memory is specified in GiB by default, but can be suffixed with an M to indicate that it is provided in MiB",
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {

				memoryStr, err := c.NextArg()
				if err != nil {
					return c.Help("No memory amount was specified")
				}

				memory, err := util.ParseSize(memoryStr)
				if err != nil || memory < 1 {
					return c.Help(fmt.Sprintf("Invalid amount of memory \"%s\"\r\n", memoryStr))
				}

				return global.Client.SetVirtualMachineMemory(c.VirtualMachineName, memory)
			}),
		}},
	})
}
