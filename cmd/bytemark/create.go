package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
	"strconv"
	"strings"
)

func init() {
	createServer := cli.Command{
		Name:      "server",
		Usage:     "bytemark create server [flags] <name> [<cores> [<memory [<disc specs>]...]]",
		UsageText: `Create a server with bytemark.`,
		Description: `Creates a Cloud Server with the given specification, defaulting to a basic server with Symbiosis installed.
flags available
    --account <name>
    --cores <num> (default 1)
    --cdrom <url>
    --disc <disc spec> - defaults to a single 25GiB sata-grade discs
    --firstboot-script-file <file name> - a script that will be run on first boot
    --force - disables the confirmation prompt
    --group <name>
    --hwprofile <profile>
    --hwprofile-locked
    --image <image name> - specify what to image the server with. Default is 'symbiosis'
    --ip <ip address> (v4 or v6) - up to one of each type may be specified
    --memory <size> (default 1, units are GiB)
    --no-image - specifies that the created server should not be imaged.
    --no-discs - specifies that the created server should not have any discs.
    --public-keys <keys> (newline seperated)
    --public-keys-file <file> (will be read & appended to --public-keys)
    --root-password <password> - if not set, will be randomly generated
    --stopped - if set, machine won't boot
    --zone <name> (default manchester)

A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs

If hwprofile-locked is set then the cloud server's virtual hardware won't be changed over time.`,
		Action: With(VirtualMachineNameProvider, AuthProvider, fn_createServer),
	}

	createDiscs := cli.Command{
		Name:    "discs",
		Aliases: []string{"disc", "disk", "disks"},
		Usage:   "bytemark create discs [--disc <disc spec>]... <cloud server>",
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
	flags := util.MakeCommonFlagSet()
	var discs util.DiscSpecFlag
	flags.Var(&discs, "disc", "")
	global.Config.ImportFlags(flags)

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
	flags := util.MakeCommonFlagSet()
	global.Config.ImportFlags(flags)

	err = global.Client.CreateGroup(c.GroupName)
	if err == nil {
		log.Logf("Group %s was created under account %s\r\n", c.GroupName.Group, c.GroupName.Account)
	}
	err = err
	return

}

func fn_createServer(c *Context) (err error) {
	flags := util.MakeCommonFlagSet()
	addImageInstallFlags(flags)
	cores := flags.Int("cores", 1, "")
	cdrom := flags.String("cdrom", "", "")
	var discs util.DiscSpecFlag
	flags.Var(&discs, "disc", "")
	hwprofile := flags.String("hwprofile", "", "")
	hwprofilelock := flags.Bool("hwprofile-locked", false, "")
	var ips util.IPFlag
	flags.Var(&ips, "ip", "")
	memorySpec := flags.String("memory", "1", "")
	noDiscs := flags.Bool("no-discs", false, "")
	noImage := flags.Bool("no-image", false, "")
	stopped := flags.Bool("stopped", false, "")
	zone := flags.String("zone", "", "")
	args := global.Config.ImportFlags(flags)

	for i, arg := range args {
		switch i {
		// the first arg is the vm name which we already have
		case 0:
			continue
		case 1:
			cores64, err := strconv.ParseInt(arg, 10, 32)
			if err != nil {
				log.Error("Cores argument given was not an int.")
				err = util.PEBKACError{}
				return err
			} else {
				*cores = int(cores64)
			}
		case 2:
			*memorySpec = arg
		default:
			if len(discs) != 0 {
				log.Error("--disc flag used along with the discs spec argument - please use only one")
				err = util.PEBKACError{}
				return
			}
			for i, spec := range strings.Split(arg, ",") {
				disc, err := util.ParseDiscSpec(spec)
				if err != nil {
					log.Errorf("Disc %d has a malformed spec - '%s' is invalid", i, spec)
					//cmds.HelpForTopic('specs')
					err = util.PEBKACError{}
					return err
				}
				discs = append(discs, *disc)
			}

		}
	}

	memory, err := util.ParseSize(*memorySpec)
	if err != nil {
		return
	}

	if *noDiscs {
		*noImage = true
	}

	if len(discs) == 0 && !*noDiscs {
		discs = append(discs, lib.Disc{Size: 25600})
	}

	for i := range discs {
		d, err := discs[i].Validate()
		if err != nil {
			return err
		}
		discs[i] = *d
	}

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

	imageInstall, _, err := prepareImageInstall(flags)
	if err != nil {
		return
	}

	if *noImage {
		imageInstall = nil
	}

	// if stopped isn't set and either cdrom or image are set, start the server
	autoreboot := !*stopped && ((imageInstall != nil) || (*cdrom != ""))

	spec := lib.VirtualMachineSpec{
		VirtualMachine: &lib.VirtualMachine{
			Name:                  c.VirtualMachineName.VirtualMachine,
			Autoreboot:            autoreboot,
			Cores:                 *cores,
			Memory:                memory,
			ZoneName:              *zone,
			CdromURL:              *cdrom,
			HardwareProfile:       *hwprofile,
			HardwareProfileLocked: *hwprofilelock,
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
