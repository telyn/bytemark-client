package main

import (
	"fmt"
	"strings"
)

// HelpForConfig outputs usage information for the bigv config command.
func (cmds *CommandSet) HelpForConfig() {
	fmt.Println("go-bigv config")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    go-bigv config")
	fmt.Println("        Outputs the current values of all variables and what source they were derived from")
	fmt.Println()
	fmt.Println("    go-bigv config set <variable> <value>")
	fmt.Println("        Sets a variable by writing to your bigv config (usually ~/.go-bigv)")
	fmt.Println()
	fmt.Println("    go-bigv config unset <variable>")
	fmt.Println("        Unsets a variable by removing data from bigv config (usually ~/.go-bigv)")
	fmt.Println()
	fmt.Println("Available variables:")
	fmt.Println("    endpoint - the BigV endpoint to connect to. https://uk0.bigv.io is the default")
	fmt.Println("    auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.")
	fmt.Println("    debug-level - the default debug level. Set to 0 unless you like lots of output")
	fmt.Println()
}

// Config provides the bigv config command, which sets variables in the user's config. See HelpForConfig for usage information.
// It's slightly more user friendly than echo "value" > ~/.go-bigv/
func (cmds *CommandSet) Config(args []string) ExitCode {
	if len(args) == 0 {
		for _, v := range cmds.config.GetAll() {
			fmt.Printf("%s\t: '%s' (%s)\r\n", v.Name, v.Value, v.Source)
		}
		return exit(nil)
	} else if len(args) == 1 {
		cmds.HelpForConfig()
		return exit(nil)
	}

	switch strings.ToLower(args[0]) {
	case "set":
		variable := strings.ToLower(args[1])

		oldVar := cmds.config.GetV(variable)

		if len(args) == 2 {
			fmt.Printf("%s: '%s' (%s)\r\n", oldVar.Name, oldVar.Value, oldVar.Source)
			return exit(nil)
		}

		// TODO(telyn): consider validating input for the set command
		cmds.config.SetPersistent(variable, args[2], "CMD set")

		if oldVar.Source == "config" && !cmds.config.GetBool("silent") {
			fmt.Printf("%s has been changed.\r\nOld value: %s\r\nNew value: %s\r\n", variable, oldVar.Value, args[1])
		} else if !cmds.config.GetBool("silent") {
			fmt.Printf("%s has been set. \r\nNew value: %s\r\n", variable, args[1])
		}

	case "unset":
		// TODO(telyn): write this...
	default:
		fmt.Printf("Unrecognised command %s\r\n", args[0])
		cmds.HelpForConfig()
	}
	return exit(nil)
}
