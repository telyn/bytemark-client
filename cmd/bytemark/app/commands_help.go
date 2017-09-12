package app

import (
	"sort"
	"strings"

	"github.com/urfave/cli"
)

func GenerateCommandsHelp(cmds []cli.Command) string {
	commandsUsage := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		commandsUsage = append(commandsUsage, GenerateSubcommandsUsage(cmd, "")...)
	}

	sort.Strings(commandsUsage)

	return "   " + strings.Join(commandsUsage, "\r\n   ")
}

func GenerateSubcommandsUsage(cmd cli.Command, prefix string) (commandsUsage []string) {
	if cmd.Subcommands == nil || len(cmd.Subcommands) == 0 {
		fullName := prefix + cmd.Name

		return []string{fullName + ": " + cmd.Usage}
	}
	commandsUsage = make([]string, 0, len(cmd.Subcommands))
	for _, subcmd := range cmd.Subcommands {
		commandsUsage = append(commandsUsage, GenerateSubcommandsUsage(subcmd, prefix+cmd.Name+" ")...)
	}
	return
}
