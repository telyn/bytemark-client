package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/urfave/cli"
)

func init() {

	commands = append(commands, cli.Command{
		Name:      "overview",
		Usage:     `overview of your Bytemark hosting`,
		UsageText: "bytemark show account [--json]",
		Description: `This command displays an overview of the hosting you have with Bytemark.

		If the --json flag is specified, prints a complete overview of the account in JSON format, including all groups and their servers.`,
		Flags: OutputFlags("account details", "object"),
		Action: With(AuthProvider, func(c *Context) error {

			allAccs, err := c.Client().GetAccounts()
			if err != nil {
				return err
			}

			accName := c.Config().GetIgnoreErr("account")
			var def lib.Account
			if accName != "" {
				def, err = c.Client().GetAccount(accName)
				if err != nil {
					return err
				}
			} else {

				def, err = c.Client().GetDefaultAccount()
				if err != nil {
					return err
				}
			}

			// TODO(telyn) refactor this to be somewhere else (ideally GetAccount/GetAccounts would fill in IsDefaultAccount automatically)
			for _, acc := range allAccs {
				if acc.Name != "" && def.Name != "" && acc.Name == def.Name {
					acc.IsDefaultAccount = true
				} else if acc.BillingID != 0 && def.BillingID != 0 && acc.BillingID == def.BillingID {
					acc.IsDefaultAccount = true
				}
			}
			overview := lib.Overview{
				Accounts:       allAccs,
				DefaultAccount: def,
				Username:       c.Client().GetSessionUser(),
			}

			return c.OutputInDesiredForm(overview)

		}),
	})
}
