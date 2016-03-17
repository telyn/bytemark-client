package main

import (
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
	"strings"
)

/*
func (cmds *CommandSet) HelpForList() util.ExitCode {
	log.Log("bytemark list")
	log.Log("")
	log.Log("usage: bytemark list servers [group | account]")
	log.Log("       bytemark list groups [account]")
	log.Log("       bytemark list accounts")
	log.Log("       bytemark list keys [user]")
	log.Log("       bytemark list discs <server>")
	return util.E_USAGE_DISPLAYED
}*/

func listDefaultAccountServers() error {
	acc, err := global.Client.GetAccount(global.Config.GetIgnoreErr("account"))
	if err != nil {
		return err
	}
	for _, group := range acc.Groups {
		for _, vm := range group.VirtualMachines {
			log.Output(vm.Hostname)
		}
	}
	return nil
}

func init() {
	commands = append(commands, cli.Command{
		Name: "list",
		Subcommands: []cli.Command{{
			Name: "accounts",
			Action: func(c *cli.Context) {
				accounts, err := global.Client.GetAccounts()

				if err != nil {
					global.Error = err
					return
				}

				for _, group := range accounts {
					log.Output(group.Name)
				}

			},
		}, {
			Name: "discs",
			Action: With(VirtualMachineProvider, func(c *Context) (err error) {
				for _, disc := range c.VirtualMachine.Discs {
					log.Outputf("%s: %dGiB %s\r\n", disc.Label, (disc.Size / 1024), disc.StorageGrade)
				}
				return
			}),
		}, {
			Name: "groups",
			Action: With(AccountProvider, func(c *Context) (err error) {
				for _, group := range c.Account.Groups {
					log.Output(group.Name)
				}
				return
			}),
		}, {
			Name: "keys",
			Action: func(c *cli.Context) {
				username := global.Config.GetIgnoreErr("user")
				if len(c.Args()) == 1 {
					username = c.Args().First()
				}

				err := EnsureAuth()
				if err != nil {
					global.Error = err
					return
				}

				user, err := global.Client.GetUser(username)
				if err != nil {
					global.Error = err
					return
				}

				for _, k := range user.AuthorizedKeys {
					log.Output(k)
				}

			},
		}, {
			Name: "servers",
			Action: func(c *cli.Context) {
				if len(c.Args()) >= 1 {
					nameStr := c.Args().First()
					name := global.Client.ParseGroupName(nameStr)

					group, err := global.Client.GetGroup(name)

					if err != nil {
						if _, ok := err.(lib.NotFoundError); ok {

							if !strings.Contains(nameStr, ".") {
								account, err := global.Client.GetAccount(nameStr)
								if err != nil {
									global.Error = err
									return
								}

								for _, g := range account.Groups {
									for _, vm := range g.VirtualMachines {
										log.Output(vm.Hostname)

									}
								}
								return
							}

						}
						return
					}

					for _, vm := range group.VirtualMachines {
						log.Output(vm.Hostname)
					}
				} else {
					global.Error = listDefaultAccountServers()
					return
				}

			},
		}},
	})
}
