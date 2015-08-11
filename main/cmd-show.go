package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// HelpForShow outputs usage information for the show commands: show, show-vm, show-group, show-account.
func (cmds *CommandSet) HelpForShow() {
	fmt.Println("go-bigv show")
	fmt.Println()
	fmt.Println("usage: go-bigv show [--json] <name>")
	fmt.Println("       go-bigv show vm [--json] <virtual machine>")
	fmt.Println("       go-bigv show group [--json] [--list-vms] [--verbose] <group>")
	fmt.Println("       go-bigv show account [--json] [--list-groups] [--list-vms] [--verbose] <account>")
	fmt.Println()
	fmt.Println("Displays information about the given virtual machine, group, or account.")
	fmt.Println("If the --verbose flag is given to bigv show group or bigv show account, full details are given for each VM.")
	fmt.Println()
}

// Show implements the show command which is a stupendous badass of a command
func (cmds *CommandSet) Show(args []string) ExitCode {
	if len(args) == 0 {
		cmds.HelpForShow()
		return E_USAGE_DISPLAYED
	}

	switch strings.ToLower(args[0]) {
	case "vm":
		return cmds.ShowVM(args[1:])
	case "account":
		return cmds.ShowAccount(args[1:])
	case "user":
		fmt.Printf("Leave me alone! I'm grumpy.")
		return 666
		//return ShowUser(args[1:])
	case "group":
		return cmds.ShowGroup(args[1:])
	case "key", "keys":
		fmt.Printf("Leave me alone, I'm grumpy!")
		return 666
		//return cmds.ShowKeys(args[1:])
	}

	name := strings.TrimSuffix(args[0], cmds.config.EndpointName())
	dots := strings.Count(name, ".")
	switch dots {
	case 2:
		return cmds.ShowVM(args)
	case 1:
		return cmds.ShowGroup(args)
	case 0:
		return cmds.ShowAccount(args)
		// TODO: should also try show-vm sprintf("%s.%s.%s", args[0], "default", config.get("user"))
	}
	return E_SUCCESS
}

// ShowVM implements the show-vm command, which is used to display information about BigV VMs. See HelpForShow for the usage information.
func (cmds *CommandSet) ShowVM(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	cmds.EnsureAuth()

	name := cmds.bigv.ParseVirtualMachineName(args[0])

	vm, err := cmds.bigv.GetVirtualMachine(name)

	if err != nil {
		return processError(err)
	}
	if !cmds.config.Silent() {
		if *jsonOut {
			js, _ := json.MarshalIndent(vm, "", "    ")
			fmt.Println(string(js))
		} else {
			fmt.Println(FormatVirtualMachine(vm))
		}
	}
	return E_SUCCESS

}

// ShowGroup implements the show-group command, which is used to show the BigV group name and ID, as well as the VMs within it.
func (cmds *CommandSet) ShowGroup(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	list := flags.Bool("list-vms", false, "")
	verbose := flags.Bool("verbose", false, "")
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name := cmds.bigv.ParseGroupName(args[0])

	cmds.EnsureAuth()

	group, err := cmds.bigv.GetGroup(name)

	if err != nil {
		return processError(err)
	}

	if !cmds.config.Silent() {

		if *jsonOut {
			js, _ := json.MarshalIndent(group, "", "    ")
			fmt.Println(string(js))
		} else {
			fmt.Printf("Group %d: %s\r\n", group.ID, group.Name)
			fmt.Println()
			if *list {
				for _, vm := range group.VirtualMachines {
					fmt.Println(vm.Name)
				}
			} else if *verbose {
				fmt.Println(strings.Join(FormatVirtualMachines(group.VirtualMachines), "\r\n"))

			}
		}
	}
	return E_SUCCESS

}

// ShowAccount implements the show-account command, which is used to show the BigV account name, as well as the groups and VMs within it.
func (cmds *CommandSet) ShowAccount(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	listgroups := flags.Bool("list-groups", false, "")
	listvms := flags.Bool("list-vms", false, "")
	verbose := flags.Bool("verbose", false, "")
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name := cmds.bigv.ParseAccountName(args[0])

	cmds.EnsureAuth()

	acc, err := cmds.bigv.GetAccount(name)

	if err != nil {
		return processError(err)
	}

	if *jsonOut {
		js, _ := json.MarshalIndent(acc, "", "    ")
		fmt.Println(string(js))
	} else {
		fmt.Printf("Account %d: %s", acc.ID, acc.Name)
		fmt.Println()
		switch {
		case *verbose:
			for _, g := range acc.Groups {
				fmt.Printf("Group %s\r\n", g.Name)
				fmt.Println(strings.Join(FormatVirtualMachines(g.VirtualMachines), "\r\n"))
			}
		case *listgroups:
			fmt.Println("Groups:")
			for _, g := range acc.Groups {
				fmt.Println("%s", g.Name)
			}
		case *listvms:
			fmt.Println("Virtual machines:")
			for _, g := range acc.Groups {
				for _, vm := range g.VirtualMachines {
					fmt.Println("%s.%s\r\n", vm.Name, g.Name)
				}
			}
		default:
			vms := 0
			for _, g := range acc.Groups {
				vms += len(g.VirtualMachines)
			}
			fmt.Println("%d groups containing %d virtual machines\r\n", len(acc.Groups), vms)
		}

	}
	return E_SUCCESS

}
