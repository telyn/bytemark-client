package main

import (
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
			Flags: append(OutputFlags("account details", "object"),
				cli.GenericFlag{
					Name:  "account",
					Usage: "The account to view",
					Value: new(AccountNameFlag),
				}),
			Action: With(OptionalArgs("account"), AccountProvider("account"), func(c *Context) error {
				return c.OutputInDesiredForm(c.Account)
			}),
		}, {
			Name:        "disc",
			Usage:       "outputs info about a disc",
			UsageText:   "bytemark show disc [--json | --table] [--table-fields help | <fields>] <server> <disc label>",
			Description: `This command displays information about a disc including any backups and backup schedules on the disc`,
			Flags: append(OutputFlags("disc details", "object"),
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to display",
					Value: new(VirtualMachineNameFlag),
				},
				cli.StringFlag{
					Name:  "disc",
					Usage: "The label or ID of the disc to show",
				},
			),
			Action: With(OptionalArgs("server", "disc"), RequiredFlags("server", "disc"), DiscProvider("server", "disc"), func(c *Context) error {
				return c.OutputInDesiredForm(c.Disc)
			}),
		}, {
			Name:      "group",
			Usage:     "outputs info about a group",
			UsageText: "bytemark show group [--json] [name]",
			Description: `This command displays information about how many servers are in the given group.
If the --json flag is specified, prints a complete overview of the group in JSON format, including all servers.`,
			Flags: append(OutputFlags("group details", "object"),
				cli.GenericFlag{
					Name:  "group",
					Usage: "The name of the group to show",
					Value: new(GroupNameFlag),
				},
			),
			Action: With(OptionalArgs("group"), GroupProvider("group"), func(c *Context) error {
				return c.OutputInDesiredForm(c.Group)
			}),
		}, {
			Name:        "server",
			Usage:       "displays details about a server",
			UsageText:   "bytemark show server [--json] <name>",
			Description: `Displays a collection of details about the server, including its full hostname, CPU and memory allocation, power status, disc capacities and IP addresses.`,
			Flags: append(OutputFlags("server details", "object"),
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to display",
					Value: new(VirtualMachineNameFlag),
				},
			),
			Action: With(OptionalArgs("server"), RequiredFlags("server"), VirtualMachineProvider("server"), func(c *Context) error {
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
			Action: With(OptionalArgs("user"), RequiredFlags("user"), UserProvider("user"), func(c *Context) error {
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
			Flags: append(OutputFlags("privileges", "array"),
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
					Value: new(GroupNameFlag),
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "The server to show the privileges of",
					Value: new(VirtualMachineNameFlag),
				},
			),
			Action: With(AuthProvider, func(c *Context) (err error) {
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
				Flags:     OutputFlags("VLANs", "array"),
				Action: With(AuthProvider, func(c *Context) error {
					vlans, err := c.Client().GetVLANs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vlans, output.Table)
				}),
			},
			{
				Name:        "disc_by_id",
				Usage:       "displays details about a disc",
				UsageText:   "bytemark show disc_by_id [--json] <id>",
				Description: `Displays a collection of details about the disc.`,
				Flags: append(OutputFlags("disc details", "object"),
					cli.IntFlag{
						Name:  "disc",
						Usage: "the disc to display",
					},
				),
				Action: With(AuthProvider, OptionalArgs("disc"), RequiredFlags("disc"), func(c *Context) error {
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
				Flags: append(OutputFlags("VLAN", "object"),
					cli.IntFlag{
						Name:  "num",
						Usage: "the num of the VLAN to display",
					},
				),
				Action: With(OptionalArgs("num"), RequiredFlags("num"), AuthProvider, func(c *Context) error {
					vlan, err := c.Client().GetVLAN(c.Int("num"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vlan, output.Table)
				}),
			},
			{
				Name:      "ip_ranges",
				Usage:     "shows all IP ranges",
				UsageText: "bytemark --admin show ip_ranges [--json]",
				Flags:     OutputFlags("ip ranges", "array"),
				Action: With(AuthProvider, func(c *Context) error {
					ipRanges, err := c.Client().GetIPRanges()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(ipRanges, output.Table)
				}),
			},
			{
				Name:      "ip_range",
				Usage:     "shows the details of an IP range",
				UsageText: "bytemark --admin show ip_range [--json] <ip_range>",
				Flags: append(OutputFlags("ip range details", "object"),
					cli.StringFlag{
						Name:  "ip_range",
						Usage: "the ID or CIDR representation of the IP range to display",
					},
				),
				Action: With(OptionalArgs("ip_range"), RequiredFlags("ip_range"), AuthProvider, func(c *Context) error {
					ipRange, err := c.Client().GetIPRange(c.String("ip_range"))
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
				Flags:     OutputFlags("heads", "array"),
				Action: With(AuthProvider, func(c *Context) error {
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
				Flags: append(OutputFlags("head details", "object"),
					cli.StringFlag{
						Name:  "head",
						Usage: "the ID of the head to display",
					},
				),
				Action: With(OptionalArgs("head"), RequiredFlags("head"), AuthProvider, func(c *Context) error {
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
				Flags:     OutputFlags("tails", "array"),
				Action: With(AuthProvider, func(c *Context) error {
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
				Flags: append(OutputFlags("tail details", "object"),
					cli.StringFlag{
						Name:  "tail",
						Usage: "the ID of the tail to display",
					},
				),
				Action: With(OptionalArgs("tail"), RequiredFlags("tail"), AuthProvider, func(c *Context) error {
					tail, err := c.Client().GetTail(c.String("tail"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(tail, output.Table)
				}),
			},
			{
				Name:      "storage_pools",
				Usage:     "shows the details of all storage pools",
				UsageText: "bytemark --admin show storage_pools [--json]",
				Flags:     OutputFlags("storage pools", "array"),
				Action: With(AuthProvider, func(c *Context) error {
					storagePools, err := c.Client().GetStoragePools()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(storagePools, output.Table)
				}),
			},
			{
				Name:      "storage_pool",
				Usage:     "shows the details of the specified storage pool",
				UsageText: "bytemark --admin show storage_pools [--json] <storage_pool>",
				Flags: append(OutputFlags("storage pool", "object"),
					cli.StringFlag{
						Name:  "storage_pool",
						Usage: "The ID or label of the storage pool to display",
					},
				),
				Action: With(OptionalArgs("storage_pool"), RequiredFlags("storage_pool"), AuthProvider, func(c *Context) error {
					storagePool, err := c.Client().GetStoragePool(c.String("storage_pool"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(storagePool, output.Table)
				}),
			},
			{
				Name:      "migrating_discs",
				Usage:     "shows a list of migrating discs",
				UsageText: "bytemark --admin show migrating_discs [--json]",
				Flags:     OutputFlags("migrating discs", "array"),
				Action: With(AuthProvider, func(c *Context) error {
					discs, err := c.Client().GetMigratingDiscs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(discs, output.Table)
				}),
			},
			{
				Name:      "migrating_vms",
				Usage:     "shows a list of migrating servers",
				UsageText: "bytemark --admin show migrating_vms [--json]",
				Flags:     OutputFlags("migrating servers", "array"),
				Action: With(AuthProvider, func(c *Context) error {
					vms, err := c.Client().GetMigratingVMs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vms, output.Table)
				}),
			},
			{
				Name:      "stopped_eligible_vms",
				Usage:     "shows a list of stopped VMs that should be running",
				UsageText: "bytemark --admin show stopped_eligible_vms [--json]",
				Flags:     OutputFlags("servers", "array"),
				Action: With(AuthProvider, func(c *Context) error {
					vms, err := c.Client().GetStoppedEligibleVMs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vms, output.Table)
				}),
			},
			{
				Name:      "recent_vms",
				Usage:     "shows a list of stopped VMs that should be running",
				UsageText: "bytemark --admin show recent_vms [--json | --table] [--table-fields <fields> | --table-fields help]",
				Flags:     OutputFlags("servers", "array"),
				Action: With(AuthProvider, func(c *Context) error {
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

func findPrivilegesForAccount(c *Context, account string, recurse bool) (privs brain.Privileges, err error) {
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

func findPrivilegesForGroup(c *Context, name lib.GroupName, recurse bool) (privs brain.Privileges, err error) {
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
