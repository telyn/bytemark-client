package update

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "disc",
		Aliases:   []string{"disk"},
		Usage:     "move or resize a cloud server's disc",
		UsageText: "update disc <server> <disc label> [--new-size <size>] [--new-server <server>]",
		Description: `resize a cloud server's disc or move it to another server

Resizes the given disc to the given size. Sizes may be specified with a + in front, in which case they are interpreted as relative. For example, '+2GB' is parsed as 'increase the disc size by 2GiB', where '2GB' is parsed as 'set the size of the disc to 2GiB'

Moving the disc to another server may require you to update your operating system configuration. Both servers must be shutdown and root discs cannot be moved. Please find documentation for moving discs at https://docs.bytemark.co.uk/`,
		Flags: []cli.Flag{
			flags.Force,
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
			cli.GenericFlag{
				Name:  "new-server",
				Usage: "the server that the disc should be moved to",
				Value: new(app.VirtualMachineNameFlag),
			},
		},
		Action: app.Action(args.Optional("server", "disc", "new-size", "new-server"), with.RequiredFlags("server", "disc"), with.Disc("server", "disc"), updateDisc),
	})
}

func updateDisc(c *app.Context) (err error) {
	vmName := c.VirtualMachineName("server")

	if c.IsSet("new-size") {
		size := c.ResizeFlag("new-size")
		newSize := size.Size

		if size.Mode == app.ResizeModeIncrease {
			newSize += c.Disc.Size
		}
		log.Logf("Resizing %s from %dGiB to %dGiB...", c.Disc.Label, c.Disc.Size/1024, newSize/1024)

		if !flags.Forced(c) && !util.PromptYesNo(c.Prompter(), fmt.Sprintf("Are you certain you wish to perform this resize?")) {
			return util.UserRequestedExit{}
		}

		err = c.Client().ResizeDisc(vmName, c.String("disc"), newSize)
		if err != nil {
			log.Logf("Failed!\r\n")
			return
		}
		log.Logf("Completed.\r\n")
	}

	if c.IsSet("new-server") {
		newVM := c.VirtualMachineName("new-server")

		log.Logf("This may require an update to the operating system configuration, please find documentation for moving discs at https://docs.bytemark.co.uk/\r\n")

		if !flags.Forced(c) && !util.PromptYesNo(c.Prompter(), fmt.Sprintf("Are you certain you wish to move the disc?")) {
			return util.UserRequestedExit{}
		}
		log.Logf("Moving %s from %s to %s...", c.Disc.Label, vmName, newVM)
		err = brainRequests.MoveDisc(c.Client(), vmName, c.String("disc"), newVM)
		if err != nil {
			log.Logf("Failed!\r\n")
			return
		}
		log.Logf("Completed.\r\n")
	}

	return
}
