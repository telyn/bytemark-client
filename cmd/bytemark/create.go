package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"strings"
)

func init() {
	createServerCmd := cli.Command{
		Name:      "server",
		Usage:     `create a new server with bytemark`,
		UsageText: "bytemark create server [flags] <name> [<cores> [<memory [<disc specs>]...]]",
		Description: `Creates a Cloud Server with the given specification, defaulting to a basic server with Symbiosis installed.
		
A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs

If hwprofile-locked is set then the cloud server's virtual hardware won't be changed over time.`,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "cores",
				Value: 1,
				Usage: "Number of CPU cores",
			},
			cli.StringFlag{
				Name:  "cdrom",
				Usage: "URL pointing to an ISO which will be attached to the cloud server as a CD",
			},
			cli.GenericFlag{
				Name:  "disc",
				Usage: "One or more disc specifications. Defaults to a single 25GiB sata-grade disc",
				Value: new(util.DiscSpecFlag),
			},
			forceFlag,
			cli.StringFlag{
				Name:  "hwprofile",
				Usage: "The hardware profile to use. Defaults to the current modern profile. See `bytemark profiles` for a list of hardware profiles available.",
			},
			cli.BoolFlag{
				Name:  "hwprofile-locked",
				Usage: "If set, the hardware profile will be 'locked', meaning that when Bytemark updates the hardware profiles your VM will keep its current one.",
			},
			cli.GenericFlag{
				Name:  "ip",
				Value: new(util.IPFlag),
				Usage: "Specify an IPv4 or IPv6 address to use. This will only be useful if you are creating the machine in a private VLAN.",
			},
			cli.BoolFlag{
				Name:  "json",
				Usage: "If set, will output the spec and created virtual machine as a JSON object.",
			},
			cli.GenericFlag{
				Name:  "memory",
				Value: new(util.SizeSpecFlag),
				Usage: "How much memory the server will have available, specified in GiB or with GiB/MiB units. Defaults to 1GiB.",
			},
			cli.BoolFlag{
				Name:  "no-image",
				Usage: "Specifies that the server should not be imaged.",
			},
			cli.BoolFlag{
				Name:  "no-discs",
				Usage: "Specifies that the server should not have discs.",
			},
			cli.BoolFlag{
				Name:  "stopped",
				Usage: "If set, the server will not be started, even to image it.",
			},
			cli.StringFlag{
				Name:  "zone",
				Usage: "Which zone the server will be created in. See `bytemark zones` for the choices.",
			},
		},

		Action: With(VirtualMachineNameProvider, AuthProvider, createServer),
	}
	for _, flag := range imageInstallFlags {
		createServerCmd.Flags = append(createServerCmd.Flags, flag)
	}

	createDiscsCmd := cli.Command{
		Name:    "discs",
		Aliases: []string{"disc", "disk", "disks"},
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "disc",
				Usage: "A disc to add. You can specify as many discs as you like by adding more --disc flags.",
				Value: new(util.DiscSpecFlag),
			},
			forceFlag,
		},
		Usage:     "create virtual discs attached to one of your cloud servers",
		UsageText: "bytemark create discs [--disc <disc spec>]... <cloud server>",
		Description: `A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs`,
		Action: With(VirtualMachineNameProvider, AuthProvider, createDiscs),
	}

	createGroupCmd := cli.Command{
		Name:        "group",
		Usage:       "create a group for organising your servers",
		UsageText:   "bytemark create group <group name>",
		Description: `Groups are part of your server's fqdn`,
		Action:      With(GroupNameProvider, AuthProvider, createGroup),
	}

	commands = append(commands, cli.Command{
		Name:      "create",
		Usage:     "creates servers, discs, etc - see `bytemark create <kind of thing> help`",
		UsageText: "bytemark create disc|group|ip|server",
		Description: `create a new disc, group, IP or server

	create disc[s] [--disc <disc spec>]... <cloud server>
	create group [--account <name>] <name>
	create ip [--reason reason] <cloud server>
	create server (see bytemark create server help)

A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			createServerCmd,
			createDiscsCmd,
			createGroupCmd,
		},
	})
}

func createDiscs(c *Context) (err error) {
	discs := c.Discs("disc")

	for i := range discs {
		d, err := discs[i].Validate()
		if err != nil {
			return err
		}
		discs[i] = *d
	}

	log.Logf("Adding %d discs to %s:\r\n", len(discs), c.VirtualMachineName)
	for _, d := range discs {
		log.Logf("    %dGiB %s...", d.Size/1024, d.StorageGrade)
		err := global.Client.CreateDisc(c.VirtualMachineName, d)
		if err != nil {
			log.Errorf("failure! %v\r\n", err.Error())
		} else {
			log.Log("success!")
		}
	}
	return
}

func createGroup(c *Context) (err error) {
	err = global.Client.CreateGroup(c.GroupName)
	if err == nil {
		log.Logf("Group %s was created under account %s\r\n", c.GroupName.Group, c.GroupName.Account)
	}
	return
}

// createServerReadArgs sets up the initial defaults, reads in the --disc, --cores and --memory flags, then reads in positional arguments for the command line.
func createServerReadArgs(c *Context) (discs []brain.Disc, cores, memory int, err error) {

	discs = c.Discs("disc")
	cores = c.Int("cores")
	memory = c.Size("memory")
	if memory == 0 {
		memory = 1024
	}

	for argNum, arg := range c.Args() {
		switch argNum {
		case 0: // cores
			tmpCores, coresErr := strconv.Atoi(arg)
			if coresErr != nil {
				err = coresErr
				return
			}
			cores = tmpCores
		case 1: // memory
			tmpMem, memErr := util.ParseSize(arg)
			if memErr != nil {
				err = memErr
				return
			}
			memory = tmpMem
		case 2: // disc
			discs = make([]brain.Disc, strings.Count(arg, ",")+1)
			for discNum, discSpec := range strings.Split(arg, ",") {
				tmpDisc, discErr := util.ParseDiscSpec(discSpec)
				if discErr != nil {
					err = discErr
					return
				}
				discs[discNum] = *tmpDisc
			}
		case 3:
			err = c.Help("Too many arguments given.")
			return
		}
	}
	return
}

// createServerReadIPs reads the IP flags and creates an IPSpec
func createServerReadIPs(c *Context) (ipspec *brain.IPSpec, err error) {
	ips := c.IPs("ip")

	if len(ips) > 2 {
		err = c.Help("A maximum of one IPv4 and one IPv6 address may be specified")
		return
	}

	if len(ips) > 0 {
		ipspec = &brain.IPSpec{}

		for _, ip := range ips {
			if ip.To4() != nil {
				if ipspec.IPv4 != "" {
					err = c.Help("A maximum of one IPv4 and one IPv6 address may be specified")
					return
				}
				ipspec.IPv4 = ip.To4().String()
			} else {
				if ipspec.IPv6 != "" {
					err = c.Help("A maximum of one IPv4 and one IPv6 address may be specified")
					return

				}
				ipspec.IPv6 = ip.String()
			}
		}
	}
	return
}

func createServerPrepSpec(c *Context) (spec brain.VirtualMachineSpec, err error) {
	noImage := c.Bool("no-image")
	if c.Bool("no-discs") {
		noImage = true
	}

	discs, cores, memory, err := createServerReadArgs(c)
	if err != nil {
		return
	}

	if len(discs) == 0 && !c.Context.Bool("no-discs") {
		discs = append(discs, brain.Disc{Size: 25600})
	}

	for i := range discs {
		d, discErr := discs[i].Validate()
		if discErr != nil {
			return spec, discErr
		}
		discs[i] = *d
	}

	ipspec, err := createServerReadIPs(c)
	if err != nil {
		return
	}

	imageInstall, _, err := prepareImageInstall(c)
	if err != nil {
		return
	}

	if noImage {
		imageInstall = nil
	}

	stopped := c.Bool("stopped")
	cdrom := c.String("cdrom")

	// if stopped isn't set and either cdrom or image are set, start the server
	autoreboot := !stopped && ((imageInstall != nil) || (cdrom != ""))

	spec = brain.VirtualMachineSpec{
		VirtualMachine: &brain.VirtualMachine{
			Name:                  c.VirtualMachineName.VirtualMachine,
			Autoreboot:            autoreboot,
			Cores:                 cores,
			Memory:                memory,
			ZoneName:              c.String("zone"),
			CdromURL:              c.String("cdrom"),
			HardwareProfile:       c.String("hwprofile"),
			HardwareProfileLocked: c.Bool("hwprofile-locked"),
		},
		Discs:   discs,
		IPs:     ipspec,
		Reimage: imageInstall,
	}
	return
}

func createServer(c *Context) (err error) {
	spec, err := createServerPrepSpec(c)
	if err != nil {
		return
	}

	groupName := c.VirtualMachineName.GroupName()

	log.Log("The following server will be created:")
	err = lib.FormatVirtualMachineSpec(os.Stderr, groupName, &spec, "specfull")
	if err != nil {
		return err
	}

	// If we're not forcing, prompt. If the prompt comes back false, exit.
	if !c.Bool("force") && !util.PromptYesNo("Are you certain you wish to continue?") {
		log.Error("Exiting.")
		return util.UserRequestedExit{}
	}

	_, err = global.Client.CreateVirtualMachine(groupName, spec)
	if err != nil {
		return err
	}
	vm, err := global.Client.GetVirtualMachine(c.VirtualMachineName)
	if err != nil {
		return
	}
	return c.IfNotMarshalJSON(map[string]interface{}{"spec": spec, "virtual_machine": vm}, func() (err error) {
		log.Log("cloud server created successfully")
		err = lib.FormatVirtualMachine(os.Stderr, vm, "server_full")
		if err != nil {
			return
		}
		if spec.Reimage != nil {
			log.Log()
			log.Logf("Root password: ") // logf so we don't get a tailing \r\n
			log.Outputf("%s\r\n", spec.Reimage.RootPassword)
		} else {
			log.Log("Machine was not imaged")
		}
		return
	})
}
