package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin/migrate"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "migrate",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: migrate.Commands,
	})
}
