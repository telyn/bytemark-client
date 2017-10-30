package main

import (
	"github.com/urfave/cli"
)

func init() {
	// both of these functions have the lists added at the end of baseAppSetup
	commands = append(commands, cli.Command{
		Name:        "admin",
		Usage:       "cluster-admin only commands",
		UsageText:   "bytemark admin",
		Description: "This is a list of cluster-admin only commands. In order to use these commands or access their help you must add the '--admin' flag after 'bytemark'.\r\n   For example, to see the help for the 'migrate server' command, 'bytemark --admin help server' should be executed.\r\n\r\nALL COMMANDS:\r\n",
		Hidden:      true,
		Action:      cli.ShowSubcommandHelp,
	}, cli.Command{
		Name:        "commands",
		Usage:       "list of all commands available",
		UsageText:   "bytemark commands",
		Description: "This is a list of all commands currently available. Toggling --admin on or off will add/remove the admin commands and enable/disable admin features in some commands.\r\nALL COMMANDS:\r\n",
		Action:      cli.ShowSubcommandHelp,
	})
}
