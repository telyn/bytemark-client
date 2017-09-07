// Package cliutil is a collection of functions to help work with urfave/cli
package cliutil

import (
	"strings"

	"github.com/urfave/cli"
)

// CreateMultiwordCommand creates a new cli.Command with subcommands for each
// word with the innermost one being identical to the original, but for the name
// and for being hidden.
//
// For example, given a command called "example command", this creates a command
// called "example" with a subcommand called "command" - where "command" is
// identical to "example command".
// This is a workaround for urfave/cli not supporting multi-word commands.
func CreateMultiwordCommand(orig cli.Command) cli.Command {
	if !strings.Contains(orig.Name, " ") {
		return orig
	}

	cmdNameParts := strings.Split(orig.Name, " ")
	lastIndex := len(cmdNameParts) - 1

	// create the innermost Command
	cmd := orig
	cmd.Name = cmdNameParts[lastIndex]

	lastIndex--
	for lastIndex >= 0 {
		name := cmdNameParts[lastIndex]
		cmd = cli.Command{
			Name:        name,
			Hidden:      true,
			Subcommands: []cli.Command{cmd},
		}
		lastIndex--
	}
	return cmd
}

// CreateMultiwordCommands is a workaround for urfave/cli not fully supporting multi-word Names for commands in a good way
// it looks at each subcommand of cmds, and creates hidden command trees to match any multi-word commands found.
func CreateMultiwordCommands(cmds []cli.Command) (newCmds []cli.Command) {
	newCmds = make([]cli.Command, 0, len(newCmds))
	for _, topLevelCommand := range cmds {
		for _, cmd := range topLevelCommand.Subcommands {
			if strings.Contains(cmd.Name, " ") {
				newCommand := CreateMultiwordCommand(cmd)

				// use mergeCommands to add our new command into the top level command
				topLevelCommand.Subcommands = MergeCommands(topLevelCommand.Subcommands, []cli.Command{newCommand})
			}
		}
		newCmds = append(newCmds, topLevelCommand)
	}
	return
}
