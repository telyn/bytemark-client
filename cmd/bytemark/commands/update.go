package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/update"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "update",
		Usage:       "modify servers and other Bytemark Cloud resources",
		UsageText:   "update server",
		Description: `update an existing server`,
		Action:      cli.ShowSubcommandHelp,
		Subcommands: update.Commands,
	})
}
