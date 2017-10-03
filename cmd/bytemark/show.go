package main

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
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
			Flags: append(app.OutputFlags("account details", "object"),
				cli.GenericFlag{
					Name:  "account",
					Usage: "The account to view",
					Value: new(app.AccountNameFlag),
				}),
			Action: app.With(args.Optional("account"), with.Account("account"), func(c *app.Context) error {
				c.Debug("show account command output")
				c.Debug("acc: %s", c.Account.String())
				return c.OutputInDesiredForm(c.Account)
			}),
		}, {
			Name:        "disc",
			Usage:       "outputs info about a disc",
			UsageText:   "bytemark show disc [--json | --table] [--table-fields help | <fields>] <server> <disc label>",
			Description: `This command displays information about a disc including any backups and backup schedules on the disc`,
			Flags: append(app.OutputFlags("disc details", "object"),
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to display",
					Value: new(app.VirtualMachineNameFlag),
				},
				cli.StringFlag{
					Name:  "disc",
					Usage: "The label or ID of the disc to show",
				},
			),
			Action: app.With(args.Optional("server", "disc"), with.RequiredFlags("server", "disc"), with.Disc("server", "disc"), func(c *app.Context) error {
				return c.OutputInDesiredForm(c.Disc)
			}),
		}, {
			Name:      "group",
			Usage:     "outputs info about a group",
			UsageText: "bytemark show group [--json] [name]",
			Description: `This command displays information about how many servers are in the given group.
If the --json flag is specified, prints a complete overview of the group in JSON format, including all servers.`,
			Flags: append(app.OutputFlags("group details", "object"),
				cli.GenericFlag{
					Name:  "group",
					Usage: "The name of the group to show",
					Value: new(app.GroupNameFlag),
				},
			),
			Action: app.With(args.Optional("group"), with.Group("group"), func(c *app.Context) error {
				return c.OutputInDesiredForm(c.Group)
			}),
		}, {
			Name:        "server",
			Usage:       "displays details about a server",
			UsageText:   "bytemark show server [--json] <name>",
			Description: `Displays a collection of details about the server, including its full hostname, CPU and memory allocation, power status, disc capacities and IP addresses.`,
			Flags: append(app.OutputFlags("server details", "object"),
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to display",
					Value: new(app.VirtualMachineNameFlag),
				},
			),
			Action: app.With(args.Optional("server"), with.RequiredFlags("server"), with.VirtualMachine("server"), func(c *app.Context) error {
				return c.OutputInDesiredForm(c.VirtualMachine)
			}),
		}, {
			Name:        "user",
			Usage:       "displays info about a user",
			UsageText:   "bytemark show user <name>",
			Description: `Currently the only details are what SSH keys are authorised for this user`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "The user to show the details of",
				},
			},
			Action: app.With(args.Optional("user"), with.RequiredFlags("user"), with.User("user"), func(c *app.Context) error {
				log.Outputf("User %s:\n\nAuthorized keys:\n", c.User.Username)
				for _, k := range c.User.AuthorizedKeys {
					log.Output(k)
				}
				return nil
			}),
		}, {
			Name:      "privileges",
			Usage:     "shows privileges for a given account, group, server, or user",
			UsageText: "bytemark show privileges",
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
					Value: new(app.GroupNameFlag),
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "The server to show the privileges of",
					Value: new(app.VirtualMachineNameFlag),
				},
			),
			Action: app.With(with.Auth, func(c *app.Context) (err error) {
				account := c.String("account")
				group := c.GroupName("group")
				server := c.VirtualMachineName("server")

				privs := make(brain.Privileges, 0)
				if account != "" {
					newPrivs, err := findPrivilegesForAccount(c, account, c.Bool("recursive"))
					if err != nil {
						return err
					}
					privs = append(privs, newPrivs...)
				}

				if group.Group != "" {
					newPrivs, err := findPrivilegesForGroup(c, group, c.Bool("recursive"))
					if err != nil {
						return err
					}
					privs = append(privs, newPrivs...)
				}

				if server.VirtualMachine != "" {
					newPrivs, err := c.Client().GetPrivilegesForVirtualMachine(server)
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

	adminCommands = append(adminCommands, cli.Command{
		Name:   "show",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "vlans",
				Usage:     "shows available VLANs",
				UsageText: "bytemark --admin show vlans [--json]",
				Flags:     app.OutputFlags("VLANs", "array"),
				Action: app.With(with.Auth, func(c *app.Context) error {
					vlans, err := c.Client().GetVLANs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vlans, output.Table)
				}),
			},
			{
				Name:        "disc by id",
				Usage:       "displays details about a disc",
				UsageText:   "bytemark show disc by id [--json] <id>",
				Description: `Displays a collection of details about the disc.`,
				Flags: append(app.OutputFlags("disc details", "object"),
					cli.IntFlag{
						Name:  "disc",
						Usage: "the disc to display",
					},
				),
				Action: app.With(with.Auth, args.Optional("disc"), with.RequiredFlags("disc"), func(c *app.Context) error {
					disc, err := c.Client().GetDiscByID(c.Int("disc"))
					if err != nil {
						return err
					}

					return c.OutputInDesiredForm(disc)
				}),
			},
			{
				Name:      "vlan",
				Usage:     "shows the details of a VLAN",
				UsageText: "bytemark --admin show vlan [--json] <num>",
				Flags: append(app.OutputFlags("VLAN", "object"),
					cli.IntFlag{
						Name:  "num",
						Usage: "the num of the VLAN to display",
					},
				),
				Action: app.With(args.Optional("num"), with.RequiredFlags("num"), with.Auth, func(c *app.Context) error {
					vlan, err := c.Client().GetVLAN(c.Int("num"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vlan, output.Table)
				}),
			},
			{
				Name:      "ip ranges",
				Usage:     "shows all IP ranges",
				UsageText: "bytemark --admin show ip ranges [--json]",
				Flags:     app.OutputFlags("ip ranges", "array"),
				Action: app.With(with.Auth, func(c *app.Context) error {
					ipRanges, err := c.Client().GetIPRanges()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(ipRanges, output.Table)
				}),
			},
			{
				Name:      "ip range",
				Usage:     "shows the details of an IP range",
				UsageText: "bytemark --admin show ip range [--json] <ip-range>",
				Flags: append(app.OutputFlags("ip range details", "object"),
					cli.StringFlag{
						Name:  "ip-range",
						Usage: "the ID or CIDR representation of the IP range to display",
					},
				),
				Action: app.With(args.Optional("ip-range"), with.RequiredFlags("ip-range"), with.Auth, func(c *app.Context) error {
					ipRange, err := c.Client().GetIPRange(c.String("ip-range"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(ipRange, output.Table)
				}),
			},
			{
				Name:      "heads",
				Usage:     "shows the details of all heads",
				UsageText: "bytemark --admin show heads [--json]",
				Flags:     app.OutputFlags("heads", "array"),
				Action: app.With(with.Auth, func(c *app.Context) error {
					heads, err := c.Client().GetHeads()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(heads, output.Table)
				}),
			},
			{
				Name:      "head",
				Usage:     "shows the details of the specified head",
				UsageText: "bytemark --admin show head <head> [--json]",
				Flags: append(app.OutputFlags("head details", "object"),
					cli.StringFlag{
						Name:  "head",
						Usage: "the ID of the head to display",
					},
				),
				Action: app.With(args.Optional("head"), with.RequiredFlags("head"), with.Auth, func(c *app.Context) error {
					head, err := c.Client().GetHead(c.String("head"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(head, output.Table)
				}),
			},
			{
				Name:      "tails",
				Usage:     "shows the details of all tails",
				UsageText: "bytemark --admin show tails [--json]",
				Flags:     app.OutputFlags("tails", "array"),
				Action: app.With(with.Auth, func(c *app.Context) error {
					tails, err := c.Client().GetTails()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(tails, output.Table)
				}),
			},
			{
				Name:      "tail",
				Usage:     "shows the details of the specified tail",
				UsageText: "bytemark --admin show tail <tail> [--json]",
				Flags: append(app.OutputFlags("tail details", "object"),
					cli.StringFlag{
						Name:  "tail",
						Usage: "the ID of the tail to display",
					},
				),
				Action: app.With(args.Optional("tail"), with.RequiredFlags("tail"), with.Auth, func(c *app.Context) error {
					tail, err := c.Client().GetTail(c.String("tail"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(tail, output.Table)
				}),
			},
			{
				Name:      "storage pools",
				Usage:     "shows the details of all storage pools",
				UsageText: "bytemark --admin show storage pools [--json]",
				Flags:     app.OutputFlags("storage pools", "array"),
				Action: app.With(with.Auth, func(c *app.Context) error {
					storagePools, err := c.Client().GetStoragePools()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(storagePools, output.Table)
				}),
			},
			{
				Name:      "storage pool",
				Usage:     "shows the details of the specified storage pool",
				UsageText: "bytemark --admin show storage pool [--json] <storage-pool>",
				Flags: append(app.OutputFlags("storage pool", "object"),
					cli.StringFlag{
						Name:  "storage-pool",
						Usage: "The ID or label of the storage pool to display",
					},
				),
				Action: app.With(args.Optional("storage-pool"), with.RequiredFlags("storage-pool"), with.Auth, func(c *app.Context) error {
					storagePool, err := c.Client().GetStoragePool(c.String("storage-pool"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(storagePool, output.Table)
				}),
			},
			{
				Name:      "migrating discs",
				Usage:     "shows a list of migrating discs",
				UsageText: "bytemark --admin show migrating_discs [--json]",
				Flags:     app.OutputFlags("migrating discs", "array"),
				Action: app.With(with.Auth, func(ctx *app.Context) error {
					discs, err := ctx.Client().GetMigratingDiscs()
					if err != nil {
						return err
					}
					// this is super horrid :|
					if ctx.String("table-fields") == "" {
						err := ctx.Context.Set("table-fields", "ID, StoragePool, NewStoragePool, StorageGrade, NewStorageGrade, Size, MigrationProgress, MigrationEta, MigrationSpeed")
						if err != nil {
							return err
						}
					}
					fmt.Fprintln(ctx.App().Writer, "Storage sizes are in MB, speeds in MB/s, and times in seconds.")
					return ctx.OutputInDesiredForm(discs, output.Table)

				}),
			},
			{
				Name:      "migrating vms",
				Usage:     "shows a list of migrating servers",
				UsageText: "bytemark --admin show migrating_vms [--json]",
				Flags:     app.OutputFlags("migrating servers", "array"),
				Action: app.With(with.Auth, func(c *app.Context) error {
					vms, err := c.Client().GetMigratingVMs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vms, output.Table)
				}),
			},
			{
				Name:      "stopped eligible vms",
				Usage:     "shows a list of stopped VMs that should be running",
				UsageText: "bytemark --admin show stopped_eligible_vms [--json]",
				Flags:     app.OutputFlags("servers", "array"),
				Action: app.With(with.Auth, func(c *app.Context) error {
					vms, err := c.Client().GetStoppedEligibleVMs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vms, output.Table)
				}),
			},
			{
				Name:      "recent vms",
				Usage:     "shows a list of stopped VMs that should be running",
				UsageText: "bytemark --admin show recent_vms [--json | --table] [--table-fields <fields> | --table-fields help]",
				Flags:     app.OutputFlags("servers", "array"),
				Action: app.With(with.Auth, func(c *app.Context) error {
					vms, err := c.Client().GetRecentVMs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vms, output.Table)
				}),
			},
		},
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
			Account: account,
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
