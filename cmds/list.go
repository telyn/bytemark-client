package cmds

import (
	"bigv.io/client/cmds/util"
	"fmt"
)

func (cmds *CommandSet) HelpForList() util.ExitCode {
	fmt.Println("bigv list")
	fmt.Println("")
	fmt.Println("usage: bigv list vms")
	fmt.Println("       bigv list groups")
	fmt.Println("       bigv list accounts")
	fmt.Println("       bigv list discs")
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) ListVMs(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(args, "group")
	if !ok {
		cmds.HelpForList()
		return util.E_PEBKAC
	}
	name := cmds.bigv.ParseGroupName(nameStr)

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	group, err := cmds.bigv.GetGroup(name)

	if err != nil {
		// TODO: try it as an account
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		for _, vm := range group.VirtualMachines {
			fmt.Println(vm.Hostname)
		}
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) ListGroups(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, ok := util.ShiftArgument(args, "account")
	if !ok {
		cmds.HelpForList()
		return util.E_PEBKAC
	}

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	account, err := cmds.bigv.GetAccount(name)

	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		for _, group := range account.Groups {
			fmt.Println(group.Name)
		}
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) ListAccounts(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	accounts, err := cmds.bigv.GetAccounts()

	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {
		for _, group := range accounts {
			fmt.Println(group.Name)
		}
	}
	return util.E_SUCCESS
}
