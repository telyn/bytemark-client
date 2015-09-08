package cmds

import (
	"bigv.io/client/cmds/util"
	bigv "bigv.io/client/lib"
	"bigv.io/client/util/log"
	"strings"
)

//HelpForCreateVM provides usage information for the create-vm command
func (cmds *CommandSet) HelpForCreateVM() util.ExitCode {
	log.Log("go-bigv create vm")
	log.Log()
	log.Log("usage: go-bigv create vm [flags] <name> [<cores> [<memory> [<disc specs>]]")
	log.Log()
	log.Log("flags available")
	log.Log("    --account <name>")
	log.Log("    --cores <num> (default 1)")
	log.Log("    --cdrom <url>")
	log.Log("    --discs <disc specs> (default 25)")
	log.Log("    --force")
	log.Log("    --group <name>")
	log.Log("    --hwprofile <profile>")
	log.Log("    --hwprofile-locked")
	log.Log("    --image <image name> (see go-bigv images)")
	log.Log("    --memory <size> (default 1, units are GiB)")
	log.Log("    --public-keys <keys> (newline seperated)")
	log.Log("    --public-keys-file <file> (will be read & appended to --public-keys)")
	log.Log("    --root-password <password>")
	log.Log("    --stopped (if set, machine won't boot)")
	log.Log("    --zone <name> (default manchester)")
	log.Log()
	log.Log("If hwprofile-locked is set then the virtual machine's hardware won't be changed over time.")
	return util.E_USAGE_DISPLAYED

}

//HelpForCreate provides usage information for the create command and its subcommands.
func (cmds *CommandSet) HelpForCreate() util.ExitCode {
	log.Log("go-bigv create")
	log.Log()
	log.Log("usage: go-bigv create disc [--account <name>] [--group <group>] [--size <size>] [--grade <storage grade>] <virtual machine> [disc specs]")
	log.Log("               create group [--account <name>] <name>")
	log.Log("               create disc[s] <disc specs> <virtual machine>")
	log.Log("               create ip [--reason reason] <virtual machine>")
	log.Log("               create vm (see go-bigv help create vm)")
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

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		return util.ProcessError(err)
	}
	var discs []bigv.Disc
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
			discs = append(discs, bigv.Disc{Size: size, StorageGrade: *gradeFlag, Label: *labelFlag})
		}

	} else {
		// if neither of flags and spec are specified, fail
		if len(args) == 0 {
			return cmds.HelpForCreate()
		} else {
			spec := strings.Join(args, " ")

			discs, err = util.ParseDiscSpec(spec, false)
			if err != nil {
				return util.ProcessError(err)
			}
		}

	}
	cmds.EnsureAuth()

	log.Logf("Adding discs to %s:\r\n", name)
	for _, d := range discs {
		log.Logf("    %d %s...", d.Size/1024, d.StorageGrade)
		err = cmds.bigv.CreateDisc(name, d)
		if err != nil {
			log.Errorf("Failure! %v\r\n", err.Error())
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
	name := cmds.bigv.ParseGroupName(nameStr)

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.bigv.CreateGroup(name)
	if err == nil {
		log.Logf("Group %s was created under account %s\r\n", name.Group, name.Account)
	}
	return util.ProcessError(err)

}

// CreateVM implements the create-vm command. See HelpForCreateVM for usage
func (cmds *CommandSet) CreateVM(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	cores := flags.Int("cores", 1, "")
	cdrom := flags.String("cdrom", "", "")
	discSpecs := flags.String("discs", "25", "")
	hwprofile := flags.String("hwprofile", "", "")
	hwprofilelock := flags.Bool("hwprofile-locked", false, "")
	image := flags.String("image", "", "")
	memorySpec := flags.String("memory", "1", "")
	pubkeys := flags.String("public-keys", "", "")
	// pubkeysfile := flags.String("public-keys-file", "", "") // TODO(telyn): --public-keys-file
	rootPassword := flags.String("root-password", "", "")
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

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		return util.ProcessError(err)
	}
	memory, err := util.ParseSize(*memorySpec)
	if err != nil {
		return util.ProcessError(err)
	}

	discs, err := util.ParseDiscSpec(*discSpecs, false)
	if err != nil {
		return util.ProcessError(err)
	}
	for i, d := range discs {
		if d.StorageGrade == "" {
			discs[i].StorageGrade = "sata"
		}
	}

	// if stopped isn't set and either cdrom or image are set, start the vm
	autoreboot := !*stopped && ((*image != "") || (*cdrom != ""))

	spec := bigv.VirtualMachineSpec{
		VirtualMachine: &bigv.VirtualMachine{
			Name:                  name.VirtualMachine,
			Autoreboot:            autoreboot,
			Cores:                 *cores,
			Memory:                memory,
			ZoneName:              *zone,
			CdromURL:              *cdrom,
			HardwareProfile:       *hwprofile,
			HardwareProfileLocked: *hwprofilelock,
		},
		Discs: discs,
		Reimage: &bigv.ImageInstall{
			Distribution: *image,
			PublicKeys:   *pubkeys,
			RootPassword: *rootPassword,
		},
	}

	groupName := bigv.GroupName{
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

	vm, err := cmds.bigv.CreateVirtualMachine(groupName, spec)
	if err != nil {
		return util.ProcessError(err)
	}
	log.Logf("Virtual machine %s created successfully\r\n", vm.Name)
	log.Log(util.FormatVirtualMachine(vm))
	return util.E_SUCCESS

}
