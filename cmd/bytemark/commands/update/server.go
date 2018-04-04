package update

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "server",
		Usage:     "update a server's configuration",
		UsageText: "update server [flags] <server>",
		Description: `Updates the configuration of an existing Cloud Server.

Note that for changes to memory or hardware profile to take effect you will need to restart the server.

Updating a server's name also allows it to be moved between groups and accounts you administer.

EXAMPLES

        bytemark update server --new-name boron oxygen
	        This will rename the server called oxygen in your default group to boron, still in your default group.

	bytemark update server --new-name sunglasses.development sunglasses
		This will move the server called sunglasses into the development group, keeping its name as sunglasses,

	bytemark update server --new-name rennes.bretagne.france charata.chaco.argentina
		This will move the server called charata in the chaco group in the argentina account, placing it in the bretagne group in the france account and rename it to rennes.`,
		Flags: append(app.OutputFlags("server", "object"),
			cli.GenericFlag{
				Name:  "memory",
				Value: new(util.SizeSpecFlag),
				Usage: "How much memory the server will have available, specified in GiB or with GiB/MiB units.",
			},
			cli.StringFlag{
				Name:  "hw-profile",
				Usage: "The hardware profile to use. See `bytemark profiles` for a list of hardware profiles available.",
			},
			cli.BoolFlag{
				Name:  "hw-profile-lock",
				Usage: "Locks the hardware profile (prevents it from being automatically upgraded when we release a newer version)",
			},
			cli.GenericFlag{
				Name:  "new-name",
				Usage: "A new name for the server",
				Value: new(app.VirtualMachineNameFlag),
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "The server to update",
				Value: new(app.VirtualMachineNameFlag),
			},
		),
		Action: app.Action(args.Optional("new-name", "hwprofile", "memory"), with.RequiredFlags("server"), with.Auth, updateServer),
	})
}

func updateMemory(c *app.Context) error {
	vmName := c.VirtualMachineName("server")
	memory := c.Size("memory")

	if memory == 0 {
		return nil
	}
	if c.VirtualMachine.Memory < memory {
		if !c.Bool("force") && !util.PromptYesNo(c.Prompter(), fmt.Sprintf("You're increasing the memory by %dGiB - this may cost more, are you sure?", (memory-c.VirtualMachine.Memory)/1024)) {
			return util.UserRequestedExit{}
		}
	}
	return c.Client().SetVirtualMachineMemory(vmName, memory)
}

func updateHwProfile(c *app.Context) error {
	vmName := c.VirtualMachineName("server")
	hwProfile := c.String("hw-profile")
	hwProfileLock := c.Bool("hw-profile-lock")

	if hwProfile == "" {
		if hwProfileLock {
			return c.Help("Must specify a hardware profile to lock")
		}
		return nil
	}
	return c.Client().SetVirtualMachineHardwareProfile(vmName, hwProfile, hwProfileLock)
}

func updateName(c *app.Context) error {
	vmName := c.VirtualMachineName("server")
	newName := c.VirtualMachineName("new-name")

	if newName.VirtualMachine == "" {
		return nil
	}
	fmt.Printf("GOT HERE! %v %+v\n", vmName, newName)
	return c.Client().MoveVirtualMachine(vmName, newName)
}

func updateServer(c *app.Context) error {
	for _, err := range []error{
		updateMemory(c),
		updateHwProfile(c),
		updateName(c), // needs to be last
	} {
		if err != nil {
			return err
		}
	}
	return nil
}
