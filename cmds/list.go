package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"strings"
)

func (cmds *CommandSet) HelpForList() util.ExitCode {
	log.Log("bytemark list")
	log.Log("")
	log.Log("usage: bytemark list vms [group | account]")
	log.Log("       bytemark list groups [account]")
	log.Log("       bytemark list accounts")
	log.Log("       bytemark list keys [user]")
	log.Log("       bytemark list discs <virtual machine>")
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) listDefaultAccountVMs() util.ExitCode {
	acc, err := cmds.bigv.GetAccount(cmds.config.GetIgnoreErr("account"))
	if err != nil {
		return util.ProcessError(err)
	}
	for _, group := range acc.Groups {
		for _, vm := range group.VirtualMachines {
			log.Output(vm.Hostname)
		}
	}
	return util.ProcessError(nil)
}

func (cmds *CommandSet) ListVMs(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr := ""
	if len(args) >= 1 {
		nameStr = args[0]
	}
	name := cmds.bigv.ParseGroupName(nameStr, cmds.config.GetGroup())

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if len(args) >= 1 {
		group, err := cmds.bigv.GetGroup(name)
		log.Debugf(5, "Error! %T: %v\r\n", err, err)

		if err != nil {
			if _, ok := err.(lib.NotFoundError); ok {

				if !strings.Contains(nameStr, ".") {
					account, err := cmds.bigv.GetAccount(nameStr)
					if err != nil {
						return util.ProcessError(err)
					}

					for _, g := range account.Groups {
						for _, vm := range g.VirtualMachines {
							log.Output(vm.Hostname)

						}
					}
					return util.E_SUCCESS
				}

			}
			return util.ProcessError(err)
		}

		for _, vm := range group.VirtualMachines {
			log.Output(vm.Hostname)
		}
	} else {
		return cmds.listDefaultAccountVMs()
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) ListGroups(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name := cmds.config.GetIgnoreErr("account")
	if len(args) >= 1 {
		name = args[0]
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
		log.Output(group.Name)
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) ListKeys(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	username := cmds.config.GetIgnoreErr("user")
	if len(args) == 1 {

		usr, ok := util.ShiftArgument(&args, "username")
		if !ok {
			cmds.HelpForShow()
			return util.E_PEBKAC
		}
		username = usr
	}

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	user, err := cmds.bigv.GetUser(username)
	if err != nil {
		return util.ProcessError(err)
	}

	for _, k := range user.AuthorizedKeys {
		log.Output(k)
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

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())

	vm, err := cmds.bigv.GetVirtualMachine(name)

	for _, disc := range vm.Discs {
		log.Outputf("%s: %dGiB %s\r\n", disc.Label, (disc.Size / 1024), disc.StorageGrade)
	}
	return util.E_SUCCESS
}
