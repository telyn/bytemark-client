package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "keys",
		Usage:       "shows all the SSH public keys associated with a user",
		UsageText:   "show keys [user]",
		Description: "Shows all the SSH public keys associated with a user, defaulting to your log-in user.",
		Action: app.Action(args.Optional("user"), with.User("user"), func(c *app.Context) error {
			// TODO(telyn): could this be rewritten using OutputInDesiredForm / is it desirable to?
			for _, k := range c.User.AuthorizedKeys {
				log.Output(k)
			}

			return nil
		}),
	})
}
