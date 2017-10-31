package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "regrade",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "disc",
				Usage:     "regrade a disc",
				UsageText: "bytemark --admin regrade disc <disc> [--new-grade]",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "disc",
						Usage: "the ID of the disc to regrade",
					},
					cli.StringFlag{
						Name:  "new-grade",
						Usage: "the new grade of the disc",
					},
				},
				Action: app.Action(args.Optional("disc", "new-grade"), with.RequiredFlags("disc", "new-grade"), with.Auth, func(c *app.Context) (err error) {
					if err := c.Client().RegradeDisc(c.Int("disc"), c.String("new-grade")); err != nil {
						return err
					}

					c.Log("Regrade started for disc %d to %s\n", c.Int("disc"), c.String("new-grade"))

					return nil
				}),
			},
		},
	})
}
