package update

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "server",
		Usage:       "update a server's configuration",
		UsageText:   "update server [flags] <name>",
		Description: `Updates the configuration of an existing Cloud Server.`,
		Flags: append(app.OutputFlags("server", "object"),
			cli.GenericFlag{
				Name:  "memory",
				Value: new(util.SizeSpecFlag),
				Usage: "How much memory the server will have available, specified in GiB or with GiB/MiB units.",
			},
			cli.StringFlag{
				Name:  "hwprofile",
				Usage: "The hardware profile to use. See `bytemark profiles` for a list of hardware profiles available.",
			},
			cli.GenericFlag{
				Name:  "newname",
				Usage: "A new name for the server",
				Value: new(app.VirtualMachineNameFlag),
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "The server to update",
				Value: new(app.VirtualMachineNameFlag),
			},
		),
		Action: app.Action(args.Optional("newname", "hwprofile", "memory"), with.RequiredFlags("server"), with.Auth, updateServer),
	})
}

func updateServer(c *app.Context) (err error) {
	//FIXME: stuff here
	return nil
}
