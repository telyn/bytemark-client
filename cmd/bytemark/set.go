package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
	"strconv"
)

func init() {
	commands = append(commands, cli.Command{
		Name: "set",
		Subcommands: []cli.Command{{
			Name: "cores",
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
				coresStr, err := c.NextArg()
				if err != nil {
					return err
				}
				cores, err := strconv.Atoi(coresStr)
				if err != nil || cores < 1 {
					log.Errorf("Invalid number of cores \"%s\"\r\n", coresStr)
					return &util.PEBKACError{}
				}
				return global.Client.SetVirtualMachineCores(c.VirtualMachineName, cores)

			}),
		}, {
			Name: "hwprofile",
			//lock_hwp := flags.Bool("lock", false, "")
			//unlock_hwp := flags.Bool("unlock", false, "")
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {

				/*// do nothing if --lock and --unlock are both specified
				    if *lock_hwp && *unlock_hwp {
					log.Log("Ambiguous command, both lock and unlock specified")
					cmds.HelpForSet()
					return util.E_PEBKAC
				    }*/
				profileStr, err := c.NextArg()
				if err != nil {
					return &util.PEBKACError{}
				}
				// if lock_hwp or unlock_hwp are specified, account this into the call
				/*if *lock_hwp {
					return global.Client.SetVirtualMachineHardwareProfile(c.VirtualMachineName, profileStr, true)
				} else if *unlock_hwp {
					return global.Client.SetVirtualMachineHardwareProfile(c.VirtualMachineName, profileStr, false)
					// otherwise omit lock
				} else {*/
				return global.Client.SetVirtualMachineHardwareProfile(c.VirtualMachineName, profileStr)
				/*}*/
			}),
		}, {
			Name: "memory",
			Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {

				memoryStr, err := c.NextArg()
				if err != nil {
					return &util.PEBKACError{}
				}

				memory, err := util.ParseSize(memoryStr)
				if err != nil || memory < 1 {
					log.Errorf("Invalid amount of memory \"%s\"\r\n", memoryStr)
					return &util.PEBKACError{}
				}

				return global.Client.SetVirtualMachineMemory(c.VirtualMachineName, memory)
			}),
		}},
	})
}
