package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "api keys",
		Aliases:     []string{"apikeys"},
		Usage:       "show all your API keys",
		UsageText:   "show api keys",
		Description: `Shows all API keys for your user.`,
		Flags:       app.OutputFlags("API keys", "array"),
		Action: app.Action(with.Auth, func(ctx *app.Context) error {
			apiKeys, err := brainRequests.GetAPIKeys(ctx.Client())
			if err != nil {
				return err
			}
			user, err := ctx.Client().GetUser(ctx.Client().GetSessionUser())
			if err != nil {
				return err
			}
			myKeys := brain.APIKeys{}
			for _, key := range apiKeys {
				if key.UserID == user.ID {
					myKeys = append(myKeys, key)
				}
			}
			return ctx.OutputInDesiredForm(myKeys, output.List)
		}),
	})
}
