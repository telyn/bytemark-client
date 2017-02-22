package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
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
			Action: With(OptionalArgs("account"), AccountProvider("account"), func(c *Context) error {
				return c.IfNotMarshalJSON(c.Account, func() error {
					def, err := global.Client.GetDefaultAccount()
					if err != nil {
						return err
					}
					if def.BrainID == c.Account.BrainID {
						c.Account.IsDefaultAccount = true
					}
					err = c.Account.PrettyPrint(os.Stderr, prettyprint.Full)
					if err != nil {
						return err
					}
					log.Output()
					log.Output()

					for _, g := range c.Account.Groups {
						for _, vm := range g.VirtualMachines {
							err := vm.PrettyPrint(os.Stderr, prettyprint.Medium)
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output group details as a JSON object",
				},
				cli.GenericFlag{
					Name:  "group",
					Usage: "The name of the group to show",
					Value: new(GroupNameFlag),
				},
			},
			Action: With(OptionalArgs("group"), GroupProvider("group"), func(c *Context) error {
				return c.IfNotMarshalJSON(c.Group, func() error {
					s := ""
					if len(c.Group.VirtualMachines) != 1 {
						s = "s"
					}
					log.Outputf("%s - Group containing %d cloud server%s\r\n", c.Group.Name, len(c.Group.VirtualMachines), s)

					log.Output()
					for _, vm := range c.Group.VirtualMachines {

						err := vm.PrettyPrint(os.Stderr, prettyprint.Medium)
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
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to display",
					Value: new(VirtualMachineNameFlag),
				},
			},
			Action: With(OptionalArgs("server"), RequiredFlags("server"), VirtualMachineProvider("server"), func(c *Context) error {
				return c.IfNotMarshalJSON(c.VirtualMachine, func() error {
					return c.VirtualMachine.PrettyPrint(os.Stderr, prettyprint.Full)
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "Output privileges as a JSON array.",
				},
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
			},
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

				return c.IfNotMarshalJSON(privs, func() error {
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
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the VLANs as a JSON array.",
					},
				},
				Action: With(AuthProvider, func(c *Context) error {
					vlans, err := global.Client.GetVLANs()
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(vlans, func() error {
						for _, vlan := range vlans {
							log.Outputf("%s\r\n", vlan.String())
						}
						return nil
					})
				}),
			},
			{
				Name:      "ip_ranges",
				Usage:     "shows all IP ranges",
				UsageText: "bytemark --admin show ip_ranges [--json]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output IP ranges as a JSON array.",
					},
				}, Action: With(AuthProvider, func(c *Context) error {
					ipRanges, err := global.Client.GetIPRanges()
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(ipRanges, func() error {
						for _, ipRange := range ipRanges {
							log.Outputf("%s\r\n", ipRange.String())
						}
						return nil
					})
				}),
			},
			{
				Name:      "ip_range",
				Usage:     "shows the details of an IP range",
				UsageText: "bytemark --admin show ip_range [--json] <ip_range>",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the IP range as a JSON object.",
					},
					cli.IntFlag{
						Name:  "ip_range",
						Usage: "the ID of the IP range to display",
					},
				}, Action: With(OptionalArgs("ip_range"), RequiredFlags("ip_range"), AuthProvider, func(c *Context) error {
					ipRange, err := global.Client.GetIPRange(c.Int("ip_range"))
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(ipRange, func() error {
						log.Outputf("%s\r\n", ipRange.String())
						return nil
					})
				}),
			},
			{
				Name:      "heads",
				Usage:     "shows the details of all heads",
				UsageText: "bytemark --admin show heads [--json]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the heads as a JSON array.",
					},
				}, Action: With(AuthProvider, func(c *Context) error {
					heads, err := global.Client.GetHeads()
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(heads, func() error {
						for _, head := range heads {
							if err := head.PrettyPrint(os.Stderr, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					})
				}),
			},
			{
				Name:      "head",
				Usage:     "shows the details of the specified head",
				UsageText: "bytemark --admin show head <head> [--json]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the head as a JSON object.",
					},
					cli.StringFlag{
						Name:  "head",
						Usage: "the ID of the head to display",
					},
				}, Action: With(OptionalArgs("head"), RequiredFlags("head"), AuthProvider, func(c *Context) error {
					head, err := global.Client.GetHead(c.String("head"))
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(head, func() error {
						return head.PrettyPrint(os.Stderr, prettyprint.Full)
					})
				}),
			},
			{
				Name:      "tails",
				Usage:     "shows the details of all tails",
				UsageText: "bytemark --admin show tails [--json]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the tails as a JSON array.",
					},
				}, Action: With(AuthProvider, func(c *Context) error {
					tails, err := global.Client.GetTails()
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(tails, func() error {
						for _, tail := range tails {
							if err := tail.PrettyPrint(os.Stderr, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					})
				}),
			},
			{
				Name:      "tail",
				Usage:     "shows the details of the specified tail",
				UsageText: "bytemark --admin show tail <tail> [--json]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the tail as a JSON object.",
					},
					cli.StringFlag{
						Name:  "tail",
						Usage: "the ID of the tail to display",
					},
				}, Action: With(OptionalArgs("tail"), RequiredFlags("tail"), AuthProvider, func(c *Context) error {
					tail, err := global.Client.GetTail(c.String("tail"))
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(tail, func() error {
						return tail.PrettyPrint(os.Stderr, prettyprint.Full)
					})
				}),
			},
			{
				Name:      "storage_pools",
				Usage:     "shows the details of all storage pools",
				UsageText: "bytemark --admin show storage_pools [--json]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the storage pools as a JSON array.",
					},
				}, Action: With(AuthProvider, func(c *Context) error {
					storagePools, err := global.Client.GetStoragePools()
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(storagePools, func() error {
						for _, storagePool := range storagePools {
							if err := storagePool.PrettyPrint(os.Stderr, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					})
				}),
			},
			{
				Name:      "storage_pool",
				Usage:     "shows the details of the specified storage pool",
				UsageText: "bytemark --admin show storage_pools [--json] <storage_pool>",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the storage pools as a JSON array.",
					},
					cli.StringFlag{
						Name:  "storage_pool",
						Usage: "The ID or label of the storage pool to display",
					},
				}, Action: With(OptionalArgs("storage_pool"), RequiredFlags("storage_pool"), AuthProvider, func(c *Context) error {
					storagePool, err := global.Client.GetStoragePool(c.String("storage_pool"))
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(storagePool, func() error {
						return storagePool.PrettyPrint(os.Stderr, prettyprint.SingleLine)
					})
				}),
			},
			{
				Name:      "migrating_vms",
				Usage:     "shows a list of migrating VMs",
				UsageText: "bytemark --admin show migrating_vms [--json]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the VMs as a JSON array.",
					},
				}, Action: With(AuthProvider, func(c *Context) error {
					vms, err := global.Client.GetMigratingVMs()
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(vms, func() error {
						for _, vm := range vms {
							if err := vm.PrettyPrint(os.Stderr, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					})
				}),
			},
			{
				Name:      "stopped_eligible_vms",
				Usage:     "shows a list of stopped VMs that should be running",
				UsageText: "bytemark --admin show stopped_eligible_vms [--json]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "json",
						Usage: "Output the VMs as a JSON array.",
					},
				}, Action: With(AuthProvider, func(c *Context) error {
					vms, err := global.Client.GetStoppedEligibleVMs()
					if err != nil {
						return err
					}
					return c.IfNotMarshalJSON(vms, func() error {
						for _, vm := range vms {
							if err := vm.PrettyPrint(os.Stderr, prettyprint.SingleLine); err != nil {
								return err
							}
						}

						return nil
					})
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
