package main

import (
	bigv "bigv.io/client/lib"
	"fmt"
)

func (cmd *CommandSet) CreateGroup(args []string) {
	flags := MakeCommonFlagSet()
	flags.Parse(args)
	args = flags.Args()

	name := bigv.GroupName{"", ""}
	if len(args) == 0 {
		name = ParseGroupName(Prompt("Group name: "))
	} else if name = ParseGroupName(args[0]); name.Group == "" {
		name = ParseGroupName(Prompt("Group name: "))
	}

	if name.Account == "" {
		// TODO(telyn): think I need to de-boilerplate flag lookup too.
		// oh yeah I was going to add Config.ReadFlags(FlagSet)
		val := flags.Lookup("account").Value
		if name.Account = val.String(); name.Account == "" {
			if name.Account = cmd.config.Get("account"); name.Account == "" {
				name.Account = Prompt("Account name: ")
			}
		}
	}

	cmd.EnsureAuth()

	err := cmd.bigv.CreateGroup(name)
	if err == nil {
		fmt.Printf("Group %s was created under account %s\r\n", name.Group, name.Account)
	}
	exit(err)

}
