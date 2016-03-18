package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
	/*"encoding/json"*/
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name: "show",
		Description: `Displays information about the given server, group, or account.
If the --verbose flag is given to bytemark show group or bytemark show account, full details are given for each server.`,
		Subcommands: []cli.Command{{
			Name: "account",
			Action: With(AccountProvider, func(c *Context) error {
				/*if *jsonOut {
					js, _ := json.MarshalIndent(acc, "", "    ")
					log.Output(string(js))
				} else {*/
				log.Output(util.FormatAccount(c.Account))

				for _, g := range c.Account.Groups {
					log.Outputf("Group %s\r\n", g.Name)
					for _, v := range util.FormatVirtualMachines(g.VirtualMachines) {
						log.Output(v)
					}
				}
				/*}*/
				return nil
			}),
		}, {
			Name: "group",
			Action: With(GroupProvider, func(c *Context) error {
				/*if *jsonOut {
					js, _ := json.MarshalIndent(group, "", "    ")
					log.Output(string(js))
				} else {*/
				s := ""
				if len(c.Group.VirtualMachines) != 1 {
					s = "s"
				}
				log.Outputf("%s - Group containing %d cloud server%s\r\n", c.Group.Name, len(c.Group.VirtualMachines), s)

				log.Output()

				for _, v := range util.FormatVirtualMachines(c.Group.VirtualMachines) {
					log.Output(v)
				}
				return nil

				/*}*/
			}),
		}, {
			Name: "server",
			Action: With(VirtualMachineProvider, func(c *Context) error {
				/*if *jsonOut {
					js, _ := json.MarshalIndent(vm, "", "    ")
					log.Output(string(js))
				} else {*/
				log.Log(util.FormatVirtualMachine(c.VirtualMachine))
				/*}*/
				return nil
			}),
		}, {
			Name: "user",
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
