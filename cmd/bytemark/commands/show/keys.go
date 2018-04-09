package show

import (
	"encoding/json"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "keys",
		Usage:       "shows all the SSH public keys associated with a user",
		UsageText:   "show keys [user]",
		Description: "Shows all the SSH public keys associated with a user, defaulting to your log-in user.",
		Flags: append([]cli.Flag{
			cli.StringFlag{
				Name:  "user",
				Usage: "the user whose keys you wish to see",
			},
		}, app.OutputFlags("keys", "array")...),
		Action: app.Action(args.Optional("user"), with.User("user"), func(ctx *app.Context) error {
			keys := ctx.User.AuthorizedKeys
			if f, _ := ctx.OutputFormat(); f == "json" {
				encoder := json.NewEncoder(ctx.Writer())
				encoder.SetIndent("", "    ")
				return encoder.Encode(keys.Strings())
			}
			return ctx.OutputInDesiredForm(keys)
		}),
	})
}
