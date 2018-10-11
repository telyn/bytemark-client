package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "cancel",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "migration",
				Usage:     "cancel a disc or server migration",
				UsageText: "--admin cancel migration (--disc <disc> | --server <server>)",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "disc",
						Usage: "the ID of the disc that is migrating",
					},
					cli.IntFlag{
						Name:  "server",
						Usage: "the ID of the server that is migrating",
					},
				},
				Action: app.Action(args.Optional("disc", "server"), with.Auth, cancelMigration),
			},
		},
	})
}

func cancelMigration(c *app.Context) error {
	discID := c.Int("disc")
	serverID := c.Int("server")
	if discID != 0 && serverID != 0 {
		return c.Help("Cannot cancel a disc and server migration simultaneously")
	}
	if discID != 0 {
		if err := c.Client().CancelDiscMigration(discID); err != nil {
			return err
		}
	} else {
		if err := c.Client().CancelVMMigration(serverID); err != nil {
			return err
		}
	}

	log.Output("Migration cancelled")

	return nil
}
