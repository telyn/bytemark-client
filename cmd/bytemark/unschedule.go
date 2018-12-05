package main

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "unschedule",
		Usage:     "unschedule automated backups",
		UsageText: "unschedule backups <server> <disc> <schedule id>",
		Description: `unschedules automated backups so that they are no longer taken
	
The <schedule id> is a number that can be found out using 'bytemark show disc <server> <disc>'
`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "backups",
				Usage:     "unschedule automated backups",
				UsageText: "unschedule backups <server> <disc> <schedule id>",
				Description: `unschedules automated backups so that they are no longer taken
	
The <schedule id> is a number that can be found out using 'bytemark show disc <server> <disc>'
`,
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "server",
						Usage: "the server to unschedule backups on",
						Value: new(flags.VirtualMachineName),
					},
					cli.StringFlag{
						Name:  "disc",
						Usage: "the disc to unschedule some backups of",
					},
					cli.IntFlag{
						Name:  "schedule-id",
						Usage: "the ID of the schedule to remove. See the output of `show disc` to find out schedule IDs.",
					},
				},
				Action: app.Action(args.Optional("server", "disc", "schedule-id"), with.RequiredFlags("server", "disc", "schedule-id"), with.Auth, func(c *app.Context) (err error) {
					if c.Int("schedule-id") < 1 {
						return fmt.Errorf("schedule-id not specified or invalid")
					}
					vmName := c.VirtualMachineName("server")
					err = c.Client().DeleteBackupSchedule(vmName, c.String("disc"), c.Int("schedule-id"))
					if err == nil {
						log.Log("Backups unscheduled.")
					}
					return
				}),
			},
		},
	})
}
