package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "accounts",
		Usage:       "show all the accounts you're able to see",
		UsageText:   "show accounts",
		Description: `This will show all the accounts that your authentication token has some form of access to.`,
		Flags:       app.OutputFlags("accounts", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			accounts, err := c.Client().GetAccounts()

			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(accounts, output.List)
		}),
	})
}
