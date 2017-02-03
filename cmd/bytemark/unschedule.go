package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strconv"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "unschedule",
		Usage:     "unschedule automated backups",
		UsageText: "bytemark unschedule backups <server> <disc> <schedule id>",
		Description: `unschedules automated backups so that they are no longer taken
	
The <schedule id> is a number that can be found out using 'bytemark show disc <server> <disc>'
`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "backups",
				Usage:     "unschedule automated backups",
				UsageText: "bytemark unschedule backups <server> <disc> <schedule id>",
				Description: `unschedules automated backups so that they are no longer taken
	
The <schedule id> is a number that can be found out using 'bytemark show disc <server> <disc>'
`,
				Action: With(VirtualMachineNameProvider, DiscLabelProvider, func(c *Context) (err error) {
					idStr, err := c.NextArg()
					if err != nil {
						return
					}

					id, err := strconv.Atoi(idStr)
					if err != nil {
						return
					}

					err = global.Client.DeleteBackupSchedule(*c.VirtualMachineName, *c.DiscLabel, id)
					if err == nil {
						log.Log("Backups unscheduled.")
					}
					return
				}),
			},
		},
	})
}
