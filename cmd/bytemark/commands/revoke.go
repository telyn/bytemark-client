package commands

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "revoke",
		Usage:     "revoke privileges on bytemark self-service objects from other users",
		UsageText: "revoke <privilege> [on] <object> [from] <user>\r\nbytemark grant cluster_admin [to] <user>",
		Action:    cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "privilege",
			Usage:       "revoke privileges on bytemark self-service objects from other users",
			UsageText:   "revoke <privilege> [on] <object> [from] <user>\r\nbytemark grant cluster_admin [to] <user>",
			Description: "Revoke a privilege from a user for a particular bytemark object\r\n\r\n" + privilegeText,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yubikey-required",
					Usage: "Set if the privilege should require a yubikey.",
				},
				cli.GenericFlag{
					Name:  "privilege",
					Usage: "the privilege to revoke",
					Value: new(flags.Privilege),
				},
			},
			Action: app.Action(args.Join("privilege"), with.RequiredFlags("privilege"), with.Privilege("privilege"), func(c *app.Context) (err error) {
				pf := c.PrivilegeFlag("privilege")
				c.Privilege.YubikeyRequired = c.Bool("yubikey-required")

				var privs brain.Privileges
				switch c.Privilege.TargetType() {
				case brain.PrivilegeTargetTypeVM:
					privs, err = c.Client().GetPrivilegesForVirtualMachine(*pf.VirtualMachineName)
					if err != nil {
						return
					}
				case brain.PrivilegeTargetTypeGroup:
					privs, err = c.Client().GetPrivilegesForGroup(*pf.GroupName)
					if err != nil {
						return
					}
				case brain.PrivilegeTargetTypeAccount:
					privs, err = c.Client().GetPrivilegesForAccount(pf.AccountName)
					if err != nil {
						return
					}
				default:
					privs, err = c.Client().GetPrivileges(pf.Username)
					if err != nil {
						return
					}
				}
				i := privs.IndexOf(c.Privilege)
				if i == -1 {
					return fmt.Errorf("Couldn't find such a privilege to revoke")
				}

				err = c.Client().RevokePrivilege(privs[i])
				if err == nil {
					log.Outputf("Revoked %s\r\n", c.PrivilegeFlag("privilege").String())

				}
				return
			}),
		}},
	})
}
