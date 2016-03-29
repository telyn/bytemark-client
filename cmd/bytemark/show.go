package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
	"encoding/json"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "show",
		Description: `Displays information about the given server, group, or account.`,
		Subcommands: []cli.Command{{
			Name:        "account",
			Description: `Displays information about the given account`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output account details as a JSON object",
				},
			},
			Action: With(AccountProvider, func(c *Context) error {
				if c.Bool("json") {
					js, _ := json.MarshalIndent(c.Account, "", "    ")
					log.Output(string(js))
				} else {
					log.Output(util.FormatAccount(c.Account))

					for _, g := range c.Account.Groups {
						log.Outputf("Group %s\r\n", g.Name)
						for _, v := range util.FormatVirtualMachines(g.VirtualMachines) {
							log.Output(v)
						}
					}
				}
				return nil
			}),
		}, {
			Name:  "group",
			Usage: "Outputs info about a group",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output group details as a JSON object",
				},
			},
			Action: With(GroupProvider, func(c *Context) error {
				if c.Bool("json") {
					js, _ := json.MarshalIndent(c.Group, "", "    ")
					log.Output(string(js))
				} else {
					s := ""
					if len(c.Group.VirtualMachines) != 1 {
						s = "s"
					}
					log.Outputf("%s - Group containing %d cloud server%s\r\n", c.Group.Name, len(c.Group.VirtualMachines), s)

					log.Output()

					for _, v := range util.FormatVirtualMachines(c.Group.VirtualMachines) {
						log.Output(v)
					}

				}
				return nil
			}),
		}, {
			Name:  "server",
			Usage: "Displays details about a server",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output server details as a JSON object.",
				},
			},
			Action: With(VirtualMachineProvider, func(c *Context) error {
				if c.Bool("json") {
					js, _ := json.MarshalIndent(c.VirtualMachine, "", "    ")
					log.Output(string(js))
				} else {
					log.Log(util.FormatVirtualMachine(c.VirtualMachine))
				}
				return nil
			}),
		}, {
			Name:  "user",
			Usage: "Displays info about a user. Currently just what SSH keys are authorised for this user",
			Action: With(UserProvider, func(c *Context) error {
				log.Outputf("User %s:\n\nAuthorized keys:\n", c.User.Username)
				for _, k := range c.User.AuthorizedKeys {
					log.Output(k)
				}
				return nil
			}),
		}},
	})
}
