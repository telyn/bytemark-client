package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util/sizespec"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/urfave/cli"
	"strconv"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "set",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "disc",
				Usage:     "disc related admin operations",
				UsageText: "bytemark --admin set disc command [command options] [arguments...]",
				Action:    cli.ShowSubcommandHelp,
				Subcommands: []cli.Command{
					{
						Name:      "iops_limit",
						Usage:     "set the IOPS limit of a disc",
						UsageText: "bytemark --admin set disc iops_limit <server> <disc> <limit>",
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:  "disc",
								Usage: "the name of the disc to alter the iops_limit of",
							},
							cli.GenericFlag{
								Name:  "server",
								Usage: "the server the disc belongs to",
								Value: new(VirtualMachineNameFlag),
							},
						},
						Action: With(OptionalArgs("server", "disc"), AuthProvider, func(c *Context) error {
							iopsLimitStr, err := c.NextArg()
							if err != nil {
								return err
							}

							iopsLimit, err := strconv.Atoi(iopsLimitStr)
							if err != nil || iopsLimit < 1 {
								return c.Help(fmt.Sprintf("Invalid number for IOPS limit \"%s\"\r\n", iopsLimitStr))
							}
							vmName := c.VirtualMachineName("server")

							return global.Client.SetDiscIopsLimit(&vmName, c.String("disc"), iopsLimit)
						}),
					},
				},
			},
		},
	})

	commands = append(commands, cli.Command{
		Name:      "set",
		Usage:     "change hardware properties of Bytemark servers",
		UsageText: "bytemark set cores|memory|hwprofile <server>",
		Description: `change hardware properties of Bytemark servers
		
These commands set various hardware properties of Bytemark servers. Note that for memory increases, cores and hwprofile to take effect you will need to restart the server.`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "cdrom",
				Usage:     "attach a cdrom to your Bytemark Cloud Server",
				UsageText: "bytemark attach cdrom <server> <cdurl>",
				Description: `attach a cdrom to your Bytemark Cloud Server

This command allows you to add a cdrom to your Bytemark server. The CD must be publicly available over HTTP in order to be attached.`,
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server to attach the CD to",
						Value: new(VirtualMachineNameFlag),
					},
				},
				Action: With(OptionalArgs("server"), AuthProvider, func(c *Context) error {
					// url should be a flag
					url, err := c.NextArg()
					if err != nil {
						return err
					}
					vmName := c.VirtualMachineName("server")
					err = global.Client.SetVirtualMachineCDROM(&vmName, url)
					if _, ok := err.(lib.InternalServerError); ok {
						return c.Help("Couldn't set the server's cdrom - check that you have provided a valid public HTTP url")
					}
					return err

				}),
			},
			{
				Name:        "cores",
				Usage:       "set the number of CPU cores on a Bytemark cloud server",
				UsageText:   "bytemark set cores <server name> <cores>",
				Description: "This command sets the number of CPU cores used by the cloud server. This will usually require a restart of the server to take effect.",
				Flags: []cli.Flag{
					forceFlag,
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server to alter",
						Value: new(VirtualMachineNameFlag),
					},
				},
				Action: With(OptionalArgs("server"), VirtualMachineProvider("server"), func(c *Context) error {
					// cores should be a flag
					vmName := c.VirtualMachineName("server")
					coresStr, err := c.NextArg()
					if err != nil {
						return err
					}
					cores, err := strconv.Atoi(coresStr)
					if err != nil || cores < 1 {
						return c.Help(fmt.Sprintf("Invalid number of cores \"%s\"\r\n", coresStr))
					}
					if c.VirtualMachine.Cores < cores {
						if !c.Bool("force") && !util.PromptYesNo(fmt.Sprintf("You are increasing the number of cores from %d to %d. This may cause your VM to cost more, are you sure?", c.VirtualMachine.Cores, cores)) {
							return util.UserRequestedExit{}
						}
					}
					return global.Client.SetVirtualMachineCores(&vmName, cores)

				}),
			}, {
				Name:        "hwprofile",
				Usage:       "set the hardware profile used by the cloud server",
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
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server whose hardware profile you wish to alter",
						Value: new(VirtualMachineNameFlag),
					},
				},
				Action: With(OptionalArgs("server"), AuthProvider, func(c *Context) error {
					if c.Bool("lock") && c.Bool("unlock") {
						return c.Help("Ambiguous command, both lock and unlock specified")
					}

					profileStr, err := c.NextArg()
					if err != nil {
						return c.Help("No hardware profile name was specified")
					}
					vmName := c.VirtualMachineName("server")
					if c.Bool("lock") {
						return global.Client.SetVirtualMachineHardwareProfile(&vmName, profileStr, true)
					} else if c.Bool("unlock") {
						return global.Client.SetVirtualMachineHardwareProfile(&vmName, profileStr, false)
					} else {
						return global.Client.SetVirtualMachineHardwareProfile(&vmName, profileStr)
					}
				}),
			}, {
				Name:        "memory",
				Usage:       "sets the amount of memory the server has",
				UsageText:   "bytemark set memory <server> <memory size>",
				Description: "Memory is specified in GiB by default, but can be suffixed with an M to indicate that it is provided in MiB",
				Flags: []cli.Flag{
					forceFlag,
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server to alter",
						Value: new(VirtualMachineNameFlag),
					},
				},
				Action: With(OptionalArgs("server"), VirtualMachineProvider("server"), func(c *Context) error {

					memoryStr, err := c.NextArg()
					if err != nil {
						return c.Help("No memory amount was specified")
					}

					memory, err := sizespec.Parse(memoryStr)
					if err != nil || memory < 1 {
						return c.Help(fmt.Sprintf("Invalid amount of memory \"%s\"\r\n", memoryStr))
					}

					if c.VirtualMachine.Memory < memory {
						if !c.Bool("force") && !util.PromptYesNo(fmt.Sprintf("You're increasing the memory by %dGiB - this may cost more, are you sure?", (memory-c.VirtualMachine.Memory)/1024)) {
							return util.UserRequestedExit{}
						}
					}

					vmName := c.VirtualMachineName("server")
					return global.Client.SetVirtualMachineMemory(&vmName, memory)
				}),
			}},
	})
}
