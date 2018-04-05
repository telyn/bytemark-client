package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin/add"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "add",
		Usage:       "add an ip range or a user - see `bytemark --admin help add <kind of thing> `",
		UsageText:   "add ip range|user",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: add.Commands,
	})
}
