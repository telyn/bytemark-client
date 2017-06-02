package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
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
				Action: With(OptionalArgs("disc", "new-grade"), RequiredFlags("disc", "new-grade"), AuthProvider, func(c *Context) (err error) {
					if err := global.Client.RegradeDisc(c.Int("disc"), c.String("new-grade")); err != nil {
						return err
					}

					log.Outputf("Regrade started for disc %d to %s\n", c.Int("disc"), c.String("new-grade"))

					return nil
				}),
			},
		},
	})
}
