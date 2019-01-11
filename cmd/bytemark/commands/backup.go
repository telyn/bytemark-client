package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "backup",
		Usage:       "create a backup - see `bytemark help backup <kind of thing> `",
		UsageText:   "backup disc",
		Description: "create a backup",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "disc",
			Usage:       "create a backup of a disc",
			UsageText:   "backup disc <server> <disc label>",
			Description: `create a backup of the disc's current state. The backup is moved to another tail in the "iceberg" storage grade.`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "disc",
					Usage: "the disc to create a backup of",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server whose disk you wish to backup",
					Value: new(flags.VirtualMachineNameFlag),
				},
			},
			Action: app.Action(args.Optional("server", "disc"), with.RequiredFlags("server", "disc"), with.Auth, func(c *app.Context) error {
				backup, err := c.Client().CreateBackup(flags.VirtualMachineName(c, "server"), c.String("disc"))
				if err != nil {
					return err
				}
				log.Errorf("Backup '%s' taken successfully!", backup.Label)
				return nil
			}),
		},
		},
	})
}
