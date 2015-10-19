package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
	"strings"
)

func (cmds *CommandSet) HelpForAdd() util.ExitCode {
	log.Log("bytemark add")
	log.Log()
	log.Log("usage: bytemark add key <user> <public key>")
	log.Log()
	log.Log("Add the given public key to the specified user. This will allow them")
	log.Log("to use that key to access management IPs they have access to using that key.")
	return util.E_USAGE_DISPLAYED
}

// AddKey implements the delete key command, which is used to remove an authorized_key from a user.
func (cmds *CommandSet) AddKey(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()

	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "user")
	if !ok {
		cmds.HelpForDelete()
		return util.E_PEBKAC
	}
	key := strings.TrimSpace(strings.Join(args, " "))

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.bigv.AddUserAuthorizedKey(nameStr, key)
	if err == nil {
		log.Log("Key added successfully")
		return util.E_SUCCESS
	} else {
		return util.ProcessError(err)
	}
}
