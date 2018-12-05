package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "groups",
		Usage:       "show all the groups in an account",
		UsageText:   "show groups [account]",
		Description: `This command shows all the groups in the given account, or in your default account if not specified.`,
		Flags: append(app.OutputFlags("groups", "array"),
			cli.GenericFlag{
				Name:  "account",
				Usage: "the account to list the groups of",
				Value: new(flags.AccountName),
			},
		),
		Action: app.Action(args.Optional("account"), with.RequiredFlags("account"), with.Account("account"), func(c *app.Context) error {
			return c.OutputInDesiredForm(c.Account.Groups, output.List)
		}),
	})
}
