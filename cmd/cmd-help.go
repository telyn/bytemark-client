package main

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
	case "config":
		cmds.HelpForConfig()
	case "debug":
		cmds.HelpForDebug()
	case "delete":
		cmds.HelpForDelete()
	case "exit":
		cmds.HelpForExitCodes()
	case "exit-codes":
		cmds.HelpForExitCodes()
	case "show":
		cmds.HelpForShow()
	}

}
