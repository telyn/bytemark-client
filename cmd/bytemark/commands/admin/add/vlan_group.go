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
		Name:      "vlan group",
		Aliases:   []string{"vlan-group"},
		Usage:     "adds groups for private VLANs",
		UsageText: "--admin add vlan group <group> [vlan-num]",
		Description: `Add a group in the specified account, with an optional VLAN specified.

Used when setting up a private VLAN for a customer.`,
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "group",
				Usage: "the name of the group to add",
				Value: new(app.GroupNameFlag),
			},
			cli.IntFlag{
				Name:  "vlan-num",
				Usage: "The VLAN number to add the group to",
			},
		},
		Action: app.Action(args.Optional("group", "vlan-num"), with.RequiredFlags("group"), with.Auth, func(c *app.Context) error {
			gp := c.GroupName("group")
			if err := c.Client().AdminCreateGroup(gp, c.Int("vlan-num")); err != nil {
				return err
			}
			log.Logf("Group %s was added under account %s\r\n", gp.Group, gp.Account)
			return nil
		}),
	})
}
