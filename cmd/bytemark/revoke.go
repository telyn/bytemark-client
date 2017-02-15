package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:        "revoke",
		Usage:       "revoke privileges on bytemark self-service objects from other users",
		UsageText:   "bytemark revoke <privilege> [on] <object> [from] <user>\r\nbytemark grant cluster_admin [to] <user>",
		Description: "Revoke a privilege from a user for a particular bytemark object\r\n\r\n" + privilegeText,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "yubikey-required",
				Usage: "Set if the privilege should require a yubikey.",
			},
			cli.GenericFlag{
				Name:  "privilege",
				Usage: "the privilege to revoke",
				Value: new(PrivilegeFlag),
			},
		},
		Action: With(JoinArgs("privilege"), RequiredFlags("privilege"), PrivilegeProvider("privilege"), func(c *Context) (err error) {
			c.Privilege.YubikeyRequired = c.Bool("yubikey-required")

			err = global.Client.RevokePrivilege(c.Privilege)
			if err == nil {
				log.Outputf("Revoked %s\r\n", c.PrivilegeFlag("privilege"))

			}
			return
		}),
	})
}
