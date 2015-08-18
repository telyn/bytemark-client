package cmds

import (
	"bigv.io/client/cmds/util"
	"fmt"
	"strings"
)

// HelpForHelp shows overall usage information for the BigV client, including a list of available commands.
func (cmds *CommandSet) HelpForHelp() util.ExitCode {
	fmt.Println("bigv command-line client (the new, cool one)")
	fmt.Println()
	fmt.Println("Usage")
	fmt.Println()
	fmt.Println("    go-bigv [flags] <command> [flags] [args]")
	fmt.Println()
	fmt.Println("Commands available")
	fmt.Println()
	fmt.Println("    help, config, create, debug, delete, list, show")
	fmt.Println("    AND MAYBE MORE OR FEWER - THIS LIST IS NOT FINAL")
	fmt.Println()
	fmt.Println("See `go-bigv help <command>` for help specific to a command")
	fmt.Println()
	return util.E_USAGE_DISPLAYED
}

// Help implements the help command, which gives usage information specific to each command. Usage: bigv help [command]
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
