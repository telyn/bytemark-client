package admin

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:   "create",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "user",
				Usage:     "creates a new cluster admin or cluster superuser",
				UsageText: "bytemark --admin create user <username> <privilege>",
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
					// Privilege is just a string and not a app.PrivilegeFlag, since it can only be "cluster_admin" or "cluster_su"
					if err := c.Client().CreateUser(c.String("username"), c.String("privilege")); err != nil {
						return err
					}
					log.Logf("User %s has been created with %s privileges\r\n", c.String("username"), c.String("privilege"))
					return nil
				}),
			},
			{
				Name:      "vlan-group",
				Usage:     "creates groups for private VLANs",
				UsageText: "bytemark --admin create vlan-group <group> [vlan-num]",
				Description: `Create a group in the specified account, with an optional VLAN specified.

Used when setting up a private VLAN for a customer.`,
				Flags: []cli.Flag{
					cli.GenericFlag{
						Name:  "group",
						Usage: "the name of the group to create",
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
					log.Logf("Group %s was created under account %s\r\n", gp.Group, gp.Account)
					return nil
				}),
			},
			{
				Name:      "ip range",
				Usage:     "create a new IP range in a VLAN",
				UsageText: "bytemark --admin create ip range <ip-range> <vlan-num>",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "ip-range",
						Usage: "the IP range to add",
					},
					cli.IntFlag{
						Name:  "vlan-num",
						Usage: "The VLAN number to add the IP range to",
					},
				},
				Action: app.Action(args.Optional("ip-range", "vlan-num"), with.RequiredFlags("ip-range", "vlan-num"), with.Auth, func(c *app.Context) error {
					if err := c.Client().CreateIPRange(c.String("ip-range"), c.Int("vlan-num")); err != nil {
						return err
					}
					log.Logf("IP range created\r\n")
					return nil
				}),
			},
		},
	})
}
