package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strconv"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "schedule",
		Usage:     "schedule backups to occur at a regular frequency",
		UsageText: "bytemark schedule backups [--start <date>] <server> <disc> <interval>",
		Description: `schedule backups to occur at a regular interval (defined in seconds)
		
EXAMPLES

To have daily backups at midnight of a server called 'fileservers' 'very-important-data' disc:
bytemark schedule backups --start 00:00 fileserver very-important-data 86400

To have hourly backups starting at 14:37 (Central European Summer Time) on the 5th of April, 2017:
bytemark schedule backups --start "2017-04-05T14:37:00+02:00" fileserver very-important-data 3600`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "backups",
				Usage:     "schedule backups to occur at a regular frequency",
				UsageText: "bytemark schedule backups [--start <date>] <server> <disc> <interval>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "start",
						Usage: "date & time the schedule starts. Assumes BST/GMT (depending on time of year) if not specified - defaults to 00:00",
					},
					cli.StringFlag{
						Name:  "disc",
						Usage: "the disc to schedule backups of",
					},
				},
				Description: `schedule backups to occur at a regular interval (defined in seconds)
		
EXAMPLES

To have daily backups at midnight of a server called 'fileservers' 'very-important-data' disc:
bytemark schedule backups --start 00:00 fileserver very-important-data 86400

To have hourly backups starting at 14:37 (Central European Summer Time) on the 5th of April, 2017:
bytemark schedule backups --start "2017-04-05T14:37:00+02:00" fileserver very-important-data 3600`,
				Action: With(VirtualMachineNameProvider, OptionalArgs("disc"), func(c *Context) (err error) {
					start := c.String("start")
					if start == "" {
						start = "00:00"
					}

					intervalStr, err := c.NextArg()
					if err != nil {
						return
					}

					interval, err := strconv.Atoi(intervalStr)
					if err != nil {
						return
					}

					_, err = global.Client.CreateBackupSchedule(*c.VirtualMachineName, c.String("disc"), start, interval)
					if err == nil {
						log.Log("Schedule set.")
					}
					return
				}),
			},
		},
	})
}
