package cmds

import (
	"bigv.io/client/cmds/util"
	bigv "bigv.io/client/lib"
	"fmt"
)

//HelpForCreateVM provides usage information for the create-vm command
func (cmds *CommandSet) HelpForCreateVM() cmd.ExitCode {
	fmt.Println("go-bigv create vm")
	fmt.Println()
	fmt.Println("usage: go-bigv create vm [flags] <name> [<cores> [<memory> [<disc specs>]]")
	fmt.Println()
	fmt.Println("flags available")
	fmt.Println("    --account <name>")
	fmt.Println("    --cores <num> (default 1)")
	fmt.Println("    --cdrom <url>")
	fmt.Println("    --discs <disc specs> (default 25)")
	fmt.Println("    --force")
	fmt.Println("    --group <name>")
	fmt.Println("    --hwprofile <profile>")
	fmt.Println("    --hwprofile-locked")
	fmt.Println("    --image <image name> (see go-bigv images)")
	fmt.Println("    --memory <num> (default 1, unit GB)")
	fmt.Println("    --public-keys <keys> (newline seperated)")
	fmt.Println("    --public-keys-file <file> (will be read & appended to --public-keys)")
	fmt.Println("    --root-password <password>")
	fmt.Println("    --stopped (if set, machine won't boot)")
	fmt.Println("    --zone <name> (default manchester)")
	fmt.Println()
	fmt.Println("If hwprofile-locked is set then the virtual machine's hardware won't be changed over time.")
	return cmd.E_USAGE_DISPLAYED

}

//HelpForCreate provides usage information for the create command and its subcommands.
func (cmds *CommandSet) HelpForCreate() cmd.ExitCode {
	fmt.Println("go-bigv create")
	fmt.Println()
	fmt.Println("usage: go-bigv create disc [--account <name>] [--group <group>] [--size <size>] [--grade <storage grade>] <virtual machine>")
	fmt.Println("               create group [--account <name>] <name>")
	fmt.Println("               create disc[s] <disc specs> <virtual machine>")
	fmt.Println("               create ip [--reason reason] <virtual machine>")
	fmt.Println("               create vm (see go-bigv help create vm)")
	fmt.Println("")
	return cmd.E_USAGE_DISPLAYED
}

// CreateGroup implements the create-group command. See HelpForCreateGroup for usage.
func (cmds *CommandSet) CreateGroup(args []string) cmd.ExitCode {
	flags := cmd.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	// extract this out to cmd.PromptForGroupName probably
	name := bigv.GroupName{"", ""}
	if len(args) == 0 {
		name = cmds.bigv.ParseGroupName(cmd.Prompt("Group name: "))
	} else if name = cmds.bigv.ParseGroupName(args[0]); name.Group == "" {
		name = cmds.bigv.ParseGroupName(cmd.Prompt("Group name: "))
	}

	var err error
	if name.Account == "" {
		if name.Account, err = cmds.config.Get("account"); name.Account == "" {
			name.Account = cmd.Prompt("Account name: ")
		}
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return cmd.ProcessError(err)
	}

	err = cmds.bigv.CreateGroup(name)
	if err == nil {
		fmt.Printf("Group %s was created under account %s\r\n", name.Group, name.Account)
	}
	return cmd.ProcessError(err)

}

// CreateVM implements the create-vm command. See HelpForCreateVM for usage
func (cmds *CommandSet) CreateVM(args []string) cmd.ExitCode {
	flags := cmd.MakeCommonFlagSet()
	cores := flags.Int("cores", 1, "")
	cdrom := flags.String("cdrom", "", "")
	discSpecs := flags.String("discs", "25", "")
	group := flags.String("group", "", "")
	hwprofile := flags.String("hwprofile", "", "")
	hwprofilelock := flags.Bool("hwprofile-locked", false, "")
	image := flags.String("image", "", "")
	memory := flags.Int("memory", 1, "")
	pubkeys := flags.String("public-keys", "", "")
	// pubkeysfile := flags.String("public-keys-file", "", "") // TODO(telyn): --public-keys-file
	rootPassword := flags.String("root-password", "", "")
	stopped := flags.Bool("stopped", false, "")
	zone := flags.String("zone", "", "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	var err error

	name := bigv.VirtualMachineName{"", "", ""}
	if len(args) > 0 {
		name, err = cmds.bigv.ParseVirtualMachineName(args[0])

	}

	if *group != "" {
		name.Group = *group
	}

	if name.Group == "" {
		name.Group = "default"
	}

	if name.Account == "" {
		name.Account, err = cmds.config.Get("account")
	}

	if name.Account == "" {
		name.Account = cmd.Prompt("Account: ")
	}

	discs, err := cmd.ParseDiscSpec(*discSpecs, false)
	if err != nil {
		return cmd.ProcessError(err)
	}
	for _, d := range discs {
		if d.StorageGrade == "" {
			d.StorageGrade = "sata"
		}
	}

	// if stopped isn't set and either cdrom or image are set, start the vm
	autoreboot := !*stopped && ((*image != "") || (*cdrom != ""))

	spec := bigv.VirtualMachineSpec{
		VirtualMachine: &bigv.VirtualMachine{
			Name:                  name.VirtualMachine,
			Autoreboot:            autoreboot,
			Cores:                 *cores,
			Memory:                *memory,
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

	if !cmds.config.Silent() {
		fmt.Println("The following VM will be created:")
		fmt.Println(cmd.FormatVirtualMachineSpec(&groupName, &spec))
	}

	// If we're not forcing, prompt. If the prompt comes back false, exit.
	if !cmds.config.Force() && !cmd.PromptYesNo("Are you certain you wish to continue?") {
		fmt.Println("Exiting.")
		return cmd.ProcessError(&cmd.UserRequestedExit{})
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return cmd.ProcessError(err)
	}

	vm, err := cmds.bigv.CreateVirtualMachine(groupName, spec)
	if err != nil {
		return cmd.ProcessError(err)
	} else if !cmds.config.Silent() {
		fmt.Printf("Virtual machine %s created successfully\n", vm.Name)
		fmt.Println(cmd.FormatVirtualMachine(vm))
	}
	return cmd.E_SUCCESS

}
