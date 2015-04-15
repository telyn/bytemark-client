package cmd

import (
	"fmt"
	"os"
	"strings"
)

func (dispatch *Dispatcher) Help(args []string) {
	if len(args) == 0 {
		fmt.Println("bigv command-line client (the new, cool one)")
		fmt.Println()
		fmt.Println("There would be some usage output here if I had actually written any.")
		os.Exit(0)
	}
	switch strings.ToLower(args[0]) {
	case "help":
		dispatch.HelpForShow()
	}

}
