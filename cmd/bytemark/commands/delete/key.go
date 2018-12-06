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
		Name:        "key",
		Usage:       "deletes the specified key",
		UsageText:   "delete key [--user <user>] <key>",
		Description: "Keys may be specified as just the comment part or as the whole key. If there are multiple keys with the comment given, an error will be returned",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "user",
				Usage: "Which user to delete the key from. Defaults to the username you log in as.",
			},
			cli.StringFlag{
				Name:  "public-key",
				Usage: "The public key to delete. Can be the comment part or the whole public key",
			},
		},
		Action: app.Action(args.Join("public-key"), with.RequiredFlags("public-key"), with.Auth, func(c *app.Context) (err error) {
			user := c.String("user")
			if user == "" {
				user = c.Config().GetIgnoreErr("user")
			}

			key := c.String("public-key")
			if key == "" {
				return c.Help("You must specify a key to delete.\r\n")
			}

			err = brainRequests.DeleteUserAuthorizedKey(c.Client(), user, key)
			if err == nil {
				c.Log("Key deleted successfully")
			}
			return
		}),
	})
}
