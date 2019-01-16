package delete

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "backup",
		Usage:       "delete the given backup",
		UsageText:   `delete backup <server name> <disc label> <backup label>`,
		Description: "Deletes the given backup. Backups cannot be recovered after deletion.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "disc",
				Usage: "the disc to delete a backup of",
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to delete a backup from",
				Value: new(flags.VirtualMachineNameFlag),
			},
			cli.StringFlag{
				Name:  "backup",
				Usage: "the name or ID of the backup to delete",
			},
		},
		Action: app.Action(args.Optional("server", "disc", "backup"), with.RequiredFlags("server", "disc", "backup"), with.Auth, func(c *app.Context) (err error) {
			err = c.Client().DeleteBackup(flags.VirtualMachineName(c, "server"), c.String("disc"), c.String("backup"))
			if err != nil {
				return
			}
			c.Log("Backup '%s' deleted successfully", c.String("backup"))
			return
		}),
	})
}
