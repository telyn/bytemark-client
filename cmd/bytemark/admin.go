package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"sort"
	"strings"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "admin",
		Usage:       "cluster-admin only commands",
		UsageText:   "bytemark admin",
		Description: "This is a list of cluster-admin only commands. In order to use these commands or access their help you must add the '--admin' flag after 'bytemark'.\r\n   For example, to see the help for the 'migrate server' command, 'bytemark --admin help server' should be executed.",
		Hidden:      true,
		Action: With(func(c *Context) error {
			log.Output(generateAdminCommandsHelp())
			return nil
		}),
	})
}

func generateAdminCommandsHelp() string {
	commandsUsage := make([]string, 0, len(adminCommands))
	for _, cmd := range adminCommands {
		commandsUsage = append(commandsUsage, generateSubcommandsUsage(cmd, "")...)
	}

	sort.Strings(commandsUsage)

	return "ADMIN COMMANDS:\r\n\r\n   " + strings.Join(commandsUsage, "\r\n   ")
}

func generateSubcommandsUsage(cmd cli.Command, prefix string) (commandsUsage []string) {
	if cmd.Subcommands == nil || len(cmd.Subcommands) == 0 {
		fullName := prefix + " " + cmd.Name

		return []string{fullName + ": " + cmd.Usage}
	}
	commandsUsage = make([]string, 0, len(cmd.Subcommands))
	for _, subcmd := range cmd.Subcommands {
		if prefix == "" {
			commandsUsage = append(commandsUsage, generateSubcommandsUsage(subcmd, cmd.Name)...)
		} else {
			commandsUsage = append(commandsUsage, generateSubcommandsUsage(subcmd, prefix+" "+cmd.Name)...)
		}
	}
	return
}
