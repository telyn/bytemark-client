package delete

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:    "api key",
		Aliases: []string{"apikey"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "api-key",
				Usage: "api key label or ID to delete",
			},
		},
		Action: app.Action(args.Optional("api-key"), with.RequiredFlags("api-key"), with.Auth, func(ctx *app.Context) error {
			err := brainRequests.DeleteAPIKey(ctx.Client(), ctx.String("api-key"))
			if err == nil {
				ctx.Log("Successfully deleted %s", ctx.String("api-key"))
			}
			return err
		}),
	})
}
