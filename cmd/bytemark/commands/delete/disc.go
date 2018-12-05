package delete

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flagsets"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "disc",
		Usage:       "delete the given disc",
		UsageText:   "delete disc [--server <virtual machine name> --label <disc label>] | [--id <disc ID>]",
		Description: "Deletes the given disc. To find out a disc's label you can use the `bytemark show server` command or `bytemark list discs` command.",
		Flags: []cli.Flag{
			flagsets.Force,
			cli.StringFlag{
				Name:  "label",
				Usage: "the disc to delete, must provide a server too",
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server whose disc you wish to delete, must provide a label too",
				Value: new(flags.VirtualMachineNameFlag),
			},
			cli.StringFlag{
				Name:  "id",
				Usage: "the ID of the disc to delete",
			},
		},
		Aliases: []string{"disk"},
		Action: app.Action(args.Optional("server", "label", "id"), with.Auth, func(c *app.Context) (err error) {
			if !c.Bool("force") && !util.PromptYesNo(c.Prompter(), "Are you sure you wish to delete this disc? It is impossible to recover.") {
				return util.UserRequestedExit{}
			}
			vmName := flags.VirtualMachineName(c, "server")
			discLabel := c.String("label")
			discID := c.String("id")

			if discID != "" {
				return brainRequests.DeleteDiscByID(c.Client(), discID)
			} else if vmName.String() != "" && discLabel != "" {
				return c.Client().DeleteDisc(vmName, discLabel)
			} else {
				return fmt.Errorf("Please include both --server and --label flags or provide --id")
			}
		}),
	})
}
