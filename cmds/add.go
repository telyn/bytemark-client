package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
	"io/ioutil"
	"strings"
)

func (cmds *CommandSet) HelpForAdd() util.ExitCode {
	log.Log("bytemark add")
	log.Log()
	log.Log("usage: bytemark add key [--public-key-file=<filename>] [--user=<user>] [<public key>]")
	log.Log()
	log.Log("Add the given public key to the given user (or the default user). This will allow them")
	log.Log("to use that key to access management IPs they have access to using that key.")
	log.Log("Specify --public-key-file=- to read the public key from stdin")
	log.Log("--public-key-file will be ignored if a public key is specified in the arguments")
	log.Log("To remove a key, use the remove key command.")
	return util.E_USAGE_DISPLAYED
}

// AddKey implements the add key command, which is used to add an authorized_key from a user.
func (cmds *CommandSet) AddKey(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	keyFile := flags.String("public-key-file", "", "")

	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	user := cmds.config.GetIgnoreErr("user")

	key := strings.TrimSpace(strings.Join(args, " "))
	if key == "" {
		if *keyFile == "" {
			cmds.HelpForAdd()
			return util.E_PEBKAC
		}

		keyBytes, err := ioutil.ReadFile(*keyFile)
		if err != nil {
			return util.ProcessError(err)
		}
		key = string(keyBytes)
	}

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.bigv.AddUserAuthorizedKey(user, key)
	if err == nil {
		log.Log("Key added successfully")
		return util.E_SUCCESS
	} else {
		return util.ProcessError(err)
	}
}
