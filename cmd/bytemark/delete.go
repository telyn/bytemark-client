package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"fmt"
	"strings"
)

// HelpForDelete outputs usage information for the delete command
func (cmds *CommandSet) HelpForDelete() util.ExitCode {
	log.Log("bytemark delete")
	log.Log()
	log.Log("usage: bytemark delete account <account>")
	log.Log("       bytemark delete disc <server> <label>")
	log.Log("       bytemark delete group [--recursive] <group>")
	//log.Log("       bytemark delete user <user>")
	log.Log("       bytemark delete key [--user=<user>] <public key identifier>")
	log.Log("       bytemark delete server [--force] [---purge] <server>")
	log.Log("       bytemark undelete server <server>")
	log.Log()
	log.Log("Deletes the given server, disc, group, account or key. Only empty groups and accounts can be deleted.")
	log.Log("If the --purge flag is given and the target is a cloud server, will permanently delete the server. Billing will cease and you will be unable to recover the server or its data.")
	log.Log("If the --force flag is given, you will not be prompted to confirm deletion.")
	log.Log()
	log.Log("The undelete server command may be used to restore a deleted (but not purged) server to its state prior to deletion.")
	log.Log()
	return util.E_USAGE_DISPLAYED
}

// DeleteServer implements the delete server command, which is used to delete and purge Bytemark servers. See HelpForDelete for usage information.
func (cmds *CommandSet) DeleteServer(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	purge := flags.Bool("purge", false, "")

	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}
	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Server name cannot be blank.")
		cmds.HelpForDelete()
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
	if vm.Deleted && !*purge {
		log.Errorf("Server %s has already been deleted.\r\nIf you wish to permanently delete it, add --purge\r\n", vm.Hostname)
		return util.E_SUCCESS
	}

	if !cmds.config.Force() {
		fstr := fmt.Sprintf("Are you certain you wish to delete %s?", vm.Hostname)
		if *purge {
			fstr = fmt.Sprintf("Are you certain you wish to permanently delete %s? You will not be able to un-delete it.", vm.Hostname)

		}
		if !util.PromptYesNo(fstr) {
			return util.ProcessError(&util.UserRequestedExit{})

		}
	}

	err = cmds.client.DeleteVirtualMachine(name, *purge)

	if err != nil {
		return util.ProcessError(err)
	}

	if *purge {
		log.Logf("Server %s purged successfully.\r\n", vm.Hostname)
	} else {
		log.Logf("Server %s deleted successfully.\r\n", vm.Hostname)
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) DeleteDisc(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	//purge := flags.Bool("purge", false, "")

	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "server")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}

	disc, ok := util.ShiftArgument(&args, "disc id")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	if !(cmds.config.Force() || util.PromptYesNo("Are you sure you wish to delete this disc? It is impossible to recover.")) {
		log.Log("Cancelling.")
		return util.E_USER_EXIT
	}

	err = cmds.client.DeleteDisc(name, disc)
	if err != nil {
		return util.ProcessError(err)
	}

	return util.E_SUCCESS

}

func (cmds *CommandSet) DeleteGroup(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()

	recursive := flags.Bool("recursive", false, "")

	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "group")
	if !ok {
		cmds.HelpForDelete()
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

	if len(group.VirtualMachines) > 0 {
		if *recursive {

			log.Log("WARNING: The following servers will be permanently deleted, without any way to recover or un-delete them:")
			for _, vm := range group.VirtualMachines {
				log.Logf("\t%s\r\n", vm.Name)
			}
			log.Log("", "")
			if util.PromptYesNo("Are you sure you want to continue?") {
				vmn := lib.VirtualMachineName{Group: name.Group, Account: name.Account}
				for _, vm := range group.VirtualMachines {
					vmn.VirtualMachine = vm.Name
					err := cmds.client.DeleteVirtualMachine(vmn, true)
					if err != nil {
						return util.ProcessError(err)
					} else {
						log.Logf("Server %s purged successfully.\r\n", name)
					}

				}
			}
		} else {
			log.Errorf("Group %s contains servers, will not be deleted without --recursive\r\n", name.Group)
			return util.E_WONT_DELETE_NONEMPTY
		}
	}
	return util.ProcessError(cmds.client.DeleteGroup(name))

}

// DeleteKey implements the delete key command, which is used to remove an authorized_key from a user.
func (cmds *CommandSet) DeleteKey(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()

	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	user := cmds.config.GetIgnoreErr("user")

	key := strings.Join(args, " ")
	if key == "" {
		log.Log("You must specify a key to delete.\r\n")
		cmds.HelpForDelete()
		return util.E_SUCCESS

	}

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.client.DeleteUserAuthorizedKey(user, key)
	if err == nil {
		log.Log("Key deleted successfullly")
		return util.E_SUCCESS
	} else {
		return util.ProcessError(err)
	}
}

// UndeleteServer implements the undelete server command, which is used to remove the deleted flag from Bytemark servers, allowing them to be reactivated.
func (cmds *CommandSet) UndeleteServer(args []string) util.ExitCode {

	nameStr, ok := util.ShiftArgument(&args, "cloud server")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}
	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		log.Error("Server name cannot be blank")
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

	if !vm.Deleted {
		log.Errorf("%s was already undeleted\r\n", vm.Hostname)
		return util.E_SUCCESS
	}

	err = cmds.client.UndeleteVirtualMachine(name)

	if err != nil {
		return util.ProcessError(err)
	}
	log.Logf("Successfully restored %s\r\n", vm.Hostname)

	return util.E_SUCCESS
}
