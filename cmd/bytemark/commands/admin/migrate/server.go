package migrate

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "server",
		Aliases:     []string{"vm"},
		Usage:       "migrate a server to a new head",
		UsageText:   "--admin migrate server <name> [new-head]",
		Description: `This command migrates a server to a new head. If a new head isn't supplied, a new one is picked automatically.`,
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to migrate",
				Value: new(app.VirtualMachineNameFlag),
			},
			cli.StringFlag{
				Name:  "new-head",
				Usage: "the head to move the server to",
			},
		},
		Action: app.Action(args.Optional("server", "new-head"), with.RequiredFlags("server"), with.Auth, func(ctx *app.Context) (err error) {
			vm := ctx.VirtualMachineName("server")
			head := ctx.String("new-head")

			if err := ctx.Client().MigrateVirtualMachine(vm, head); err != nil {
				return err
			}

			log.Outputf("Migration for server %s initiated\n", vm.String())

			return nil
		}),
	})
}
