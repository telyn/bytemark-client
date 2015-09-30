package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
	"encoding/json"
)

// HelpForShow outputs usage information for the show commands: show, show-vm, show-group, show-account.
func (cmds *CommandSet) HelpForShow() util.ExitCode {
	log.Log("bytemark show")
	log.Log()
	log.Log("usage: bytemark show [--json] <name>")
	log.Log("       bytemark show vm [--json] <virtual machine>")
	log.Log("       bytemark show group [--json] [--list-vms] [--verbose] <group>")
	log.Log("       bytemark show account [--json] [--list-groups] [--list-vms] [--verbose] <account>")
	log.Log()
	log.Log("Displays information about the given virtual machine, group, or account.")
	log.Log("If the --verbose flag is given to bytemark show group or bytemark show account, full details are given for each VM.")
	log.Log()
	return util.E_USAGE_DISPLAYED
}

// ShowVM implements the show-vm command, which is used to display information about BigV VMs. See HelpForShow for the usage information.
func (cmds *CommandSet) ShowVM(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForShow()
		return util.E_PEBKAC
	}
	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)
	if err != nil {
		log.Error("Virtual machine name cannnot be blank")
		return util.E_PEBKAC
	}

	cmds.EnsureAuth()
	vm, err := cmds.bigv.GetVirtualMachine(name)

	if err != nil {
		return util.ProcessError(err)
	}
	if !cmds.config.Silent() {
		if *jsonOut {
			js, _ := json.MarshalIndent(vm, "", "    ")
			log.Output(string(js))
		} else {
			log.Log(util.FormatVirtualMachine(vm))
		}
	}
	return util.E_SUCCESS

}

// ShowGroup implements the show-group command, which is used to show the BigV group name and ID, as well as the VMs within it.
func (cmds *CommandSet) ShowGroup(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	list := flags.Bool("list-vms", false, "")
	verbose := flags.Bool("verbose", false, "")
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "group")
	if !ok {
		cmds.HelpForShow()
		return util.E_PEBKAC
	}
	name := cmds.bigv.ParseGroupName(nameStr)

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	group, err := cmds.bigv.GetGroup(name)

	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {

		if *jsonOut {
			js, _ := json.MarshalIndent(group, "", "    ")
			log.Output(string(js))
		} else {
			log.Outputf("Group %d: %s\r\n\r\n", group.ID, group.Name)

			if *list {
				for _, vm := range group.VirtualMachines {
					log.Output(vm.Name)
				}
			} else if *verbose {
				for _, v := range util.FormatVirtualMachines(group.VirtualMachines) {
					log.Output(v)
				}

			}
		}
	}
	return util.E_SUCCESS

}

// ShowAccount implements the show-account command, which is used to show the BigV account name, as well as the groups and VMs within it.
func (cmds *CommandSet) ShowAccount(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	listgroups := flags.Bool("list-groups", false, "")
	listvms := flags.Bool("list-vms", false, "")
	verbose := flags.Bool("verbose", false, "")
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "account")
	if !ok {
		cmds.HelpForShow()
		return util.E_PEBKAC
	}
	name := cmds.bigv.ParseAccountName(nameStr)

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	acc, err := cmds.bigv.GetAccount(name)

	if err != nil {
		return util.ProcessError(err)
	}

	if *jsonOut {
		js, _ := json.MarshalIndent(acc, "", "    ")
		log.Output(string(js))
	} else {
		log.Outputf("Account %d: %s\r\n", acc.ID, acc.Name)
		switch {
		case *verbose:
			for _, g := range acc.Groups {
				log.Outputf("Group %s\r\n", g.Name)
				for _, v := range util.FormatVirtualMachines(g.VirtualMachines) {
					log.Output(v)
				}
			}
		case *listgroups:
			log.Output("Groups:")
			for _, g := range acc.Groups {
				log.Output(g.Name)
			}
		case *listvms:
			log.Output("Virtual machines:")
			for _, g := range acc.Groups {
				for _, vm := range g.VirtualMachines {
					log.Outputf("%s.%s\r\n", vm.Name, g.Name)
				}
			}
		default:
			vms := 0
			for _, g := range acc.Groups {
				vms += len(g.VirtualMachines)
			}
			log.Outputf("%d groups containing %d virtual machines\r\n", len(acc.Groups), vms)
		}

	}
	return util.E_SUCCESS

}
