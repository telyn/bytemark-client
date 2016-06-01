package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
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

func listServersInGroup(g *lib.Group) {
	for _, vm := range g.VirtualMachines {
		if vm.Deleted {
			log.Output(vm.Hostname + " (deleted)")
		} else {
			log.Output(vm.Hostname)
		}
	}
}

func init() {
	commands = append(commands, cli.Command{
		Name:        "list",
		Usage:       "scripting-friendly lists of your assets at Bytemark",
		UsageText:   "bytemark list accounts|discs|groups|keys|servers",
		Description: `This commmand will list the kind of object you request, one per line. Perfect for piping into a bash while loop!`,
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "accounts",
			Usage:       "list all the accounts you're able to see",
			UsageText:   "bytemark list accounts",
			Description: `This will list all the accounts that your authentication token has some form of access to.`,
			Action: With(AuthProvider, func(c *Context) error {
				accounts, err := global.Client.GetAccounts()

				if err != nil {
					return err
				}

				for _, group := range accounts {
					log.Output(group.Name)
				}
				return nil
			}),
		}, {
			Name:        "discs",
			Usage:       "list all the discs attached to a given virtual machine",
			UsageText:   "bytemark list discs <virtual machine>",
			Description: `This command lists all the discs attached to the given virtual machine. They're presented in the following format: 'LABEL: SIZE GRADE', where size is an integer number of megabytes. Add the --human flag to output the size in GiB (rounded down to the nearest GiB)`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "human",
					Usage: "output disc size in GiB, suffixed",
				},
			},
			Action: With(VirtualMachineProvider, AuthProvider, func(c *Context) (err error) {
				for _, disc := range c.VirtualMachine.Discs {
					if c.Bool("human") {
						log.Outputf("%s: %dGiB %s\r\n", disc.Label, (disc.Size / 1024), disc.StorageGrade)
					} else {
						log.Outputf("%s: %d %s\r\n", disc.Label, disc.Size, disc.StorageGrade)
					}
				}
				return
			}),
		}, {
			Name:      "groups",
			Usage:     "list all the groups in an account",
			UsageText: "bytemark list groups <account>",
			Description: `This command lists all the groups in the given account, or in your default account.
Your default account is determined by the --account flag, the account variable in your config, falling back to the account with the same name as you log in with.`,
			Action: With(AccountProvider, AuthProvider, func(c *Context) (err error) {
				for _, group := range c.Account.Groups {
					log.Output(group.Name)
				}
				return
			}),
		}, {
			Name:        "keys",
			Usage:       "list all the SSH public keys associated with a user",
			UsageText:   "bytemark list keys [user]",
			Description: "Lists all the SSH public keys associated with a user, defaulting to your log-in user.",
			Action: func(c *cli.Context) error {
				username := global.Config.GetIgnoreErr("user")
				if len(c.Args()) == 1 {
					username = c.Args().First()
				}

				err := EnsureAuth()
				if err != nil {
					return err
				}

				user, err := global.Client.GetUser(username)
				if err != nil {
					return err
				}

				for _, k := range user.AuthorizedKeys {
					log.Output(k)
				}

				return nil
			},
		}, {
			Name:      "servers",
			Usage:     "list all the servers in an account",
			UsageText: "bytemark list servers [account]",
			Description: `This command lists all the servers in the given account, or in your default account if you didn't specify an account on the command-line.
Deleted servers are included in the list, with ' (deleted)' appended.`,
			Action: With(AuthProvider, func(c *Context) error {
				var account *lib.Account
				var err error

				if len(c.Args()) >= 1 {
					nameStr, _ := c.NextArg()
					name := global.Client.ParseAccountName(nameStr, global.Config.GetIgnoreErr("account"))
					account, err = global.Client.GetAccount(name)
				} else {
					account, err = global.Client.GetDefaultAccount()
				}

				if err != nil {
					return err
				}
				for _, g := range account.Groups {
					listServersInGroup(g)
				}
				return nil
			}),
		}},
	})
}
