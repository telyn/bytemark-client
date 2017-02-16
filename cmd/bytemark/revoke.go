package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:        "revoke",
		Usage:       "revoke privileges on bytemark self-service objects from other users",
		UsageText:   "bytemark revoke <privilege> [on] <object> [from|to] <user>\r\nbytemark grant cluster_admin [to] <user>",
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
			pf := c.PrivilegeFlag("privilege")
			c.Privilege.YubikeyRequired = c.Bool("yubikey-required")

			var privs brain.Privileges
			switch strings.SplitN(string(c.Privilege.Level), "_", 2)[0] {
			case "vm":
				privs, err = global.Client.GetPrivilegesForVirtualMachine(*pf.VirtualMachineName)
				if err != nil {
					return
				}
			case "group":
				privs, err = global.Client.GetPrivilegesForGroup(*pf.GroupName)
				if err != nil {
					return
				}
			case "account":
				privs, err = global.Client.GetPrivilegesForAccount(pf.AccountName)
				if err != nil {
					return
				}
			default:
				privs, err = global.Client.GetPrivileges(pf.Username)
				if err != nil {
					return
				}
			}
			i := privs.IndexOf(c.Privilege)
			if i == -1 {
				return fmt.Errorf("Couldn't find such a privilege to revoke")
			}

			err = global.Client.RevokePrivilege(*privs[i])
			if err == nil {
				log.Outputf("Revoked %s\r\n", c.PrivilegeFlag("privilege"))

			}
			return
		}),
	})
}
