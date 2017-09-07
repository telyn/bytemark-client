// cliutil is a collection of functions to help work with urfave/cli
package cliutil

import (
	"strings"

	"github.com/urfave/cli"
)

func CreateMultiwordCommand(orig cli.Command) cli.Command {
	if !strings.Contains(orig.Name, " ") {
		return orig
	}

	cmdNameParts := strings.Split(orig.Name, " ")
	lastIndex := len(cmdNameParts) - 1

	// create the innermost Command
	cmd := cli.Command{
		Name:        cmdNameParts[lastIndex],
		Usage:       orig.Usage,
		UsageText:   orig.UsageText,
		Description: orig.Description,
		Flags:       orig.Flags,
		Action:      orig.Action,
		Hidden:      true,
	}

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
