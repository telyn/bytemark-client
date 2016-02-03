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
	log.Log("usage: bytemark list servers [group | account]")
	log.Log("       bytemark list groups [account]")
	log.Log("       bytemark list accounts")
	log.Log("       bytemark list keys [user]")
	log.Log("       bytemark list discs <server>")
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) listDefaultAccountServers() util.ExitCode {
	acc, err := cmds.client.GetAccount(cmds.config.GetIgnoreErr("account"))
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

func (cmds *CommandSet) ListServers(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr := ""
	if len(args) >= 1 {
		nameStr = args[0]
	}
	name := cmds.client.ParseGroupName(nameStr, cmds.config.GetGroup())

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if len(args) >= 1 {
		group, err := cmds.client.GetGroup(name)

		if err != nil {
			if _, ok := err.(lib.NotFoundError); ok {

				if !strings.Contains(nameStr, ".") {
					account, err := cmds.client.GetAccount(nameStr)
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
		return cmds.listDefaultAccountServers()
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

	account, err := cmds.client.GetAccount(name)

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

	accounts, err := cmds.client.GetAccounts()

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

	user, err := cmds.client.GetUser(username)
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

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForList()
		return util.E_PEBKAC
	}

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())

	vm, err := cmds.client.GetVirtualMachine(name)

	for _, disc := range vm.Discs {
		log.Outputf("%s: %dGiB %s\r\n", disc.Label, (disc.Size / 1024), disc.StorageGrade)
	}
	return util.E_SUCCESS
}
