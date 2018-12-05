package add

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"

	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "group",
		Usage:       "add a group for organising your servers",
		UsageText:   "add group <group name>",
		Description: `Groups are part of your server's fqdn`,
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "group",
				Usage: "the name of the group to add",
				Value: new(flags.GroupNameFlag),
			},
		},
		Action: app.Action(args.Optional("group"), with.RequiredFlags("group"), with.Auth, createGroup),
	})
}

func createGroup(c *app.Context) (err error) {
	gp := flags.GroupName(c, "group")
	err = c.Client().CreateGroup(gp)
	if err == nil {
		log.Logf("Group %s was created under account %s\r\n", gp.Group, gp.Account)
	}
	return
}
