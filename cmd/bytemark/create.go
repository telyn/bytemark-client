package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"os"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "create",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "user",
				Usage:     "creates a new cluster admin or cluster superuser",
				UsageText: "bytemark --admin create user <username> <privilege>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "username",
						Usage: "The username of the new user",
					},
					cli.StringFlag{
						Name:  "privilege",
						Usage: "The privilege to grant to the new user",
					},
				},
				Action: With(OptionalArgs("username", "privilege"), RequiredFlags("username", "privilege"), AuthProvider, func(c *Context) error {
					// Privilege is just a string and not a PrivilegeFlag, since it can only be "cluster_admin" or "cluster_su"
					if err := global.Client.CreateUser(c.String("username"), c.String("privilege")); err != nil {
						return err
					}
					log.Logf("User %s has been created with %s privileges\r\n", c.String("username"), c.String("privilege"))
					return nil
				}),
			},
			{
				Name:      "vlan_group",
				Usage:     "creates groups for private VLANs",
				UsageText: "bytemark --admin create vlan_group <group> [vlan_num]",
				Description: `Create a group in the specified account, with an optional VLAN specified.

Used when setting up a private VLAN for a customer.`,
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "group",
						Usage: "the name of the group to create",
						Value: new(GroupNameFlag),
					},
					cli.IntFlag{
						Name:  "vlan_num",
						Usage: "The VLAN number to add the group to",
					},
				},
				Action: With(OptionalArgs("group", "vlan_num"), RequiredFlags("group"), AuthProvider, func(c *Context) error {
					gp := c.GroupName("group")
					if err := global.Client.AdminCreateGroup(gp, c.Int("vlan_num")); err != nil {
						return err
					}
					log.Logf("Group %s was created under account %s\r\n", gp.Group, gp.Account)
					return nil
				}),
			},
			{
				Name:      "ip_range",
				Usage:     "create a new IP range in a VLAN",
				UsageText: "bytemark --admin create ip_range <ip_range> <vlan_num>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "ip_range",
						Usage: "the IP range to add",
					},
					cli.IntFlag{
						Name:  "vlan_num",
						Usage: "The VLAN number to add the IP range to",
					},
				},
				Action: With(OptionalArgs("ip_range", "vlan_num"), RequiredFlags("ip_range", "vlan_num"), AuthProvider, func(c *Context) error {
					if err := global.Client.CreateIPRange(c.String("ip_range"), c.Int("vlan_num")); err != nil {
						return err
					}
					log.Logf("IP range created\r\n")
					return nil
				}),
			},
		},
	})

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
		Flags: append(OutputFlags("server", "object"),
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
			cli.GenericFlag{
				Name:  "memory",
				Value: new(util.SizeSpecFlag),
				Usage: "How much memory the server will have available, specified in GiB or with GiB/MiB units. Defaults to 1GiB.",
			},
			cli.GenericFlag{
				Name:  "name",
				Usage: "The new server's name",
				Value: new(VirtualMachineNameFlag),
			},
			cli.BoolFlag{
				Name:  "no-image",
				Usage: "Specifies that the server should not be imaged.",
			},
			cli.BoolFlag{
				Name:  "stopped",
				Usage: "If set, the server will not be started, even to image it.",
			},
			cli.StringFlag{
				Name:  "zone",
				Usage: "Which zone the server will be created in. See `bytemark zones` for the choices.",
			},
		),
		Action: With(OptionalArgs("name", "cores", "memory", "disc"), RequiredFlags("name"), AuthProvider, createServer),
	}
	createServerCmd.Flags = append(createServerCmd.Flags, imageInstallFlags...)

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
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to add the disc to",
				Value: new(VirtualMachineNameFlag),
			},
		},
		Usage:     "create virtual discs attached to one of your cloud servers",
		UsageText: "bytemark create discs [--disc <disc spec>]... <cloud server>",
		Description: `A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs`,
		Action: With(OptionalArgs("server", "cores", "memory", "disc"), AuthProvider, createDiscs),
	}

	createGroupCmd := cli.Command{
		Name:        "group",
		Usage:       "create a group for organising your servers",
		UsageText:   "bytemark create group <group name>",
		Description: `Groups are part of your server's fqdn`,
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "group",
				Usage: "the name of the group to create",
				Value: new(GroupNameFlag),
			},
		},
		Action: With(OptionalArgs("group"), RequiredFlags("group"), AuthProvider, createGroup),
	}

	createBackupCmd := cli.Command{
		Name:        "backup",
		Usage:       "create a backup of a disc's current state",
		UsageText:   "bytemark create backup <server name> <disc label>",
		Description: `Creates a backup of the disc's current state. The backup is moved to another tail in the "iceberg" storage grade.`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "disc",
				Usage: "the disc to create a backup of",
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server whose disk you wish to backup",
				Value: new(VirtualMachineNameFlag),
			},
		},
		Action: With(OptionalArgs("server", "disc"), RequiredFlags("server", "disc"), AuthProvider, func(c *Context) error {
			backup, err := global.Client.CreateBackup(c.VirtualMachineName("server"), c.String("disc"))
			if err != nil {
				return err
			}
			log.Errorf("Backup '%s' taken successfully!", backup.Label)
			return nil
		}),
	}

	commands = append(commands, cli.Command{
		Name:      "create",
		Usage:     "creates servers, discs, etc - see `bytemark create <kind of thing> help`",
		UsageText: "bytemark create disc|group|ip|server",
		Description: `create a new disc, group, IP or server

	create disc[s] [--disc <disc spec>]... <cloud server>
	create group [--account <name>] <name>
	create ip [--reason reason] <cloud server>
	create server (see bytemark help create server)

A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			createServerCmd,
			createDiscsCmd,
			createGroupCmd,
			createBackupCmd,
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
	vmName := c.VirtualMachineName("server")

	log.Logf("Adding %d discs to %s:\r\n", len(discs), vmName)
	for _, d := range discs {
		log.Logf("    %dGiB %s...", d.Size/1024, d.StorageGrade)
		err := global.Client.CreateDisc(&vmName, d)
		if err != nil {
			log.Errorf("failure! %v\r\n", err.Error())
		} else {
			log.Log("success!")
		}
	}
	return
}

func createGroup(c *Context) (err error) {
	gp := c.GroupName("group")
	err = global.Client.CreateGroup(&gp)
	if err == nil {
		log.Logf("Group %s was created under account %s\r\n", gp.Group, gp.Account)
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

	discs, cores, memory, err := createServerReadArgs(c)
	if err != nil {
		return
	}

	if len(discs) == 0 {
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
			Name:                  c.VirtualMachineName("name").VirtualMachine,
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
	name := c.VirtualMachineName("name")
	spec, err := createServerPrepSpec(c)
	if err != nil {
		return
	}

	groupName := name.GroupName()

	log.Log("The following server will be created:")
	err = spec.PrettyPrint(os.Stderr, prettyprint.Full)
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
	vm, err := global.Client.GetVirtualMachine(&name)
	if err != nil {
		return
	}
	return c.OutputInDesiredForm(map[string]interface{}{"spec": spec, "virtual_machine": vm}, func() (err error) {
		log.Log("cloud server created successfully")
		err = vm.PrettyPrint(os.Stderr, prettyprint.Full)
		if err != nil {
			return
		}
		if spec.Reimage != nil {
			log.Log()
			log.Logf("Root password: ") // logf so we don't get a trailing \r\n
			log.Outputf("%s\r\n", spec.Reimage.RootPassword)
		} else {
			log.Log("Machine was not imaged")
		}
		return
	})
}
