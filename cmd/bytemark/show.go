package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
)

func init() {

	commands = append(commands, cli.Command{
		Name:        "show",
		Action:      cli.ShowSubcommandHelp,
		Usage:       `Displays information about the given server, group, or account.`,
		UsageText:   "bytemark show account|server|group|user [flags] <name>",
		Description: `Displays information about the given server, group, or account.`,
		Subcommands: []cli.Command{{
			Name:      "account",
			Usage:     `Displays information about the given account`,
			UsageText: "bytemark show account [--json] <name>",
			Description: `This command displays information about the given account, including contact details and how many servers are in it across its groups.
			
If the --json flag is specified, prints a complete overview of the account in JSON format, including all groups and their servers.`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output account details as a JSON object",
				},
			},
			Action: With(AccountProvider, func(c *Context) error {
				return c.IfNotMarshalJSON(c.Account, func() error {
					log.Output(util.FormatAccount(c.Account))

					for _, g := range c.Account.Groups {
						log.Outputf("Group %s\r\n", g.Name)
						for _, v := range util.FormatVirtualMachines(g.VirtualMachines) {
							log.Output(v)
						}
					}
					return nil
				})
			}),
		}, {
			Name:      "group",
			Usage:     "Outputs info about a group",
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

					for _, v := range util.FormatVirtualMachines(c.Group.VirtualMachines) {
						log.Output(v)
					}

					return nil
				})
			}),
		}, {
			Name:        "server",
			Usage:       "Displays details about a server",
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
					log.Log(util.FormatVirtualMachine(c.VirtualMachine))
					return nil
				})
			}),
		}, {
			Name:        "user",
			Usage:       "Displays info about a user.",
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
