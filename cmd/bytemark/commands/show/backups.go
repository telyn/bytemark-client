package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "backups",
		Usage:       "show all the backups of a server or disc",
		UsageText:   "show backups <server name> [disc label]",
		Description: "Shows all the backups of all the discs in the given server, or if you also give a disc label, just the backups of that disc.",
		Flags: append(app.OutputFlags("backups", "array"),
			cli.StringFlag{
				Name:  "disc",
				Usage: "the disc you wish to list the backups of",
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server you wish to list the backups of",
				Value: new(app.VirtualMachineNameFlag),
			},
		),
		Action: app.Action(args.Optional("server", "disc"), with.RequiredFlags("server", "disc"), with.Auth, func(c *app.Context) (err error) {
			vmName := c.VirtualMachineName("server")
			label := c.String("disc")
			var backups brain.Backups

			if label != "" {
				backups, err = c.Client().GetBackups(vmName, label)
				if err != nil {
					return
				}
			} else {
				err = with.VirtualMachine("server")(c)
				if err != nil {
					return
				}
				for _, disc := range c.VirtualMachine.Discs {
					discbackups, err := c.Client().GetBackups(vmName, disc.Label)
					if err != nil {
						return err
					}
					backups = append(backups, discbackups...)
				}
			}
			return c.OutputInDesiredForm(backups, output.List)
		}),
	})
}
