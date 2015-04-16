package cmd

import (
	"fmt"
	"strings"
)

func (dispatch *Dispatcher) HelpForHelp() {
	fmt.Println("bigv command-line client (the new, cool one)")
	fmt.Println()
	fmt.Println("There would be some usage output here if I had actually written any.")
	fmt.Println()
}

func (dispatch *Dispatcher) Help(args []string) {
	if len(args) == 0 {
		dispatch.HelpForHelp()
		return
	}

	// please try and keep these in alphabetical order
	switch strings.ToLower(args[0]) {
	case "debug":
		dispatch.HelpForDebug()
	case "set":
		dispatch.HelpForSet()
	case "show":
		dispatch.HelpForShow()
	case "unset":
		dispatch.HelpForUnset()
	}

}
