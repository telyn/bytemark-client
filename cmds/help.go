package cmds

import (
	"bigv.io/client/cmds/util"
	"bigv.io/client/util/log"
	"strings"
)

// HelpForHelp shows overall usage information for the BigV client, including a list of available commands.
func (cmds *CommandSet) HelpForHelp() util.ExitCode {
	log.Log("Bytemark command-line client")
	log.Log()
	log.Log("Usage")
	log.Log()
	log.Log("    bytemark [flags] <command> [flags] [args]")
	log.Log()
	log.Log("Commands available")
	log.Log()
	log.Log("    help, config, create, debug, delete,")
	log.Log("    list, resize, show, undelete")
	log.Log("    AND MAYBE MORE OR FEWER - THIS LIST IS NOT FINAL")
	log.Log()
	log.Log("See `bytemark help <command>` for help specific to a command")
	log.Log()
	return util.E_USAGE_DISPLAYED
}

// Help implements the help command, which gives usage information specific to each command. Usage: bytemark help [command]
func (cmds *CommandSet) Help(args []string) util.ExitCode {
	if len(args) == 0 {
		return cmds.HelpForHelp()
	}

	// please try and keep these in alphabetical order
	switch strings.ToLower(args[0]) {
	case "config":
		return cmds.HelpForConfig()
	case "create":
		return cmds.HelpForCreate()
	case "debug":
		return cmds.HelpForDebug()
	case "delete":
		return cmds.HelpForDelete()
	case "exit":
		return util.HelpForExitCodes()
	case "exit-codes":
		return util.HelpForExitCodes()
	case "show":
		return cmds.HelpForShow()
	}
	return cmds.HelpForHelp()

}
