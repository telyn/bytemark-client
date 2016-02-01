package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"strconv"
	"strings"
)

//HelpForCreateVM provides usage information for the create-vm command
func (cmds *CommandSet) HelpForCreateVM() util.ExitCode {
	log.Log("bytemark create vm")
	log.Log()
	log.Log("usage: bytemark create vm [flags] <name> [<cores> [<memory> [<disc specs>]...]")
	log.Log()
	log.Log("flags available")
	log.Log("    --account <name>")
	log.Log("    --cores <num> (default 1)")
	log.Log("    --cdrom <url>")
	log.Log("    --disc <disc spec> - defaults to a single 25GiB sata-grade discs")
	log.Log("    --firstboot-script-file <file name> - a script that will be run on first boot")
	log.Log("    --force - disables the confirmation prompt")
	log.Log("    --group <name>")
	log.Log("    --hwprofile <profile>")
	log.Log("    --hwprofile-locked")
	log.Log("    --image <image name> - specify what to image the server with. Default is 'symbiosis'")
	log.Log("    --ip <ip address> (v4 or v6) - up to one of each type may be specified")
	log.Log("    --memory <size> (default 1, units are GiB)")
	log.Log("    --no-image - specifies that the created server should not be imaged.")
	log.Log("    --no-discs - specifies that the created server should not have any discs.")
	log.Log("    --public-keys <keys> (newline seperated)")
	log.Log("    --public-keys-file <file> (will be read & appended to --public-keys)")
	log.Log("    --root-password <password> (if not set, will be randomly generated)")
	log.Log("    --stopped (if set, machine won't boot)")
	log.Log("    --zone <name> (default manchester)")
	log.Log()
	log.Log("A disc spec looks like the following: label:grade:size")
	log.Log("The label and grade fields are optional. If grade is empty, defaults to sata.")
	log.Log("If there are two fields, they are assumed to be grade and size.")
	log.Log("Multiple --disc flags can be used to create multiple discs")
	log.Log()
	log.Log("If hwprofile-locked is set then the virtual machine's hardware won't be changed over time.")
	return util.E_USAGE_DISPLAYED

}

//HelpForCreate provides usage information for the create command and its subcommands.
func (cmds *CommandSet) HelpForCreate() util.ExitCode {
	log.Log("bytemark create")
	log.Log()
	log.Log("usage: bytemark create disc [--account <name>] [--group <group>] [--size <size>] [--grade <storage grade>] <virtual machine> [disc specs]")
	log.Log("               create group [--account <name>] <name>")
	log.Log("               create disc[s] <disc specs> <virtual machine>")
	log.Log("               create ip [--reason reason] <virtual machine>")
	log.Log("               create vm (see bytemark help create vm)")
	log.Log("")
	log.Log("Disc specs are a comma seperated list of size:storage grade pairs. Sizes are in GB by default but can be specified in M")
	log.Log("")
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) CreateDiscs(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	sizeFlag := flags.String("size", "", "")
	gradeFlag := flags.String("grade", "", "")
	labelFlag := flags.String("label", "", "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForCreate()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		return util.ProcessError(err)
	}

	var discs []lib.Disc
	if *sizeFlag != "" || *gradeFlag != "" {
		// if both flags and spec are specified, fail
		if len(args) >= 1 {
			if !cmds.config.Silent() {
				log.Error("Ambiguous command given - please only specify disc specs as arguments or flags, not both")
			}
			return util.E_PEBKAC

		} else {
			// only flags
			size, err := util.ParseSize(*sizeFlag)
			if err != nil {
				return util.ProcessError(err)
			}
			discs = append(discs, lib.Disc{Size: size, StorageGrade: *gradeFlag, Label: *labelFlag})
		}

	} else {
		// if neither of flags and spec are specified, fail
		if len(args) == 0 {
			return cmds.HelpForCreate()
		} else {
			specs := strings.Split(strings.Join(args, " "), ",")
			for _, spec := range specs {

				disc, err := util.ParseDiscSpec(spec)
				if err != nil {
					return util.ProcessError(err)
				}
				discs = append(discs, *disc)
			}
		}

	}
	for i := range discs {
		d, err := discs[i].Validate()
		if err != nil {
			return util.ProcessError(err)
		}
		discs[i] = *d
	}
	cmds.EnsureAuth()

	log.Logf("Adding discs to %s:\r\n", name)
	for _, d := range discs {
		log.Logf("    %dGiB %s...", d.Size/1024, d.StorageGrade)
		err = cmds.client.CreateDisc(name, d)
		if err != nil {
			log.Errorf("failure! %v\r\n", err.Error())
		} else {
			log.Log("success!")
		}
	}
	return util.ProcessError(err)

}

// CreateGroup implements the create-group command. See HelpForCreateGroup for usage.
func (cmds *CommandSet) CreateGroup(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "group")
	if !ok {
		cmds.HelpForCreate()
		return util.E_PEBKAC
	}
	name := cmds.client.ParseGroupName(nameStr, cmds.config.GetGroup())

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.client.CreateGroup(name)
	if err == nil {
		log.Logf("Group %s was created under account %s\r\n", name.Group, name.Account)
	}
	return util.ProcessError(err)

}

// CreateVM implements the create-vm command. See HelpForCreateVM for usage
func (cmds *CommandSet) CreateVM(args []string) util.ExitCode {
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
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	var err error
	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForCreateVM()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		return util.ProcessError(err)
	}
	for i, arg := range args {
		switch i {
		case 0:
			cores64, err := strconv.ParseInt(arg, 10, 32)
			if err != nil {
				log.Error("Cores argument given was not an int.")
				cmds.HelpForCreateVM()
				return util.E_PEBKAC
			} else {
				*cores = int(cores64)
			}
		case 1:
			*memorySpec = arg
		default:
			if len(discs) != 0 {
				log.Error("--disc flag used along with the discs spec argument - please use only one")
				cmds.HelpForCreateVM()
				return util.E_PEBKAC
			}
			for i, spec := range strings.Split(arg, ",") {
				disc, err := util.ParseDiscSpec(spec)
				if err != nil {
					log.Errorf("Disc %d has a malformed spec - '%s' is invalid", i, spec)
					//cmds.HelpForTopic('specs')
					return util.E_PEBKAC
				}
				discs = append(discs, *disc)
			}

		}
	}

	memory, err := util.ParseSize(*memorySpec)
	if err != nil {
		return util.ProcessError(err)
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
			return util.ProcessError(err)
		}
		discs[i] = *d
	}

	if len(ips) > 2 {
		log.Debugf(1, "%d IP addresses were specified", len(ips))
		log.Log("A maximum of one IPv4 and one IPv6 address may be specified")
		return util.E_PEBKAC
	}

	var ipspec *lib.IPSpec
	if len(ips) > 0 {
		ipspec = &lib.IPSpec{}

		for _, ip := range ips {
			if ip.To4() != nil {
				if ipspec.IPv4 != "" {
					log.Debugf(1, "Multiple IPv4 addresses were specified\n")
					log.Log("A maximum of one IPv4 and one IPv6 address may be specified")
					return util.E_PEBKAC
				}
				ipspec.IPv4 = ip.To4().String()
			} else {
				if ipspec.IPv6 != "" {
					log.Debugf(1, "Multiple IPv6 addresses were specified\n")
					log.Log("A maximum of one IPv4 and one IPv6 address may be specified")
					return util.E_PEBKAC

				}
				ipspec.IPv6 = ip.String()
			}
		}
	}

	imageInstall, _, err := prepareImageInstall(flags)
	if err != nil {
		return util.ProcessError(err)
	}

	if *noImage {
		imageInstall = nil
	}

	// if stopped isn't set and either cdrom or image are set, start the vm
	autoreboot := !*stopped && ((imageInstall != nil) || (*cdrom != ""))

	spec := lib.VirtualMachineSpec{
		VirtualMachine: &lib.VirtualMachine{
			Name:                  name.VirtualMachine,
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

	groupName := lib.GroupName{
		Group:   name.Group,
		Account: name.Account,
	}

	log.Log("The following VM will be created:")
	log.Log(util.FormatVirtualMachineSpec(&groupName, &spec))

	// If we're not forcing, prompt. If the prompt comes back false, exit.
	if !cmds.config.Force() && !util.PromptYesNo("Are you certain you wish to continue?") {
		log.Error("Exiting.")
		return util.ProcessError(&util.UserRequestedExit{})
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	_, err = cmds.client.CreateVirtualMachine(groupName, spec)
	if err != nil {
		return util.ProcessError(err)
	}
	vm, err := cmds.client.GetVirtualMachine(name)
	if err != nil {
		return util.ProcessError(err)
	}
	log.Log("Virtual machine created successfully", "")
	log.Log(util.FormatVirtualMachine(vm))
	if imageInstall != nil {
		log.Logf("Root password:") // logf so we don't get a tailing \r\n
		log.Outputf("%s\r\n", imageInstall.RootPassword)
	} else {
		log.Log("Machine was not imaged")
	}
	return util.E_SUCCESS

}
