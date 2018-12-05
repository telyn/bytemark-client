package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "discs",
		Usage:       "show all the discs attached to a given virtual machine",
		UsageText:   "show discs <virtual machine>",
		Description: `This command shows all the discs attached to the given virtual machine. They're presented in the following format: 'LABEL: SIZE GRADE', where size is an integer number of megabytes. Add the --human flag to output the size in GiB (rounded down to the nearest GiB)`,
		Flags: append(app.OutputFlags("discs", "array"),
			cli.BoolFlag{
				Name:  "human",
				Usage: "output disc size in GiB, suffixed",
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server whose discs you wish to list",
				Value: new(flags.VirtualMachineName),
			},
		),
		Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.VirtualMachine("server"), func(c *app.Context) error {
			return c.OutputInDesiredForm(c.VirtualMachine.Discs, output.List)
		}),
	})
}
