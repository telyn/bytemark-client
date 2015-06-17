package main

import (
	bigv "bigv.io/client/lib"
	"fmt"
)

//HelpForCreateVM provides usage information for the create-vm command
func (cmd *CommandSet) HelpForCreateVM() {
	fmt.Println("go-bigv create vm")
	fmt.Println()
	fmt.Println("usage: go-bigv create vm [flags] <name>")
	fmt.Println()
	fmt.Println("flags available")
	fmt.Println("    --account <name>")
	fmt.Println("    --cores <num> (default 1)")
	fmt.Println("    --cdrom <url>")
	fmt.Println("    --discs <disc spec> (default 25)")
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

}

//HelpForCreate provides usage information for the create command and its subcommands.
func (cmd *CommandSet) HelpForCreate() {
	fmt.Println("go-bigv create")
	fmt.Println()
	fmt.Println("usage: go-bigv create disc [--account <name>] [--group <group>] [--size <size>] [--grade <storage grade>] <virtual machine name>")
	fmt.Println("               create group [--account <name>] <name>")
	fmt.Println("               create discs <disc specs>")
	fmt.Println("               create ip [--reason reason] <virtual machine>")
	fmt.Println("               create vm (see go-bigv help create vm)")
	fmt.Println("")
}

// CreateGroup implements the create-group command. See HelpForCreateGroup for usage.
func (cmd *CommandSet) CreateGroup(args []string) {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmd.config.ImportFlags(flags)

	// extract this out to PromptForGroupName probably
	name := bigv.GroupName{"", ""}
	if len(args) == 0 {
		name = cmd.bigv.ParseGroupName(Prompt("Group name: "))
	} else if name = cmd.bigv.ParseGroupName(args[0]); name.Group == "" {
		name = cmd.bigv.ParseGroupName(Prompt("Group name: "))
	}

	if name.Account == "" {
		if name.Account = cmd.config.Get("account"); name.Account == "" {
			name.Account = Prompt("Account name: ")
		}
	}

	cmd.EnsureAuth()

	err := cmd.bigv.CreateGroup(name)
	if err == nil {
		fmt.Printf("Group %s was created under account %s\r\n", name.Group, name.Account)
	}
	exit(err)

}

// CreateVM implements the create-vm command. See HelpForCreateVM for usage
func (cmd *CommandSet) CreateVM(args []string) {
	flags := MakeCommonFlagSet()
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
	args = cmd.config.ImportFlags(flags)

	name := bigv.VirtualMachineName{"", "", ""}
	if len(args) > 0 {
		name = cmd.bigv.ParseVirtualMachineName(args[0])
	}

	if *group != "" {
		name.Group = *group
	}

	if name.Group == "" {
		name.Group = "default"
	}

	if name.Account == "" {
		name.Account = cmd.config.Get("account")
	}

	if name.Account == "" {
		name.Account = Prompt("Account: ")
	}

	discs, err := ParseDiscSpec(*discSpecs, false)
	if err != nil {
		exit(err)
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

	if !cmd.config.GetBool("silent") {
		fmt.Println("The following VM will be created:")
		fmt.Println(FormatVirtualMachineSpec(&groupName, &spec))
	}

	// If we're not forcing, prompt. If the prompt comes back false, exit.
	if !cmd.config.GetBool("force") && !PromptYesNo("Are you certain you wish to continue?") {
		fmt.Println("Exiting.")
		exit(&UserRequestedExit{})
	}

	cmd.EnsureAuth()

	vm, err := cmd.bigv.CreateVirtualMachine(groupName, spec)
	if err != nil {
		exit(err)
	} else if !cmd.config.GetBool("silent") {
		fmt.Printf("Virtual machine %s created successfully\n", vm.Name)
		fmt.Println(FormatVirtualMachine(vm))
	}
	exit(nil)

}
