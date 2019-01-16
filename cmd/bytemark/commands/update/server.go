package update

import (
	"errors"
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flagsets"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

// err.. this whole file is a bit of a fudge at the moment.
// a simpler UpdateVirtualMachine(VirtualMachineName, brain.VirtualMachine)
// would do the job of all the functions completely adequately

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "server",
		Usage:     "update a server's configuration",
		UsageText: "update server [flags] <server>",
		Description: `Updates the configuration of an existing Cloud Server.

    Note that for changes to cores, memory or hardware profile to take effect you will need to restart the server.

    Hardware profiles can be thought of as what virtual motherboard your server has, and in general it is best to use a
    pretty recent one for maximum speed. If the server is running an old or experimental OS without support for virtio
    drivers, or installing an older windows from an ISO without the virtio drivers compiled in, you may require the
    compatibility profile. See "bytemark show hwprofiles" for which ones are currently available.

    Memory is specified in GiB by default, but can be suffixed with an M to indicate that it is provided in MiB.

    Updating a server's name also allows it to be moved between groups and accounts you administer.

EXAMPLES
    bytemark update server --memory 768m --hwprofile virtio2018 small-server
        Changes small-server's memory to 768MiB, and its hwprofile to virtio2018

    bytemark update server --server app-dev --swap-ips-with app-production
        Swaps the primary IPs (those given to a server upon creation) with another server. Specify --swap-extra-ips to
        swap both primary and extra IPs with another server. For more granular IP alterations, use the panel for now or
        petition us to add an 'update ip'. Before swapping the IPs the servers ought to be reconfigured to expect the
        new IPs and then shut down, otherwise the serial/VNC console will have to be used (see bytemark-console) to
        configure the networking.

    bytemark update server --server app-dev --swap-ips-with app-production --swap-extra-ips
        Swaps all the IPs between app-dev and app-production. Before swapping the IPs the servers ought to be
        reconfigured to expect the new IPs and then shut down, otherwise the serial/VNC console will have to be used
        (see bytemark-console) to configure the networking.

    bytemark update server --new-name boron oxygen
        This will rename the server called oxygen in your default group to boron, still in your default group.

    bytemark update server --new-name sunglasses.development sunglasses
        This will move the server called sunglasses into the development group, keeping its name as sunglasses,

    bytemark update server --new-name rennes.bretagne.france charata.chaco.argentina
        This will move the server called charata in the chaco group in the argentina account, placing it in the bretagne
        group in the france account and rename it to rennes.`,
		Flags: append(app.OutputFlags("server", "object"),
			flagsets.Force,
			cli.GenericFlag{
				Name:  "memory",
				Value: new(flags.SizeSpecFlag),
				Usage: "How much memory the server will have available, specified in GiB or with GiB/MiB units.",
			},
			cli.StringFlag{
				Name:  "hwprofile",
				Usage: "The hardware profile to use. See `bytemark profiles` for a list of hardware profiles available.",
			},
			cli.BoolFlag{
				Name:  "lock-hwprofile",
				Usage: "Locks the hardware profile (prevents it from being automatically upgraded when we release a newer version)",
			},
			cli.BoolFlag{
				Name:  "unlock-hwprofile",
				Usage: "Locks the hardware profile (allows it to be automatically upgraded when we release a newer version)",
			},
			cli.GenericFlag{
				Name:  "new-name",
				Usage: "A new name for the server",
				Value: new(flags.VirtualMachineNameFlag),
			},
			cli.GenericFlag{
				Name:  "swap-ips-with",
				Usage: "A server to swap IP addresses with. Both v4 and v6 are swapped. See description and --swap-extra-ips",
				Value: new(flags.VirtualMachineNameFlag),
			},
			cli.BoolFlag{
				Name:  "swap-extra-ips",
				Usage: "Swaps extra IPs with the target of --swap-ips-with. When --swap-ips-with is unspecified, this flag is ignored",
			},
			cli.IntFlag{
				Name:  "cores",
				Usage: "the number of cores that should be available to the server",
			},
			cli.StringFlag{
				Name:  "cd-url",
				Usage: "An HTTP(S) URL for an ISO image file to attach as a cdrom",
			},
			cli.BoolFlag{
				Name:  "remove-cd",
				Usage: "Removes any current cdrom, as if the cd were ejected",
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "The server to update",
				Value: new(flags.VirtualMachineNameFlag),
			},
		),
		Action: app.Action(args.Optional("new-name", "hwprofile", "memory"),
			with.RequiredFlags("server"),
			with.VirtualMachine("server"),
			updateServer),
	})
}

func updateMemory(c *app.Context) error {
	vmName := flags.VirtualMachineName(c, "server")
	memory := flags.Size(c, "memory")

	if memory == 0 {
		return nil
	}
	if c.VirtualMachine.Memory < memory {
		if !flagsets.Forced(c) && !util.PromptYesNo(c.Prompter(), fmt.Sprintf("You're increasing the memory by %dGiB - this may cost more, are you sure?", (memory-c.VirtualMachine.Memory)/1024)) {
			return util.UserRequestedExit{}
		}
	}
	return c.Client().SetVirtualMachineMemory(vmName, memory)
}

func updateHwProfile(c *app.Context) error {
	vmName := flags.VirtualMachineName(c, "server")
	hwProfile := c.String("hwprofile")
	if hwProfile == "" {
		return nil
	}

	return c.Client().SetVirtualMachineHardwareProfile(vmName, hwProfile)
}

func updateLock(c *app.Context) error {
	server := flags.VirtualMachineName(c, "server")

	lockProfile := c.Bool("lock-hwprofile")
	unlockProfile := c.Bool("unlock-hwprofile")
	if lockProfile && unlockProfile {
		return errors.New("--lock-hwprofile and --unlock-hwprofile were both specified - only one may be specified at a time")
	} else if lockProfile {
		return c.Client().SetVirtualMachineHardwareProfileLock(server, true)
	} else if unlockProfile {
		return c.Client().SetVirtualMachineHardwareProfileLock(server, false)
	}
	return nil
}

func updateCores(c *app.Context) error {
	vmName := flags.VirtualMachineName(c, "server")
	cores := c.Int("cores")

	if cores == 0 {
		return nil
	}
	if c.VirtualMachine.Cores < cores {
		if !flagsets.Forced(c) && !util.PromptYesNo(c.Prompter(), fmt.Sprintf("You are increasing the number of cores from %d to %d. This may cause your VM to cost more, are you sure?", c.VirtualMachine.Cores, cores)) {
			return util.UserRequestedExit{}
		}
	}
	return c.Client().SetVirtualMachineCores(vmName, cores)
}

func updateName(c *app.Context) error {
	vmName := flags.VirtualMachineName(c, "server")
	newName := flags.VirtualMachineName(c, "new-name")

	if newName.VirtualMachine == "" {
		return nil
	}
	return c.Client().MoveVirtualMachine(vmName, newName)
}

func updateCdrom(c *app.Context) error {
	vmName := flags.VirtualMachineName(c, "server")
	cdURL := c.String("cd-url")
	removeCD := c.Bool("remove-cd")

	if cdURL == "" && !removeCD {
		return nil
	}
	err := c.Client().SetVirtualMachineCDROM(vmName, cdURL)
	if _, ok := err.(lib.InternalServerError); ok {
		return c.Help("Couldn't set the server's cdrom - check that you have provided a valid public HTTP url")
	}
	return err
}

func swapIPs(ctx *app.Context) error {
	if !ctx.IsSet("swap-ips-with") {
		return nil
	}

	// this is hacky - having to call GetVirtualMachine for the --server vm twice because
	// SwapVirtualMachineIPs calls it internally and we can't pass an ID to it,
	// when ctx.VirtualMachine.ID is already set.
	// TODO: once VirtualMachinePather is in, add VirtualMachineIDer and this
	// function somewhere (brainRequests? It does make a request after all).
	// Then rewrite the hack in the test.
	//
	//         func GetVirtualMachineID(client lib.Client, pather VirtualMachinePather) (id int, err error) {
	//             if ider, ok := pather.(VirtualMachineIDer); ok {
	//                 id = ider.VirtualMachineID()
	//                 return
	//             }
	//             var vm brain.VirtualMachine
	//             vm, err = GetVirtualMachine(client, pather)
	//             return vm.ID, err
	//         }
	server := flags.VirtualMachineName(ctx, "server")
	target := flags.VirtualMachineName(ctx, "swap-ips-with")
	swapExtra := ctx.Bool("swap-extra-ips")
	return brainRequests.SwapVirtualMachineIPs(ctx.Client(), server, target, swapExtra)
}

func updateServer(c *app.Context) error {
	for _, f := range [](func(*app.Context) error){
		swapIPs,
		updateMemory,
		updateHwProfile,
		updateLock,
		updateCores,
		updateCdrom,
		updateName, // needs to be last
	} {
		err := f(c)
		if err != nil {
			return err
		}
	}
	return nil
}
