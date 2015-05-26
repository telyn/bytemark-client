package main

import (
	"fmt"
)

// HelpForShow outputs usage information for the show commands: show, show-vm, show-group, show-account.
func (cmds *CommandSet) HelpForShow() {
	// TODO(telyn): Replace instances of bigv with $0, however you get $0 in go?
	fmt.Println("bigv show")
	fmt.Println()
	fmt.Println("usage: bigv show [-j | --json] <name>")
	fmt.Println("       bigv show vm [-j | --json] <virtual machine>")
	fmt.Println("       bigv show group [-j | --json] [-v | --verbose] <group>")
	fmt.Println("       bigv show account [-j | --json] [-v | --verbose] <account>")
	fmt.Println()
	fmt.Println("Displays information about the given virtual machine, group, or account.")
	fmt.Println("If the --verbose flag is given to bigv show group or bigv show account, full details are given for each VM.")
	fmt.Println()
}

// ShowVM implements the show-vm command, which is used to display information about BigV VMs. See HelpForShow for the usage information.
func (cmds *CommandSet) ShowVM(args []string) {
	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	vm, err := cmds.bigv.GetVirtualMachine(name)

	if err != nil {
		exit(err)
	}
	if cmds.config.Get("silent") != "true" {
		fmt.Println(FormatVirtualMachine(vm))
	}

}

// ShowAccount implements the show-account command, which is used to show the BigV account name, as well as the groups and VMs within it.
func (cmds *CommandSet) ShowAccount(args []string) {
	name := cmds.bigv.ParseAccountName(args[0])

	acc, err := cmds.bigv.GetAccount(name)

	if err != nil {
		exit(err)
	}

	fmt.Printf("Account %d: %s", acc.Id, acc.Name)

}
