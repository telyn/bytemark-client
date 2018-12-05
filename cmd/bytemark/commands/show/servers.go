package show

import (
	"errors"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "servers",
		Usage:     "show all the servers in an account",
		UsageText: "show servers [--group <group> | --account <account>] [group]",
		Description: `This command shows all the servers in the given account, or in your default account if not specified.
Deleted servers are included in the list, with ' (deleted)' appended.

If --group and --account are specified, the group will be displayed and the account will be ignored.`,
		Flags: append(app.OutputFlags("servers", "array"),
			cli.GenericFlag{
				Name:  "group",
				Usage: "the group to list the servers of",
				Value: new(flags.GroupNameFlag),
			},
			// TODO: change to AccountNameFlag
			cli.StringFlag{
				Name:  "account",
				Usage: "the account to show all the servers of",
			},
		),
		Action: app.Action(args.Optional("group"), with.Auth, func(c *app.Context) error {
			servers := brain.VirtualMachines{}
			if c.IsSet("group") {
				groupName := flags.GroupName(c, "group")
				group, err := c.Client().GetGroup(groupName)
				if err != nil {
					return err
				}
				return c.OutputInDesiredForm(brain.VirtualMachines(group.VirtualMachines), output.List)
			}
			if c.IsSet("account") {
				err := with.Account("account")(c)
				if err != nil {
					return err
				}
				for _, g := range c.Account.Groups {
					servers = append(servers, g.VirtualMachines...)
				}
				return c.OutputInDesiredForm(servers, output.List)
			}
			return errors.New("A group or account must be specified")
		}),
	})
}
