package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "reap",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "servers",
				Aliases:   []string{"vm", "vms"},
				Usage:     "triggers server reaping, purging all deleted servers and discs",
				UsageText: "bytemark --admin reap servers",
				Action: app.Action(with.Auth, func(c *app.Context) error {
					if err := c.Client().ReapVMs(); err != nil {
						return err
					}

					log.Output("Reap initiated")

					return nil
				}),
			},
		},
	})
}
