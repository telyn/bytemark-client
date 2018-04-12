package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "reify",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "disc",
				Usage:     "reify a disc",
				UsageText: "--admin reify disc <disc>",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "disc",
						Usage: "the ID of the disc to reify",
					},
				},
				Action: app.Action(args.Optional("disc"), with.RequiredFlags("disc"), with.Auth, func(c *app.Context) (err error) {
					if err := c.Client().ReifyDisc(c.Int("disc")); err != nil {
						return err
					}

					c.Log("Reification started for disc %d\n", c.Int("disc"))

					return nil
				}),
			},
		},
	})
}
