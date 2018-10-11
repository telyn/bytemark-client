package commands

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/wait"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "restart",
		Usage:       "power off a server and start it again",
		UsageText:   "restart server [--rescue | --appliance <appliance>] <server>",
		Description: "This command will power down a server and then start it back up again.",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:      "server",
			Usage:     "power off a server and start it again",
			UsageText: "restart server [--rescue | --appliance <appliance>] <server>",
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

				c.Log("Shutting down %v...", vmName)
				err = c.Client().ShutdownVirtualMachine(vmName, true)
				if err != nil {
					return
				}
				err = wait.VMPowerOff(c, vmName)
				if err != nil {
					return
				}

				c.Log("Done!\n\nStarting %s back up.", vmName)
				if appliance != "" {
					err = brainRequests.StartVirtualMachineWithAppliance(c.Client(), vmName, appliance)
					c.Log("Server has now started. Use `bytemark console %v` or visit %v to connect.", c.String("server"), c.Config().PanelURL())
				} else {
					err = c.Client().StartVirtualMachine(vmName)
				}

				return
			}),
		}},
	})
}
