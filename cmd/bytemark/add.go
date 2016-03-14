package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"strings"
)

func init() {

	commands = append(commands, cli.Command{
		Name:            "add",
		SkipFlagParsing: true,
		Subcommands: []cli.Command{{
			Name:  "key",
			Usage: "bytemark add key [--public-key-file=<filename>] [--user=<user>] [<public key>]",
			UsageText: `Add the given public key to the given user (or the default user). This will allow them
to use that key to access management IPs they have access to using that key.
To remove a key, use the remove key command.`,
			Description: `Specify --public-key-file= to read the public key from stdin
--public-key-file will be ignored if a public key is specified in the arguments`,
			SkipFlagParsing: true,
			Action: func(ctx *cli.Context) {
				flags := util.MakeCommonFlagSet()
				keyFile := flags.String("public-key-file", "", "")

				flags.Parse(ctx.Args())
				args := global.Config.ImportFlags(flags)

				user := global.Config.GetIgnoreErr("user")

				key := strings.TrimSpace(strings.Join(args, " "))
				if key == "" {
					if *keyFile == "" {
						cli.ShowCommandHelp(ctx, ctx.Command.Name)
						global.Error = &PEBKACError{}
						return
					}

					keyBytes, err := ioutil.ReadFile(*keyFile)
					if err != nil {
						global.Error = err
						return
					}
					key = string(keyBytes)
				}

				err := EnsureAuth()
				if err != nil {
					global.Error = err
					return
				}

				err = global.Client.AddUserAuthorizedKey(user, key)
				if err == nil {
					log.Log("Key added successfully")
					return
				} else {
					global.Error = err
					return
				}
			},
		}},
	})
}
