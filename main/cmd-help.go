package main

import (
	"fmt"
	"strings"
)

// HelpForHelp shows overall usage information for the BigV client, including a list of available commands.
func (cmds *CommandSet) HelpForHelp() {
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
}

// Help implements the help command, which gives usage information specific to each command. Usage: bigv help [command]
func (cmds *CommandSet) Help(args []string) {
	if len(args) == 0 {
		cmds.HelpForHelp()
		return
	}

	// please try and keep these in alphabetical order
	switch strings.ToLower(args[0]) {
	case "config":
		cmds.HelpForConfig()
	case "create":
		cmds.HelpForCreate()
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
