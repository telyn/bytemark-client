package main

import (
	"fmt"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "reset",
		Usage:       "restart a server as though the reset button had been pushed",
		UsageText:   "bytemark reset <server>",
		Description: "For cloud servers, this does not cause the qemu process to be restarted. This means that the server will remain on the same head and will not notice hardware changes.",
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to reset",
				Value: new(app.VirtualMachineNameFlag),
			},
		},
		Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
			vmName := c.VirtualMachineName("server")
			log.Logf("Attempting to reset %v...\r\n", vmName)
			err = c.Client().ResetVirtualMachine(vmName)
			if err != nil {
				return err
			}

			log.Errorf("%v reset successfully.\r\n", vmName)
			return
		}),
	}, cli.Command{
		Name:        "restart",
		Usage:       "power off a server and start it again",
		UsageText:   "bytemark restart <server> [--rescue || --appliance <appliance>]",
		Description: "This command will power down a server and then start it back up again.",
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to restart",
				Value: new(app.VirtualMachineNameFlag),
			},
			cli.BoolFlag{
				Name:  "rescue",
				Usage: "boots the server using the rescue appliance",
			},
			cli.StringFlag{
				Name:  "appliance",
				Usage: "the appliance to boot into when the server starts",
			},
		},
		Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
			vmName := c.VirtualMachineName("server")
			appliance := c.String("appliance")

			if appliance != "" && c.Bool("rescue") {
				return fmt.Errorf("--appliance and --rescue have both been set when only one is allowed")
			}

			if c.Bool("rescue") {
				appliance = "rescue"
			}

			fmt.Fprintf(c.App().Writer, "Shutting down %v...", vmName)
			err = c.Client().ShutdownVirtualMachine(vmName, true)
			if err != nil {
				return
			}
			err = waitForShutdown(c, vmName)
			if err != nil {
				return
			}

			c.Log("Done!\n\nStarting %s back up.", vmName)
			return startServer(c, vmName, appliance)
		}),
	}, cli.Command{
		Name:        "shutdown",
		Usage:       "cleanly shut down a server",
		UsageText:   "bytemark shutdown <server>",
		Description: "This command sends the ACPI shutdown signal to the server, causing a clean shut down. This is like pressing the power button on a computer you have physical access to.",
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to shutdown",
				Value: new(app.VirtualMachineNameFlag),
			},
		},
		Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
			vmName := c.VirtualMachineName("server")
			fmt.Fprintf(c.App().Writer, "Shutting down %v...", vmName)
			err = c.Client().ShutdownVirtualMachine(vmName, true)
			if err != nil {
				return
			}

			err = waitForShutdown(c, vmName)
			if err != nil {
				return
			}

			c.Log("Done!", vmName)
			return
		}),
	}, cli.Command{
		Name:        "start",
		Usage:       "start a stopped server",
		UsageText:   "bytemark start <server>",
		Description: "This command will start a server that is not currently running.",
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to start",
				Value: new(app.VirtualMachineNameFlag),
			},
			cli.BoolFlag{
				Name:  "rescue",
				Usage: "boots the server using the rescue appliance",
			},
			cli.StringFlag{
				Name:  "appliance",
				Usage: "the appliance to boot into when the server starts",
			},
		},
		Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
			vmName := c.VirtualMachineName("server")
			appliance := c.String("appliance")

			if appliance != "" && c.Bool("rescue") {
				return fmt.Errorf("--appliance and --rescue have both been set when only one is allowed")
			}

			if c.Bool("rescue") {
				appliance = "rescue"
			}
			c.LogErr("Attempting to start %s...\r\n", vmName)
			return startServer(c, vmName, appliance)
		}),
	}, cli.Command{
		Name:        "stop",
		Usage:       "stop a server, as though pulling the power cable out",
		UsageText:   "bytemark stop <server>",
		Description: "This command will instantly power down a server. Note that this may cause data loss, particularly on servers with unjournaled file systems (e.g. ext2)",
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to stop",
				Value: new(app.VirtualMachineNameFlag),
			},
		},
		Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
			vmName := c.VirtualMachineName("server")
			log.Logf("Attempting to stop %s...\r\n", vmName)
			err = c.Client().StopVirtualMachine(vmName)
			if err != nil {
				return
			}

			log.Logf("%s stopped successfully.\r\n", vmName)
			return
		}),
	})
}
func waitForShutdown(c *app.Context, name lib.VirtualMachineName) (err error) {
	vm := brain.VirtualMachine{PowerOn: true}

	for vm.PowerOn {
		if !c.IsTest() {
			time.Sleep(5 * time.Second)
		}
		fmt.Fprint(c.App().Writer, ".")

		vm, err = c.Client().GetVirtualMachine(name)
		if err != nil {
			return
		}
	}
	return
}

func startServer(c *app.Context, vmName lib.VirtualMachineName, appliance string) (err error) {
	if appliance != "" {
		err = brainMethods.StartVirtualMachineWithAppliance(c.Client(), vmName, appliance)
	} else {
		err = c.Client().StartVirtualMachine(vmName)
	}
	if err == nil {
		c.Log("Server has been started.")
		if appliance != "" {
			c.Log("Use `bytemark console %v` or visit %v to connect.", c.String("server"), c.Config().PanelURL())
		}
	}
	return
}
