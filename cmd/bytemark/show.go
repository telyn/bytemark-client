package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"os"
)

func init() {

	commands = append(commands, cli.Command{
		Name:        "show",
		Action:      cli.ShowSubcommandHelp,
		Usage:       `displays information about the given server, group, or account`,
		UsageText:   "bytemark show account|server|group|user [flags] <name>",
		Description: `displays information about the given server, group, or account`,
		Subcommands: []cli.Command{{
			Name:      "account",
			Usage:     `displays information about the given account`,
			UsageText: "bytemark show account [--json] [name]",
			Description: `This command displays information about the given account, including contact details and how many servers are in it across its groups.
If no account is specified, it uses your default account.
			
If the --json flag is specified, prints a complete overview of the account in JSON format, including all groups and their servers.`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output account details as a JSON object",
				},
			},
			Action: With(AccountProvider(false), func(c *Context) error {
				return c.IfNotMarshalJSON(c.Account, func() error {
					def, err := global.Client.GetDefaultAccount()
					if err != nil {
						return err
					}
					err = lib.FormatAccount(os.Stderr, c.Account, def, "account_overview")
					if err != nil {
						return err
					}
					log.Output()
					log.Output()

					for _, g := range c.Account.Groups {
						for _, vm := range g.VirtualMachines {
							err := lib.FormatVirtualMachine(os.Stderr, vm, lib.TwoLine)
							log.Output()
							log.Output()
							if err != nil {
								return err
							}
						}
					}
					return nil
				})
			}),
		}, {
			Name:      "group",
			Usage:     "outputs info about a group",
			UsageText: "bytemark show group [--json] <name>",
			Description: `This command displays information about how many servers are in the given group.
If the --json flag is specified, prints a complete overview of the group in JSON format, including all servers.`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output group details as a JSON object",
				},
			},
			Action: With(GroupProvider, func(c *Context) error {
				return c.IfNotMarshalJSON(c.Group, func() error {
					s := ""
					if len(c.Group.VirtualMachines) != 1 {
						s = "s"
					}
					log.Outputf("%s - Group containing %d cloud server%s\r\n", c.Group.Name, len(c.Group.VirtualMachines), s)

					log.Output()
					for _, vm := range c.Group.VirtualMachines {

						err := lib.FormatVirtualMachine(os.Stderr, vm, lib.TwoLine)
						log.Output()
						log.Output()
						if err != nil {
							return err
						}
					}

					return nil
				})
			}),
		}, {
			Name:        "server",
			Usage:       "displays details about a server",
			UsageText:   "bytemark show server [--json] <name>",
			Description: `Displays a collection of details about the server, including its full hostname, CPU and memory allocation, power status, disc capacities and IP addresses.`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output server details as a JSON object.",
				},
			},
			Action: With(VirtualMachineProvider, func(c *Context) error {
				return c.IfNotMarshalJSON(c.VirtualMachine, func() error {
					return lib.FormatVirtualMachine(os.Stderr, c.VirtualMachine, lib.All)
				})
			}),
		}, {
			Name:        "user",
			Usage:       "displays info about a user",
			UsageText:   "bytemark show user <name>",
			Description: `Currently the only details are what SSH keys are authorised for this user`,
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
