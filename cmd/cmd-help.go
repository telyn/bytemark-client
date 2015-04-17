package cmd

import (
	"fmt"
	"strings"
)

func (cmds *CommandSet) HelpForHelp() {
	fmt.Println("bigv command-line client (the new, cool one)")
	fmt.Println()
	fmt.Println("There would be some usage output here if I had actually written any.")
	fmt.Println()
}

func (cmds *CommandSet) Help(args []string) {
	if len(args) == 0 {
		cmds.HelpForHelp()
		return
	}

	// please try and keep these in alphabetical order
	switch strings.ToLower(args[0]) {
	case "debug":
		cmds.HelpForDebug()
	case "set":
		cmds.HelpForSet()
	case "show":
		cmds.HelpForShow()
	case "unset":
		cmds.HelpForUnset()
	}

}
