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
	switch strings.ToLower(args[0]) {
	case "show":
		dispatch.HelpForShow()
		return
	case "set":
		dispatch.HelpForSet()
		return
	case "unset":
		dispatch.HelpForUnset()
		return
	}

}
