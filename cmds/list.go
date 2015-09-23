package cmds

import (
	"bigv.io/client/cmds/util"
	"bigv.io/client/util/log"
)

func (cmds *CommandSet) HelpForList() util.ExitCode {
	log.Log("bytemark list")
	log.Log("")
	log.Log("usage: bytemark list vms <group>")
	log.Log("       bytemark list groups <account>")
	log.Log("       bytemark list accounts")
	log.Log("       bytemark list discs <virtual machine>")
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) ListVMs(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "group")
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

	for _, vm := range group.VirtualMachines {
		log.Log(vm.Hostname)
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) ListGroups(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name, ok := util.ShiftArgument(&args, "account")
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

	for _, group := range account.Groups {
		log.Output(group.Name)
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

	for _, group := range accounts {
		log.Log(group.Name)
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) ListDiscs(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForList()
		return util.E_PEBKAC
	}

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr)

	vm, err := cmds.bigv.GetVirtualMachine(name)

	for _, disc := range vm.Discs {
		log.Logf("%s: %dGiB %s\r\n", disc.Label, (disc.Size / 1024), disc.StorageGrade)
	}
	return util.E_SUCCESS
}
