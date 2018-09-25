package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/delete"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "delete",
		Usage:     "delete a given server, disc, group, account or key",
		UsageText: "delete account|disc|group|key|server",
		Description: `delete a given server, disc, group, account or key

   Only empty groups and accounts can be deleted.

   The restore server command may be used to restore a deleted (but not purged) server to its state prior to deletion.`,
		Action:      cli.ShowSubcommandHelp,
		Subcommands: delete.Commands,
	})
}
