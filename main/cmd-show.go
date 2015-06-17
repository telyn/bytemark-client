package main

import (
	"encoding/json"
	"fmt"
)

// HelpForShow outputs usage information for the show commands: show, show-vm, show-group, show-account.
func (cmds *CommandSet) HelpForShow() {
	fmt.Println("go-bigv show")
	fmt.Println()
	fmt.Println("usage: go-bigv show [--json] <name>")
	fmt.Println("       go-bigv show vm [--json] <virtual machine>")
	fmt.Println("       go-bigv show group [--json] [--verbose] <group>")
	fmt.Println("       go-bigv show account [--json] [--verbose] <account>")
	fmt.Println()
	fmt.Println("Displays information about the given virtual machine, group, or account.")
	fmt.Println("If the --verbose flag is given to bigv show group or bigv show account, full details are given for each VM.")
	fmt.Println()
}

// ShowVM implements the show-vm command, which is used to display information about BigV VMs. See HelpForShow for the usage information.
func (cmds *CommandSet) ShowVM(args []string) {
	flags := MakeCommonFlagSet()
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	vm, err := cmds.bigv.GetVirtualMachine(name)

	if err != nil {
		exit(err)
	}
	if !cmds.config.GetBool("silent") {
		if *jsonOut {
			fmt.Println(json.MarshalIndent(vm, "", "    "))
		} else {
			fmt.Println(FormatVirtualMachine(vm))
		}
	}

}

// ShowAccount implements the show-account command, which is used to show the BigV account name, as well as the groups and VMs within it.
func (cmds *CommandSet) ShowAccount(args []string) {
	flags := MakeCommonFlagSet()
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name := cmds.bigv.ParseAccountName(args[0])

	acc, err := cmds.bigv.GetAccount(name)

	if err != nil {
		exit(err)
	}

	if *jsonOut {
		fmt.Println(json.MarshalIndent(acc, "", "    "))
	} else {
		fmt.Printf("Account %d: %s", acc.ID, acc.Name)
	}

}
