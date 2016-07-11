package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/urfave/cli"
	"os"
)

func init() {

	commands = append(commands, cli.Command{
		Name:      "overview",
		Usage:     `overview of your Bytemark hosting`,
		UsageText: "bytemark show account [--json] [name]",
		Description: `This command displays an overview of the hosting you have with Bytemark.

		If the --json flag is specified, prints a complete overview of the account in JSON format, including all groups and their servers.`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "Output account details as a JSON object",
			},
		},
		Action: With(AuthProvider, func(c *Context) error {

			allAccs, err := global.Client.GetAccounts()
			if err != nil {
				return err
			}

			accName := global.Config.GetIgnoreErr("account")
			var def *lib.Account
			if accName != "" {
				def, err = global.Client.GetAccount(accName)
				if err != nil {
					return err
				}
			} else {
				def, err = global.Client.GetDefaultAccount()
				if err != nil {
					return err
				}
			}

			return lib.FormatOverview(os.Stdout, allAccs, def, global.Client.GetSessionUser())

		}),
	})
}
