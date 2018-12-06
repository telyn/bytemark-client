package add

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "user",
		Usage:       "adds a new cluster admin or cluster superuser",
		UsageText:   "--admin add user <username> <privilege>",
		Description: `adds a new cluster admin or superuser. The privilege field must be either cluster_admin or cluster_su.`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "username",
				Usage: "The username of the new user",
			},
			cli.StringFlag{
				Name:  "privilege",
				Usage: "The privilege to grant to the new user",
			},
		},
		Action: app.Action(args.Optional("username", "privilege"), with.RequiredFlags("username", "privilege"), with.Auth, func(c *app.Context) error {
			// Privilege is just a string and not a flags.Privilege, since it can only be "cluster_admin" or "cluster_su"
			if err := c.Client().CreateUser(c.String("username"), c.String("privilege")); err != nil {
				return err
			}
			log.Logf("User %s has been added with %s privileges\r\n", c.String("username"), c.String("privilege"))
			return nil
		}),
	})
}
