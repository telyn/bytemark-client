package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
	"encoding/json"
)

// HelpForShow outputs usage information for the show commands: show, show server, show group, show account.
func (cmds *CommandSet) HelpForShow() util.ExitCode {
	log.Log("bytemark show")
	log.Log()
	log.Log("usage: bytemark show [--json] <name>")
	log.Log("       bytemark show [--json] <server>")
	log.Log("       bytemark show group [--json] [--verbose] <group>")
	log.Log("       bytemark show account [--json] [--verbose] <account>")
	log.Log()
	log.Log("Displays information about the given server, group, or account.")
	log.Log("If the --verbose flag is given to bytemark show group or bytemark show account, full details are given for each server.")
	log.Log()
	return util.E_USAGE_DISPLAYED
}

// ShowServer implements the show server command, which is used to display information about Bytemark servers. See HelpForShow for the usage information.
func (cmds *CommandSet) ShowServer(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForShow()
		return util.E_PEBKAC
	}
	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("server name cannnot be blank")
		return util.E_PEBKAC
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}
	vm, err := cmds.client.GetVirtualMachine(name)

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

// ShowGroup implements the show-group command, which is used to show the group name and ID, as well as the servers within it.
func (cmds *CommandSet) ShowGroup(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	verbose := flags.Bool("verbose", false, "")
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "group")
	if !ok {
		cmds.HelpForShow()
		return util.E_PEBKAC
	}
	name := cmds.client.ParseGroupName(nameStr, cmds.config.GetGroup())

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	group, err := cmds.client.GetGroup(name)

	if err != nil {
		return util.ProcessError(err)
	}

	if !cmds.config.Silent() {

		if *jsonOut {
			js, _ := json.MarshalIndent(group, "", "    ")
			log.Output(string(js))
		} else {
			s := ""
			if len(group.VirtualMachines) != 1 {
				s = "s"
			}
			log.Outputf("%s - Group containing %d cloud server%s\r\n", group.Name, len(group.VirtualMachines), s)

			if *verbose || len(group.VirtualMachines) <= 3 {
				log.Output()
				for _, v := range util.FormatVirtualMachines(group.VirtualMachines) {
					log.Output(v)
				}

			}
		}
	}
	return util.E_SUCCESS

}

// ShowAccount implements the show-account command, which is used to show the client account name, as well as the groups and servers within it.
func (cmds *CommandSet) ShowAccount(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	verbose := flags.Bool("verbose", false, "")
	jsonOut := flags.Bool("json", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "account")
	if !ok {
		cmds.HelpForShow()
		return util.E_PEBKAC
	}
	name := cmds.client.ParseAccountName(nameStr, cmds.config.GetIgnoreErr("account"))

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	acc, err := cmds.client.GetAccount(name)

	if err != nil {
		return util.ProcessError(err)
	}

	if *jsonOut {
		js, _ := json.MarshalIndent(acc, "", "    ")
		log.Output(string(js))
	} else {
		log.Output(util.FormatAccount(acc))

		switch {
		case *verbose:
			for _, g := range acc.Groups {
				log.Outputf("Group %s\r\n", g.Name)
				for _, v := range util.FormatVirtualMachines(g.VirtualMachines) {
					log.Output(v)
				}
			}
		}
	}
	return util.E_SUCCESS

}

func (cmds *CommandSet) ShowUser(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	username, ok := util.ShiftArgument(&args, "username")
	if !ok {
		cmds.HelpForShow()
		return util.E_PEBKAC
	}

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	user, err := cmds.client.GetUser(username)
	if err != nil {
		return util.ProcessError(err)
	}

	log.Outputf("User %s:\n\nAuthorized keys:\n", user.Username)
	for _, k := range user.AuthorizedKeys {
		log.Output(k)
	}
	return util.E_SUCCESS
}
