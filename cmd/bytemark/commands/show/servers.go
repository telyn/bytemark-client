package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "servers",
		Usage:     "show all the servers in an account",
		UsageText: "show servers [account]",
		Description: `This command shows all the servers in the given account, or in your default account if not specified.
Deleted servers are included in the list, with ' (deleted)' appended.`,
		Flags: append(app.OutputFlags("servers", "array"),
			cli.GenericFlag{
				Name:  "account",
				Usage: "the account to list the servers of",
				Value: new(app.AccountNameFlag),
			},
		),
		Action: app.Action(args.Optional("account"), with.Account("account"), with.Auth, func(c *app.Context) error {
			servers := brain.VirtualMachines{}

			for _, g := range c.Account.Groups {
				servers = append(servers, g.VirtualMachines...)
			}
			return c.OutputInDesiredForm(servers, output.List)
		}),
	})
}
