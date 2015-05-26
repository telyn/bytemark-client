package main

import (
	bigv "bigv.io/client/lib"
	"fmt"
)

//TODO(telyn): create-vm should have a --ansible=FILE flag that appends the VM's hostname to the given file, adds

// CreateGroup implements the create-group command. See HelpForCreateGroup for usage.
func (cmd *CommandSet) CreateGroup(args []string) {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = cmd.config.ImportFlags(flags)

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
