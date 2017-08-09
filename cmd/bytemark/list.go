package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func listServersInGroup(g brain.Group) {
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
		Name:      "list",
		Usage:     "scripting-friendly lists of your assets at Bytemark",
		UsageText: "bytemark list accounts|discs|groups|keys|servers",
		Description: `scripting-friendly lists of your assets at Bytemark

This commmand will list the kind of object you request, one per line. Perfect for piping into a bash while loop!`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "accounts",
			Usage:       "list all the accounts you're able to see",
			UsageText:   "bytemark list accounts",
			Description: `This will list all the accounts that your authentication token has some form of access to.`,
			Flags:       OutputFlags("accounts", "array"),
			Action: With(AuthProvider, func(c *Context) error {
				accounts, err := global.Client.GetAccounts()

				if err != nil {
					return err
				}
				return c.OutputInDesiredForm(accounts, "list")
			}),
		}, {
			Name:        "discs",
			Usage:       "list all the discs attached to a given virtual machine",
			UsageText:   "bytemark list discs <virtual machine>",
			Description: `This command lists all the discs attached to the given virtual machine. They're presented in the following format: 'LABEL: SIZE GRADE', where size is an integer number of megabytes. Add the --human flag to output the size in GiB (rounded down to the nearest GiB)`,
			Flags: append(OutputFlags("discs", "array"),
				cli.BoolFlag{
					Name:  "human",
					Usage: "output disc size in GiB, suffixed",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server whose discs you wish to list",
					Value: new(VirtualMachineNameFlag),
				},
			),
			Action: With(OptionalArgs("server"), RequiredFlags("server"), VirtualMachineProvider("server"), func(c *Context) error {
				return c.OutputInDesiredForm(c.VirtualMachine.Discs, "list")
			}),
		}, {
			Name:        "groups",
			Usage:       "list all the groups in an account",
			UsageText:   "bytemark list groups [account]",
			Description: `This command lists all the groups in the given account, or in your default account if not specified.`,
			Flags: append(OutputFlags("groups", "array"),
				cli.GenericFlag{
					Name:  "account",
					Usage: "the account to list the groups of",
					Value: new(AccountNameFlag),
				},
			),
			Action: With(OptionalArgs("account"), RequiredFlags("account"), AccountProvider("account"), func(c *Context) error {
				return c.OutputInDesiredForm(c.Account.Groups, "list")
			}),
		}, {
			Name:        "keys",
			Usage:       "list all the SSH public keys associated with a user",
			UsageText:   "bytemark list keys [user]",
			Description: "Lists all the SSH public keys associated with a user, defaulting to your log-in user.",
			Action: With(OptionalArgs("user"), UserProvider("user"), func(c *Context) error {
				// TODO(telyn): could this be rewritten using OutputInDesiredForm / is it desirable to?
				for _, k := range c.User.AuthorizedKeys {
					log.Output(k)
				}

				return nil
			}),
		}, {
			Name:      "servers",
			Usage:     "list all the servers in an account",
			UsageText: "bytemark list servers [account]",
			Description: `This command lists all the servers in the given account, or in your default account if not specified.
Deleted servers are included in the list, with ' (deleted)' appended.`,
			Flags: append(OutputFlags("servers", "array"),
				cli.GenericFlag{
					Name:  "account",
					Usage: "the account to list the servers of",
					Value: new(AccountNameFlag),
				},
			),
			Action: With(OptionalArgs("account"), AccountProvider("account"), AuthProvider, func(c *Context) error {
				servers := brain.VirtualMachines(make([]brain.VirtualMachine, 0))

				for _, g := range c.Account.Groups {
					servers = append(servers, g.VirtualMachines...)
				}
				return c.OutputInDesiredForm(servers)
			}),
		}, {
			Name:        "backups",
			Usage:       "list all the backups of a server or disc",
			UsageText:   "bytemark list backups <server name> [disc label]",
			Description: "Lists all the backups of all the discs in the given server, or if you also give a disc label, just the backups of that disc.",
			Flags: append(OutputFlags("backups", "array"),
				cli.StringFlag{
					Name:  "disc",
					Usage: "the disc you wish to list the backups of",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server you wish to list the backups of",
					Value: new(VirtualMachineNameFlag),
				},
			),
			Action: With(OptionalArgs("server", "disc"), RequiredFlags("server", "disc"), AuthProvider, func(c *Context) (err error) {
				vmName := c.VirtualMachineName("server")
				label := c.String("disc")
				var backups brain.Backups

				if label != "" {
					backups, err = global.Client.GetBackups(vmName, label)
					if err != nil {
						return
					}
				} else {
					err = VirtualMachineProvider("server")(c)
					if err != nil {
						return
					}
					for _, disc := range c.VirtualMachine.Discs {
						discbackups, err := global.Client.GetBackups(vmName, disc.Label)
						if err != nil {
							return err
						}
						backups = append(backups, discbackups...)
					}
				}
				return c.OutputInDesiredForm(backups)
			}),
		}},
	})
}
