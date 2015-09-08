package cmds

import (
	"bigv.io/client/cmds/util"
	"bigv.io/client/util/log"
	"strings"
)

// HelpForConfig outputs usage information for the bigv config command.
func (cmds *CommandSet) HelpForConfig() util.ExitCode {
	log.Log("go-bigv config")
	log.Log()
	log.Log("Usage:")
	log.Log("    go-bigv config")
	log.Log("        Outputs the current values of all variables and what source they were derived from")
	log.Log()
	log.Log("    go-bigv config set <variable> <value>")
	log.Log("        Sets a variable by writing to your bigv config (usually ~/.go-bigv)")
	log.Log()
	log.Log("    go-bigv config unset <variable>")
	log.Log("        Unsets a variable by removing data from bigv config (usually ~/.go-bigv)")
	log.Log()
	log.Log("Available variables:")
	log.Log("    endpoint - the BigV endpoint to connect to. https://uk0.bigv.io is the default")
	log.Log("    auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.")
	log.Log("    debug-level - the default debug level. Set to 0 unless you like lots of output")
	log.Log("    token - the token used for authentication.") // You can get one using bigv auth.")
	log.Log()
	return util.E_USAGE_DISPLAYED
}

// Config provides the bigv config command, which sets variables in the user's config. See HelpForConfig for usage information.
// It's slightly more user friendly than echo "value" > ~/.go-bigv/
func (cmds *CommandSet) Config(args []string) util.ExitCode {
	if len(args) == 0 {
		vars, err := cmds.config.GetAll()
		if err != nil {
			return util.ProcessError(err)
		}
		for _, v := range vars {
			log.Logf("%s\t: '%s' (%s)\r\n", v.Name, v.Value, v.Source)
		}
		return util.E_SUCCESS
	} else if len(args) == 1 {
		cmds.HelpForConfig()
		return util.E_SUCCESS
	}

	switch strings.ToLower(args[0]) {
	case "set":
		variable := strings.ToLower(args[1])

		oldVar, err := cmds.config.GetV(variable)
		if err != nil {
			if e, ok := err.(*util.ConfigReadError); ok {
				log.Errorf("Couldn't read the old value of %s - %v\r\n", e.Name, e.Err)
			} else {
				log.Errorf("Couldn't read the old value of %s - %v\r\n", variable, err)
			}
			return util.E_CANT_READ_CONFIG
		}

		if len(args) == 2 {
			log.Logf("%s: '%s' (%s)\r\n", oldVar.Name, oldVar.Value, oldVar.Source)
			return util.E_SUCCESS
		}

		// TODO(telyn): consider validating input for the set command
		err = cmds.config.SetPersistent(variable, args[2], "CMD set")
		if err != nil {
			if e, ok := err.(*util.ConfigReadError); ok {
				log.Errorf("Couldn't set %s - %v\r\n", e.Name, e.Err)
			} else {
				log.Errorf("Couldn't set %s - %v\r\n", variable, err)
			}
			return util.E_CANT_WRITE_CONFIG
		}

		if oldVar.Source == "config" {
			log.Logf("%s has been changed.\r\nOld value: %s\r\nNew value: %s\r\n", variable, oldVar.Value, args[1])
		} else {
			log.Logf("%s has been set. \r\nNew value: %s\r\n", variable, args[1])
		}

	case "unset":
		variable := strings.ToLower(args[0])
		err := cmds.config.Unset(variable)
		return util.ProcessError(err)
	default:
		log.Errorf("Unrecognised command %s\r\n", args[0])
		cmds.HelpForConfig()
	}
	return util.E_SUCCESS
}
