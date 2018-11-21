package commands

import (
	"errors"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

const privilegeText = `PRIVILEGES

Privileges are listed in descending order of privilege. Users can only create or modify privileges which have a lower level than themselves. For example, an account_admin can only create or modify group_admin and lower privileges. This means that to grant or revoke another use account_admin on an account you have account_admin on, you must contact support.

cluster_admin - full administration rights to the whole cluster
account_admin - full administration rights to a whole account
group_admin - create, modify & delete servers within a group
vm_admin - modify & administer a single server
vm_console - access to a server's console`

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "grant",
		Usage:     "grant privileges on bytemark self-service objects to other users",
		UsageText: "grant privilege <privilege> [on] <object> [to] <user>\r\nbytemark grant cluster_admin [to] <user>",
		Action:    cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "privilege",
			Usage:       "grant privileges on bytemark self-service objects to other users",
			UsageText:   "grant privilege <privilege> [on] <object> [to] <user>\r\nbytemark grant cluster_admin [to] <user>",
			Description: "Grant a privilege to a user for a particular bytemark object\r\n\r\n" + privilegeText,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "api-key-id",
					Usage: "ID of an API key that should be allowed to use this privilege. Leave blank when not using api keys",
				},
				cli.BoolFlag{
					Name:  "yubikey-required",
					Usage: "Set if the privilege should require a yubikey.",
				},
				cli.GenericFlag{
					Name:  "privilege",
					Usage: "A privilege written out like '<level> [on] <object> [to] <user>",
					Value: new(app.PrivilegeFlag),
				},
			},
			Action: app.Action(args.Join("privilege"), with.RequiredFlags("privilege"), with.Privilege("privilege"), func(c *app.Context) (err error) {
				c.Privilege.YubikeyRequired = c.Bool("yubikey-required")
				if c.Bool("yubikey-required") && c.IsSet("api-key-id") {
					return errors.New("Only one of --api-key-id and --yubikey-required may be set at a time")
				}
				c.Privilege.APIKeyID = c.Int("api-key-id")

				err = c.Client().GrantPrivilege(c.Privilege)
				if err == nil {
					log.Outputf("Granted %s\r\n", c.PrivilegeFlag("privilege").String())
				}
				return
			}),
		}},
	})
}
