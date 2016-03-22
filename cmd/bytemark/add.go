package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
	"strings"
)

func init() {
	publicKeyFile := util.FileFlag{}
	commands = append(commands, cli.Command{
		Name: "add",
		Subcommands: []cli.Command{{
			Name:        "key",
			Usage:       "Add public SSH keys to a Bytemark user",
			Description: `Add the given public key to the given user (or the default user). This will allow them to use that key to access management IPs they have access to using that key. To remove a key, use the remove key command. --public-key-file will be ignored if a public key is specified in the arguments`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "Which user to add the key to. Defaults to the username you log in as.",
				},
				cli.GenericFlag{
					Name:  "public-key-file",
					Usage: "The public key file to add to the account",
					Value: &publicKeyFile,
				},
			},
			Action: With(AuthProvider, func(ctx *Context) (err error) {
				user := global.Config.GetIgnoreErr("user")

				key := strings.TrimSpace(strings.Join(ctx.Args(), " "))
				if key == "" {
					key = publicKeyFile.Value
				}

				err = global.Client.AddUserAuthorizedKey(user, key)
				if err == nil {
					log.Log("Key added successfully")
					return
				} else {
					return
				}
			}),
		}},
	})
}
