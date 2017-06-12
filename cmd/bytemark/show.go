package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
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
			Flags: append(OutputFlags("account details", "object", DefaultAccountTableFields),
				cli.GenericFlag{
					Name:  "account",
					Usage: "The account to view",
					Value: new(AccountNameFlag),
				}),
			Action: With(OptionalArgs("account"), AccountProvider("account"), func(c *Context) error {
				return c.OutputInDesiredForm(c.Account, func() error {
					err := c.Account.PrettyPrint(global.App.Writer, prettyprint.Full)
					if err != nil {
						return err
					}
					log.Output()
					log.Output()

					for _, g := range c.Account.Groups {
						for _, vm := range g.VirtualMachines {
							err := vm.PrettyPrint(global.App.Writer, prettyprint.Medium)
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
			UsageText: "bytemark show group [--json] [name]",
			Description: `This command displays information about how many servers are in the given group.
If the --json flag is specified, prints a complete overview of the group in JSON format, including all servers.`,
			Flags: append(OutputFlags("group details", "object", DefaultGroupTableFields),
				cli.GenericFlag{
					Name:  "group",
					Usage: "The name of the group to show",
					Value: new(GroupNameFlag),
				},
			),
			Action: With(OptionalArgs("group"), GroupProvider("group"), func(c *Context) error {
				return c.OutputInDesiredForm(c.Group, func() error {
					s := ""
					if len(c.Group.VirtualMachines) != 1 {
						s = "s"
					}
					log.Outputf("%s - Group containing %d cloud server%s\r\n", c.Group.Name, len(c.Group.VirtualMachines), s)

					log.Output()
					for _, vm := range c.Group.VirtualMachines {

						err := vm.PrettyPrint(global.App.Writer, prettyprint.Medium)
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
			Flags: append(OutputFlags("server details", "object", DefaultServerTableFields),
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to display",
					Value: new(VirtualMachineNameFlag),
				},
			),
			Action: With(OptionalArgs("server"), RequiredFlags("server"), VirtualMachineProvider("server"), func(c *Context) error {
				return c.OutputInDesiredForm(c.VirtualMachine, func() error {
					return c.VirtualMachine.PrettyPrint(global.App.Writer, prettyprint.Full)
				})
			}),
		}, {
			Name:        "user",
			Usage:       "displays info about a user",
			UsageText:   "bytemark show user <name>",
			Description: `Currently the only details are what SSH keys are authorised for this user`,
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
			Flags: append(OutputFlags("privileges", "array", DefaultPrivilegeTableFields),
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
					newPrivs, err := findPrivilegesForAccount(account, c.Bool("recursive"))
					if err != nil {
						return err
					}
					privs = append(privs, newPrivs...)
				}

				if group.Group != "" {
					newPrivs, err := findPrivilegesForGroup(group, c.Bool("recursive"))
					if err != nil {
						return err
					}
					privs = append(privs, newPrivs...)
				}

				if server.VirtualMachine != "" {
					newPrivs, err := global.Client.GetPrivilegesForVirtualMachine(server)
					if err != nil {
						return err
					}
					privs = append(privs, newPrivs...)
				}
				if c.String("user") != "" || (server.VirtualMachine == "" && group.Group == "" && account == "") {

					privs, err = global.Client.GetPrivileges(c.String("user"))
					if err != nil {
						return
					}
				}

				return c.OutputInDesiredForm(privs, func() error {
					for _, p := range privs {
						log.Outputf("%s\r\n", p.String())
					}
					return nil
				})
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
				Flags:     OutputFlags("VLANs", "array", DefaultVLANTableFields),
				Action: With(AuthProvider, func(c *Context) error {
					vlans, err := global.Client.GetVLANs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vlans, func() error {
						for _, vlan := range vlans {
							if err := vlan.PrettyPrint(global.App.Writer, prettyprint.SingleLine); err != nil {
								return err
							}
						}
						return nil
					}, "table")
				}),
			},
			{
				Name:      "vlan",
				Usage:     "shows the details of a VLAN",
				UsageText: "bytemark --admin show vlan [--json] <num>",
				Flags: append(OutputFlags("VLAN", "object", DefaultVLANTableFields),
					cli.IntFlag{
						Name:  "num",
						Usage: "the num of the VLAN to display",
					},
				),
				Action: With(OptionalArgs("num"), RequiredFlags("num"), AuthProvider, func(c *Context) error {
					vlan, err := global.Client.GetVLAN(c.Int("num"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vlan, func() error {
						return vlan.PrettyPrint(global.App.Writer, prettyprint.Full)
					}, "table")
				}),
			},
			{
				Name:      "ip_ranges",
				Usage:     "shows all IP ranges",
				UsageText: "bytemark --admin show ip_ranges [--json]",
				Flags:     OutputFlags("ip ranges", "array", DefaultIPRangeTableFields),
				Action: With(AuthProvider, func(c *Context) error {
					ipRanges, err := global.Client.GetIPRanges()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(ipRanges, func() error {
						for _, ipRange := range ipRanges {
							if err := ipRange.PrettyPrint(global.App.Writer, prettyprint.SingleLine); err != nil {
								return err
							}
						}
						return nil
					}, "table")
				}),
			},
			{
				Name:      "ip_range",
				Usage:     "shows the details of an IP range",
				UsageText: "bytemark --admin show ip_range [--json] <ip_range>",
				Flags: append(OutputFlags("ip range details", "object", DefaultIPRangeTableFields),
					cli.StringFlag{
						Name:  "ip_range",
						Usage: "the ID or CIDR representation of the IP range to display",
					},
				),
				Action: With(OptionalArgs("ip_range"), RequiredFlags("ip_range"), AuthProvider, func(c *Context) error {
					ipRange, err := global.Client.GetIPRange(c.String("ip_range"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(ipRange, func() error {
						return ipRange.PrettyPrint(global.App.Writer, prettyprint.Full)
					}, "table")
				}),
			},
			{
				Name:      "heads",
				Usage:     "shows the details of all heads",
				UsageText: "bytemark --admin show heads [--json]",
				Flags:     OutputFlags("heads", "array", DefaultHeadTableFields),
				Action: With(AuthProvider, func(c *Context) error {
					heads, err := global.Client.GetHeads()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(heads, func() error {
						for _, head := range heads {
							if err := head.PrettyPrint(global.App.Writer, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					}, "table")
				}),
			},
			{
				Name:      "head",
				Usage:     "shows the details of the specified head",
				UsageText: "bytemark --admin show head <head> [--json]",
				Flags: append(OutputFlags("head details", "object", DefaultHeadTableFields),
					cli.StringFlag{
						Name:  "head",
						Usage: "the ID of the head to display",
					},
				),
				Action: With(OptionalArgs("head"), RequiredFlags("head"), AuthProvider, func(c *Context) error {
					head, err := global.Client.GetHead(c.String("head"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(head, func() error {
						return head.PrettyPrint(global.App.Writer, prettyprint.Full)
					}, "table")
				}),
			},
			{
				Name:      "tails",
				Usage:     "shows the details of all tails",
				UsageText: "bytemark --admin show tails [--json]",
				Flags:     OutputFlags("tails", "array", DefaultTailTableFields),
				Action: With(AuthProvider, func(c *Context) error {
					tails, err := global.Client.GetTails()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(tails, func() error {
						for _, tail := range tails {
							if err := tail.PrettyPrint(global.App.Writer, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					}, "table")
				}),
			},
			{
				Name:      "tail",
				Usage:     "shows the details of the specified tail",
				UsageText: "bytemark --admin show tail <tail> [--json]",
				Flags: append(OutputFlags("tail details", "object", DefaultTailTableFields),
					cli.StringFlag{
						Name:  "tail",
						Usage: "the ID of the tail to display",
					},
				),
				Action: With(OptionalArgs("tail"), RequiredFlags("tail"), AuthProvider, func(c *Context) error {
					tail, err := global.Client.GetTail(c.String("tail"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(tail, func() error {
						return tail.PrettyPrint(global.App.Writer, prettyprint.Full)
					}, "table")
				}),
			},
			{
				Name:      "storage_pools",
				Usage:     "shows the details of all storage pools",
				UsageText: "bytemark --admin show storage_pools [--json]",
				Flags:     OutputFlags("storage pools", "array", DefaultStoragePoolTableFields),
				Action: With(AuthProvider, func(c *Context) error {
					storagePools, err := global.Client.GetStoragePools()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(storagePools, func() error {
						for _, storagePool := range storagePools {
							if err := storagePool.PrettyPrint(global.App.Writer, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					}, "table")
				}),
			},
			{
				Name:      "storage_pool",
				Usage:     "shows the details of the specified storage pool",
				UsageText: "bytemark --admin show storage_pools [--json] <storage_pool>",
				Flags: append(OutputFlags("storage pool", "object", DefaultStoragePoolTableFields),
					cli.StringFlag{
						Name:  "storage_pool",
						Usage: "The ID or label of the storage pool to display",
					},
				),
				Action: With(OptionalArgs("storage_pool"), RequiredFlags("storage_pool"), AuthProvider, func(c *Context) error {
					storagePool, err := global.Client.GetStoragePool(c.String("storage_pool"))
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(storagePool, func() error {
						return storagePool.PrettyPrint(global.App.Writer, prettyprint.Full)
					}, "table")
				}),
			},
			{
				Name:      "migrating_vms",
				Usage:     "shows a list of migrating servers",
				UsageText: "bytemark --admin show migrating_vms [--json]",
				Flags:     OutputFlags("migrating servers", "array", DefaultServerTableFields),
				Action: With(AuthProvider, func(c *Context) error {
					vms, err := global.Client.GetMigratingVMs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vms, func() error {
						for _, vm := range vms {
							if err := vm.PrettyPrint(global.App.Writer, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					}, "table")
				}),
			},
			{
				Name:      "stopped_eligible_vms",
				Usage:     "shows a list of stopped VMs that should be running",
				UsageText: "bytemark --admin show stopped_eligible_vms [--json]",
				Flags:     OutputFlags("servers", "array", DefaultServerTableFields),
				Action: With(AuthProvider, func(c *Context) error {
					vms, err := global.Client.GetStoppedEligibleVMs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vms, func() error {
						for _, vm := range vms {
							if err := vm.PrettyPrint(global.App.Writer, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					}, "table")
				}),
			},
			{
				Name:      "recent_vms",
				Usage:     "shows a list of stopped VMs that should be running",
				UsageText: "bytemark --admin show recent_vms [--json | --table] [--table-fields <fields> | --table-fields help]",
				Flags:     OutputFlags("servers", "array", DefaultServerTableFields),
				Action: With(AuthProvider, func(c *Context) error {
					vms, err := global.Client.GetRecentVMs()
					if err != nil {
						return err
					}
					return c.OutputInDesiredForm(vms, func() error {
						for _, vm := range vms {
							if err := vm.PrettyPrint(global.App.Writer, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					}, "table")
				}),
			},
		},
	})
}

func findPrivilegesForAccount(account string, recurse bool) (privs brain.Privileges, err error) {
	privs, err = global.Client.GetPrivilegesForAccount(account)
	if !recurse || err != nil {
		return
	}
	acc, err := global.Client.GetAccount(account)
	if err != nil {
		return
	}

	for _, group := range acc.Groups {
		newPrivs, err := findPrivilegesForGroup(lib.GroupName{
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

func findPrivilegesForGroup(name lib.GroupName, recurse bool) (privs brain.Privileges, err error) {
	privs, err = global.Client.GetPrivilegesForGroup(name)
	if !recurse || err != nil {
		return
	}
	group, err := global.Client.GetGroup(&name)
	if err != nil {
		return
	}
	for _, vm := range group.VirtualMachines {
		vmName := lib.VirtualMachineName{
			VirtualMachine: vm.Name,
			Group:          name.Group,
			Account:        name.Account,
		}
		newPrivs, err := global.Client.GetPrivilegesForVirtualMachine(vmName)
		if err != nil {
			return privs, err
		}
		privs = append(privs, newPrivs...)
	}
	return
}
