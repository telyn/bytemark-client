package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {
	createServer := cli.Command{
		Name:      "server",
		Usage:     `Create a new server with bytemark.`,
		UsageText: "bytemark create server [flags] <name> [<cores> [<memory [<disc specs>]...]]",
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
				Usage: "One of more disc specifications. Defaults to a single 25GiB sata-grade disc",
				Value: new(util.DiscSpecFlag),
			},
			cli.GenericFlag{
				Name:  "firstboot-script-file",
				Usage: "Path to a script which will be run the first time the server boots after imaging",
			},
			cli.BoolFlag{
				Name:  "force",
				Usage: "Disables the confirmation prompt",
			},
			cli.StringFlag{
				Name:  "hwprofile",
				Usage: "The hardware profile to use. Defaults to the current modern profile. See `bytemark profiles` for a list of hardware profiles available.",
			},
			cli.BoolFlag{
				Name:  "hwprofile-locked",
				Usage: "If set, the hardware profile will be 'locked', meaning that when Bytemark updates the hardware profiles your VM will keep its current one.",
			},
			cli.StringFlag{
				Name:  "image",
				Usage: "Which operating system image to use. See `bytemark images` for a list of images available.",
			},
			cli.GenericFlag{
				Name:  "ip",
				Value: new(util.IPFlag),
				Usage: "Specify an IPv4 or IPv6 address to use. This will only be useful if you are creating the machine in a private VLAN.",
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
			cli.GenericFlag{
				Name:  "authorized-keys-file",
				Value: new(util.FileFlag),
				Usage: "Specifies SSH authorized keys for the root user. Only affects linux images.",
			},
			cli.StringFlag{
				Name:  "authorized-keys",
				Usage: "Specifies SSH authorized keys for the root user. Only affects linux images.",
			},
			cli.StringFlag{
				Name:  "root-password",
				Usage: "Specifies the password for the root/Administrator user.",
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

		Description: `Creates a Cloud Server with the given specification, defaulting to a basic server with Symbiosis installed.
		
A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs

If hwprofile-locked is set then the cloud server's virtual hardware won't be changed over time.`,
		Action: With(VirtualMachineNameProvider, AuthProvider, fn_createServer),
	}

	createDiscs := cli.Command{
		Name:      "discs",
		Aliases:   []string{"disc", "disk", "disks"},
		Usage:     "create virtual discs attached to one of your cloud servers",
		UsageText: "bytemark create discs [--disc <disc spec>]... <cloud server>",
		Description: `A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs`,
		Action: With(VirtualMachineNameProvider, AuthProvider, fn_createDisc),
	}

	createGroup := cli.Command{
		Name:   "group",
		Usage:  "bytemark create group <group name>",
		Action: With(GroupNameProvider, AuthProvider, fn_createGroup),
	}

	commands = append(commands, cli.Command{
		Name:      "create",
		Usage:     "bytemark create disc|group|ip|server",
		UsageText: "Creates various kinds of things. See `bytemark create <kind of thing> help`",
		Description: `	    bytemark create disc[s] [--disc <disc spec>]... <cloud server>
	create group [--account <name>] <name>
	create ip [--reason reason] <cloud server>
	create server (see bytemark create server help)

A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs`,
		Subcommands: []cli.Command{
			createServer,
			createDiscs,
			createGroup,
		},
	})
}

func fn_createDisc(c *Context) (err error) {
	//flags.Var(&discs, "disc", "")
	discs := c.Discs("disc")

	for i := range discs {
		d, err := discs[i].Validate()
		if err != nil {
			return err
		}
		discs[i] = *d
	}

	log.Logf("Adding discs to %s:\r\n", c.VirtualMachineName)
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

func fn_createGroup(c *Context) (err error) {
	err = global.Client.CreateGroup(c.GroupName)
	if err == nil {
		log.Logf("Group %s was created under account %s\r\n", c.GroupName.Group, c.GroupName.Account)
	}
	err = err
	return

}

func fn_createServer(c *Context) (err error) {
	noImage := c.Bool("no-image")
	if c.Bool("no-discs") {
		noImage = true
	}

	discs := c.Discs("disc")

	if len(discs) == 0 && !c.Context.Bool("no-discs") {
		discs = append(discs, lib.Disc{Size: 25600})
	}

	for i := range discs {
		d, err := discs[i].Validate()
		if err != nil {
			return err
		}
		discs[i] = *d
	}

	ips := c.IPs("ip")

	if len(ips) > 2 {
		log.Log("A maximum of one IPv4 and one IPv6 address may be specified")
		err = &util.PEBKACError{}
		return
	}

	var ipspec *lib.IPSpec
	if len(ips) > 0 {
		ipspec = &lib.IPSpec{}

		for _, ip := range ips {
			if ip.To4() != nil {
				if ipspec.IPv4 != "" {
					log.Log("A maximum of one IPv4 and one IPv6 address may be specified")
					err = &util.PEBKACError{}
					return
				}
				ipspec.IPv4 = ip.To4().String()
			} else {
				if ipspec.IPv6 != "" {
					log.Log("A maximum of one IPv4 and one IPv6 address may be specified")
					err = &util.PEBKACError{}
					return

				}
				ipspec.IPv6 = ip.String()
			}
		}
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

	spec := lib.VirtualMachineSpec{
		VirtualMachine: &lib.VirtualMachine{
			Name:                  c.VirtualMachineName.VirtualMachine,
			Autoreboot:            autoreboot,
			Cores:                 c.Int("cores"),
			Memory:                c.Size("memory"),
			ZoneName:              c.String("zone"),
			CdromURL:              c.String("cdrom"),
			HardwareProfile:       c.String("hwprofile"),
			HardwareProfileLocked: c.Bool("hwprofilelock"),
		},
		Discs:   discs,
		IPs:     ipspec,
		Reimage: imageInstall,
	}

	groupName := c.VirtualMachineName.GroupName()

	log.Log("The following server will be created:")
	log.Log(util.FormatVirtualMachineSpec(groupName, &spec))

	// If we're not forcing, prompt. If the prompt comes back false, exit.
	if !global.Config.Force() && !util.PromptYesNo("Are you certain you wish to continue?") {
		log.Error("Exiting.")
		err = &util.UserRequestedExit{}
		return
	}

	_, err = global.Client.CreateVirtualMachine(groupName, spec)
	if err != nil {
		return
	}
	vm, err := global.Client.GetVirtualMachine(c.VirtualMachineName)
	if err != nil {
		return
	}
	log.Log("cloud server created successfully", "")
	log.Log(util.FormatVirtualMachine(vm))
	if imageInstall != nil {
		log.Log()
		log.Logf("Root password:") // logf so we don't get a tailing \r\n
		log.Outputf("%s\r\n", imageInstall.RootPassword)
	} else {
		log.Log("Machine was not imaged")
	}
	return
}
