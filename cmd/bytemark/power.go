package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "reset",
		Usage:       "Restart a server as though the reset button had been pushed",
		UsageText:   "bytemark reset <server>",
		Description: "For cloud servers, this does not cause the qemu process to be restarted. This means that the server will remain on the same head and will not notice hardware changes.",
		Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) (err error) {
			log.Logf("Attempting to reset %v...\r\n", c.VirtualMachineName)
			err = global.Client.ResetVirtualMachine(c.VirtualMachineName)
			if err != nil {
				return err
			}

			log.Errorf("%v reset successfully.\r\n", c.VirtualMachineName)
			return
		}),
	}, cli.Command{
		Name:        "restart",
		Usage:       "Power off a server and start it again",
		UsageText:   "bytemark restart <server>",
		Description: "This command will power down a server, cleanly if possible, and then start it back up again. For cloud servers this can cause the server to be started on a different head to the one it was running on, but this is not guaranteed.",
		Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) (err error) {
			log.Logf("Attempting to restart %v...\r\n", c.VirtualMachineName)
			err = global.Client.RestartVirtualMachine(c.VirtualMachineName)
			if err != nil {
				return
			}

			log.Logf("%s restarted successfully.\r\n", c.VirtualMachineName)
			return
		}),
	}, cli.Command{
		Name:        "shutdown",
		Usage:       "Cleanly shut down a server",
		UsageText:   "bytemark shutdown <server>",
		Description: "This command sends the ACPI shutdown signal to the server, causing a clean shut down. This is like pressing the power button on a computer you have physical access to.",
		Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) (err error) {
			log.Logf("Attempting to shutdown %v...\r\n", c.VirtualMachineName)
			err = global.Client.ShutdownVirtualMachine(c.VirtualMachineName, true)
			if err != nil {
				return
			}

			log.Logf("%s was shutdown successfully.\r\n", c.VirtualMachineName)
			return
		}),
	}, cli.Command{
		Name:        "start",
		Usage:       "Start a stopped server",
		UsageText:   "bytemark start <server>",
		Description: "This command will start a server that is not currently running.",
		Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) (err error) {
			log.Logf("Attempting to start %s...\r\n", c.VirtualMachineName)
			err = global.Client.StartVirtualMachine(c.VirtualMachineName)
			if err != nil {
				return
			}

			log.Logf("%s started successfully.\r\n", c.VirtualMachineName)
			return
		}),
	}, cli.Command{
		Name:        "stop",
		Usage:       "Stop a server, as though pulling the power cable out.",
		UsageText:   "bytemark stop <server>",
		Description: "This command will instantly power down a server. Note that this may cause data loss, particularly on servers with unjournaled file systems (e.g. ext2)",
		Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) (err error) {
			log.Logf("Attempting to stop %s...\r\n", c.VirtualMachineName)
			err = global.Client.StopVirtualMachine(c.VirtualMachineName)
			if err != nil {
				return
			}

			log.Logf("%s stopped successfully.\r\n", c.VirtualMachineName)
			return
		}),
	})
}
