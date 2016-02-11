package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
)

func (c *CommandManager) HelpForSignup() util.ExitCode {
	log.Log("usage: bytemark signup [--new-user]")
	log.Log("")
	log.Log("Create a new Bytemark account. If you already have a login,")
	log.Log("you can use this to create a new account attached to you,")
	log.Log("or specify --new-user to create a new account with a new user")
	log.Log("as if signing up for the first time again")
}

func (cmds *CommandManager) Signup(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	newUser := flags.Bool("disc", false, "")
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	// TODO(telyn): check a terminal is attached to stdin to try to help prevent fraudy/spammy crap just in case
	token := cmds.Config.GetIgnoreErr("token")
	if token == "" {
		user, err := cmds.Config.Get("user")
		if err != nil || user.Source != "ENV USER" {
			*newUser = true
		}
	}

	if *newUser {

	}
}
