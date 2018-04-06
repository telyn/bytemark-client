package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/add"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "add",
		Usage:     "add servers, discs, etc - see `bytemark help add <kind of thing> `",
		UsageText: "add server|group|disc",
		Description: `add a new group, server or disc

  add disc[s] [--disc <disc spec>]... <cloud server>
  add group [--account <name>] <name>
  add server (see bytemark help create server)

A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to create multiple discs`,
		Action:      cli.ShowSubcommandHelp,
		Subcommands: add.Commands,
	})
}
