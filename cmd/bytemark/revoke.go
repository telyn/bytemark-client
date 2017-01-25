package main

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "revoke",
		Usage:       "revoke privileges on bytemark self-service objects from other users",
		UsageText:   "bytemark revoke <privilege> [on] <object> [from|to] <user>\r\nbytemark grant cluster_admin [to] <user>",
		Description: "Revoke a privilege from a user for a particular bytemark object\r\n\r\n" + privilegeText,
		Action: With(func(c *Context) (err error) {
			priv, _, err := parsePrivilege(c)

			if err != nil {
				return
			}

			err = global.Client.RevokePrivilege(priv)
			if err == nil {
				log.Outputf("Revoked %s from %s\r\n", priv.Level, priv.Username)

			}
			return
		}),
	})
}
