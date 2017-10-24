package main

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "resize",
		Usage:     "resize a cloud server's disc",
		UsageText: "bytemark resize disc <server> <disc label> <size>",
		Description: `resize a cloud server's disc

Resizes the given disc to the given size. Sizes may be specified with a + in front, in which case they are interpreted as relative. For example, '+2GB' is parsed as 'increase the disc size by 2GiB', where '2GB' is parsed as 'set the size of the disc to 2GiB'`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "disc",
			Aliases:     []string{"disk"},
			Usage:       "resize a cloud server's disc",
			UsageText:   "bytemark resize disc <server> <disc label> <size>",
			Description: "Resizes the given server's disc to the given size. Sizes may be specified with a + in front, in which case they are interpreted as relative. For example, '+2GB' is parsed as 'increase the disc size by 2GiB', where '2GB' is parsed as 'set the size of the disc to 2GiB'",
			Flags: []cli.Flag{
				forceFlag,
				cli.StringFlag{
					Name:  "disc",
					Usage: "the disc to resize",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server that the disc is attached to",
					Value: new(app.VirtualMachineNameFlag),
				},
				cli.GenericFlag{
					Name:  "new-size",
					Usage: "the new size for the disc. Prefix with + to indicate 'increase by'",
					Value: new(app.ResizeFlag),
				},
			},
			Action: app.Action(args.Optional("server", "disc", "new-size"), with.RequiredFlags("server", "disc", "new-size"), with.Disc("server", "disc"), func(c *app.Context) (err error) {
				vmName := c.VirtualMachineName("server")
				size := c.ResizeFlag("new-size")
				newSize := size.Size

				if size.Mode == app.ResizeModeIncrease {
					newSize += c.Disc.Size
				}
				log.Logf("Resizing %s from %dGiB to %dGiB...", c.Disc.Label, c.Disc.Size/1024, newSize/1024)

				if !c.Bool("force") && !util.PromptYesNo(fmt.Sprintf("Are you certain you wish to perform this resize?")) {
					return util.UserRequestedExit{}
				}

				err = c.Client().ResizeDisc(vmName, c.String("disc"), newSize)
				if err != nil {
					log.Logf("Failed!\r\n")
					return
				}
				log.Logf("Completed.\r\n")
				return
			}),
		}},
	})
}
