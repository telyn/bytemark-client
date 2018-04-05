package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/show"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "show",
		Usage:     "displays information about your account and of your assets at Bytemark",
		UsageText: "show accounts, discs, groups, etc - see `bytemark help show <kind of thing> `",
		Description: `displays information about the given server, group, or account

Plurals are scripting-friendly lists of your assets at Bytemark, showing the kind of object you request, one per line.

Perfect for piping into a bash while loop!`,
		Action:      cli.ShowSubcommandHelp,
		Subcommands: show.Commands,
	})
}
