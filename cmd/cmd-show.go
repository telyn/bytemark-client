package main

import (
	"fmt"
)

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

func (cmds *CommandSet) ShowVM(args []string) {
	cmds.EnsureAuth()

	name := ParseVirtualMachineName(args[0])

	vm, err := cmds.bigv.GetVirtualMachine(name)

	if err != nil {
		exit(err)
	}

	fmt.Println(FormatVirtualMachine(vm))

}

func (cmds *CommandSet) ShowAccount(args []string) {
	name := ParseAccountName(args[0])

	acc, err := cmds.bigv.GetAccount(name)

	if err != nil {
		exit(err)
	}

	fmt.Printf("Account %d: %s", acc.Id, acc.Name)

}
