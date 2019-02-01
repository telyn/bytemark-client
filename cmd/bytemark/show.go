package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {

	commands = append(commands, cli.Command{
		Name:        "show",
		Action:      cli.ShowSubcommandHelp,
		Usage:       `displays information about the given server, group, or account`,
		UsageText:   "show account|server|group|user [flags] <name>",
		Description: `displays information about the given server, group, or account`,
		Subcommands: []cli.Command{{
			Name:      "account",
			Usage:     `displays information about the given account`,
			UsageText: "show account [--json] [name]",
			Description: `This command displays information about the given account, including contact details and how many servers are in it across its groups.
If no account is specified, it uses your default account.
			
If the --json flag is specified, prints a complete overview of the account in JSON format, including all groups and their servers.`,
			Flags: append(app.OutputFlags("account details", "object"),
				cli.GenericFlag{
					Name:  "account",
					Usage: "The account to view",
					Value: new(flags.AccountNameFlag),
				}),
			Action: app.Action(args.Optional("account"), with.Account("account"), func(c *app.Context) error {
				c.Debug("show account command output")
				c.Debug("acc: %s", c.Account.String())
				return c.OutputInDesiredForm(c.Account)
			}),
		}, {
			Name:        "disc",
			Usage:       "outputs info about a disc",
			UsageText:   "show disc [--json | --table] [--table-fields help | <fields>] [server [disc label]]",
			Description: `This command displays information about a disc including any backups and backup schedules on the disc`,
			Flags: append(app.OutputFlags("disc details", "object"),
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to display",
					Value: new(flags.VirtualMachineNameFlag),
				},
				cli.StringFlag{
					Name:  "disc",
					Usage: "The label or ID of the disc to show",
				},
			),
			Action: app.Action(args.Optional("server", "disc"), with.RequiredFlags("disc"), with.Disc("server", "disc"), func(c *app.Context) error {
				return c.OutputInDesiredForm(c.Disc)
			}),
		}, {
			Name:      "group",
			Usage:     "outputs info about a group",
			UsageText: "show group [--json] [name]",
			Description: `This command displays information about how many servers are in the given group.
If the --json flag is specified, prints a complete overview of the group in JSON format, including all servers.`,
			Flags: append(app.OutputFlags("group details", "object"),
				cli.GenericFlag{
					Name:  "group",
					Usage: "The name of the group to show",
					Value: new(flags.GroupNameFlag),
				},
			),
			Action: app.Action(args.Optional("group"), with.Group("group"), func(c *app.Context) error {
				return c.OutputInDesiredForm(c.Group)
			}),
		}, {
			Name:        "server",
			Usage:       "displays details about a server",
			UsageText:   "show server [--json] <name>",
			Description: `Displays a collection of details about the server, including its full hostname, CPU and memory allocation, power status, disc capacities and IP addresses.`,
			Flags: append(app.OutputFlags("server details", "object"),
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to display",
					Value: new(flags.VirtualMachineNameFlag),
				},
			),
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.VirtualMachine("server"), func(c *app.Context) error {
				return c.OutputInDesiredForm(c.VirtualMachine)
			}),
		}, {
			Name:        "user",
			Usage:       "displays info about a user",
			UsageText:   "show user <name>",
			Description: `Currently the only details are what SSH keys are authorised for this user`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "The user to show the details of",
				},
			},
			Action: app.Action(args.Optional("user"), with.RequiredFlags("user"), with.User("user"), func(c *app.Context) error {
				log.Outputf("User %s:\n\nAuthorized keys:\n", c.User.Username)
				for _, k := range c.User.AuthorizedKeys {
					log.Output(k)
				}
				return nil
			}),
		}, {
			Name:      "privileges",
			Usage:     "shows privileges for a given account, group, server, or user",
			UsageText: "show privileges",
			Description: `Displays a list of all the privileges for a given account, group, server or user. If none are specified, shows the privileges for your user.

Setting --recursive will cause a lot of extra requests to be made and may take a long time to run.

Privileges will be output in no particular order.`,
			Flags: append(app.OutputFlags("privileges", "array"),
				cli.BoolFlag{
					Name:  "recursive",
					Usage: "for account & group, will also find all privileges for all groups in the account and virtual machines in the group",
				},
				cli.StringFlag{
					Name:  "user",
					Usage: "The user to show the privileges of",
				},
				cli.StringFlag{
					Name:  "account",
					Usage: "The account to show the privileges of",
				},
				cli.GenericFlag{
					Name:  "group",
					Usage: "The group to show the privileges of",
					Value: new(flags.GroupNameFlag),
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "The server to show the privileges of",
					Value: new(flags.VirtualMachineNameFlag),
				},
			),
			Action: app.Action(with.Auth, func(c *app.Context) (err error) {
				account := c.String("account")
				group := flags.GroupName(c, "group")
				server := flags.VirtualMachineName(c, "server")

				privs := make(brain.Privileges, 0)
				var newPrivs brain.Privileges
				if account != "" {
					newPrivs, err = findPrivilegesForAccount(c, account, c.Bool("recursive"))
					if err != nil {
						return err
					}
					privs = append(privs, newPrivs...)
				}

				if group.Group != "" {
					newPrivs, err = findPrivilegesForGroup(c, group, c.Bool("recursive"))
					if err != nil {
						return err
					}
					privs = append(privs, newPrivs...)
				}

				if server.VirtualMachine != "" {
					newPrivs, err = c.Client().GetPrivilegesForVirtualMachine(server)
					if err != nil {
						return err
					}
					privs = append(privs, newPrivs...)
				}
				if c.String("user") != "" || (server.VirtualMachine == "" && group.Group == "" && account == "") {

					privs, err = c.Client().GetPrivileges(c.String("user"))
					if err != nil {
						return
					}
				}

				return c.OutputInDesiredForm(privs, output.List)
			}),
		}},
	})

}

func findPrivilegesForAccount(c *app.Context, account string, recurse bool) (privs brain.Privileges, err error) {
	privs, err = c.Client().GetPrivilegesForAccount(account)
	if !recurse || err != nil {
		return
	}
	acc, err := c.Client().GetAccount(account)
	if err != nil {
		return
	}

	for _, group := range acc.Groups {
		newPrivs, err := findPrivilegesForGroup(c, lib.GroupName{
			Group:   group.Name,
			Account: pathers.AccountName(account),
		}, recurse) // recurse is always true at this point but maybe I'd like to make two flags? recurse-account and recurse-group?
		if err != nil {
			return privs, err
		}
		privs = append(privs, newPrivs...)
	}
	return
}

func findPrivilegesForGroup(c *app.Context, name lib.GroupName, recurse bool) (privs brain.Privileges, err error) {
	privs, err = c.Client().GetPrivilegesForGroup(name)
	if !recurse || err != nil {
		return
	}
	group, err := c.Client().GetGroup(name)
	if err != nil {
		return
	}
	for _, vm := range group.VirtualMachines {
		vmName := lib.VirtualMachineName{
			VirtualMachine: vm.Name,
			Group:          name.Group,
			Account:        name.Account,
		}
		newPrivs, err := c.Client().GetPrivilegesForVirtualMachine(vmName)
		if err != nil {
			return privs, err
		}
		privs = append(privs, newPrivs...)
	}
	return
}
