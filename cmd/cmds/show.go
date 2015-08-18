package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// HelpForShow outputs usage information for the show commands: show, show-vm, show-group, show-account.
func (cmds *CommandSet) HelpForShow() ExitCode {
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
	return E_USAGE_DISPLAYED
}

// ShowVM implements the show-vm command, which is used to display information about BigV VMs. See HelpForShow for the usage information.
func (cmds *CommandSet) ShowVM(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	err := cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	name, err := cmds.bigv.ParseVirtualMachineName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Virtual machine name cannnot be blank\r\n")
		return E_PEBKAC
	}

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

	err := cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

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

	err := cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

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
