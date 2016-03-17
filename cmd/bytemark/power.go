package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

/*
start: Starts a stopped server.

shutdown: Sends the ACPI shutdown signal, as if you had
          pressed the power/standby button. Allows the
          operating system to gracefully shut down.
          Hardware changes will be applied after the
          machine has been started again.

stop: Stops a running server, as if you had just pulled the
      cord out. Hardware changes will be applied when the
      machine has been started again.

restart: Stops and then starts a running server, as if you had
         pulled the cord out, then plugged it in and
         powered the machine on again.

reset: Instantly restarts a running server, as if you had
       pressed the reset button. Doesn't apply hardware
       changes.
*/

func init() {
	commands = append(commands, cli.Command{
		Name: "reset",
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
		Name: "restart",
		Action: With(VirtualMachineNameProvider, func(c *Context) (err error) {
			log.Logf("Attempting to restart %v...\r\n", c.VirtualMachineName)
			err = global.Client.RestartVirtualMachine(c.VirtualMachineName)
			if err != nil {
				return
			}

			log.Logf("%s restarted successfully.\r\n", c.VirtualMachineName)
			return
		}),
	}, cli.Command{
		Name: "shutdown",
		Action: With(VirtualMachineNameProvider, func(c *Context) (err error) {
			log.Logf("Attempting to shutdown %v...\r\n", c.VirtualMachineName)
			err = global.Client.ShutdownVirtualMachine(c.VirtualMachineName, true)
			if err != nil {
				return
			}

			log.Logf("%s was shutdown successfully.\r\n", c.VirtualMachineName)
			return
		}),
	}, cli.Command{
		Name: "start",
		Action: With(VirtualMachineNameProvider, func(c *Context) (err error) {
			log.Logf("Attempting to start %s...\r\n", c.VirtualMachineName)
			err = global.Client.StartVirtualMachine(c.VirtualMachineName)
			if err != nil {
				return
			}

			log.Logf("%s started successfully.\r\n", c.VirtualMachineName)
			return
		}),
	}, cli.Command{
		Name: "stop",
		Action: With(VirtualMachineNameProvider, func(c *Context) (err error) {
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
