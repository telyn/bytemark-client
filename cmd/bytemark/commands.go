package main

import (
	"github.com/urfave/cli"
	"sort"
	"strings"
)

func init() {
	// both of these functions have the lists added at the end of baseAppSetup
	commands = append(commands, cli.Command{
		Name:        "admin",
		Usage:       "cluster-admin only commands",
		UsageText:   "bytemark admin",
		Description: "This is a list of cluster-admin only commands. In order to use these commands or access their help you must add the '--admin' flag after 'bytemark'.\r\n   For example, to see the help for the 'migrate server' command, 'bytemark --admin help server' should be executed.\r\n\r\nALL COMMANDS:\r\n\r\n",
		Hidden:      true,
		Action:      cli.ShowSubcommandHelp,
	}, cli.Command{
		Name:        "commands",
		Usage:       "list of all commands available",
		UsageText:   "bytemark commands",
		Description: "ALL COMMANDS:\r\n\r\n",
		Action:      cli.ShowSubCommandHelp,
	})
}

func generateCommandsHelp(cmds []cli.Command) string {
	commandsUsage := make([]string, 0, len(commands))
	for _, cmd := range cmds {
		commandsUsage = append(commandsUsage, generateSubcommandsUsage(cmd, "")...)
	}

	sort.Strings(commandsUsage)

	return "   " + strings.Join(commandsUsage, "\r\n   ")
}

func generateSubcommandsUsage(cmd cli.Command, prefix string) (commandsUsage []string) {
	if cmd.Subcommands == nil || len(cmd.Subcommands) == 0 {
		fullName := prefix + cmd.Name

		return []string{fullName + ": " + cmd.Usage}
	}
	commandsUsage = make([]string, 0, len(cmd.Subcommands))
	for _, subcmd := range cmd.Subcommands {
		commandsUsage = append(commandsUsage, generateSubcommandsUsage(subcmd, prefix+cmd.Name+" ")...)
	}
	return
}
