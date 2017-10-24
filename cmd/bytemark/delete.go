package main

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:   "delete",
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			{
				Name:      "vlan",
				Usage:     "delete a given VLAN",
				UsageText: "bytemark --admin delete vlan <id>",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "id",
						Usage: "the ID of the VLAN to delete",
					},
				},
				Action: app.Action(args.Optional("id"), with.RequiredFlags("id"), with.Auth, func(c *app.Context) error {
					if err := c.Client().DeleteVLAN(c.Int("id")); err != nil {
						return err
					}

					log.Output("VLAN deleted")

					return nil
				}),
			},
		},
	})

	commands = append(commands, cli.Command{
		Name:      "delete",
		Usage:     "delete a given server, disc, group, account or key",
		UsageText: "bytemark delete account|disc|group|key|server",
		Description: `delete a given server, disc, group, account or key

   Only empty groups and accounts can be deleted.

   The restore server command may be used to restore a deleted (but not purged) server to its state prior to deletion.`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "disc",
			Usage:       "delete the given disc",
			UsageText:   "bytemark delete disc <virtual machine name> <disc label>",
			Description: "Deletes the given disc. To find out a disc's label you can use the `bytemark show server` command or `bytemark list discs` command.",
			Flags: []cli.Flag{
				forceFlag,
				cli.StringFlag{
					Name:  "disc",
					Usage: "the disc to delete",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server whose disc you wish to delete",
					Value: new(app.VirtualMachineNameFlag),
				},
			},
			Aliases: []string{"disk"},
			Action: app.Action(args.Optional("server", "disc"), with.RequiredFlags("server", "disc"), with.Auth, func(c *app.Context) (err error) {
				if !c.Bool("force") && !util.PromptYesNo("Are you sure you wish to delete this disc? It is impossible to recover.") {
					return util.UserRequestedExit{}
				}
				vmName := c.VirtualMachineName("server")
				return c.Client().DeleteDisc(vmName, c.String("disc"))
			}),
		}, {
			Name:      "group",
			Usage:     "deletes the given group",
			UsageText: "bytemark delete group [--force] [--recursive] <group name>",
			Description: `Deletes the given group.
If --recursive is specified, all servers in the group will be purged. Otherwise, if there are servers in the group, will return an error.`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "recursive",
					Usage: "If set, all servers in the group will be irrevocably deleted.",
				},
				cli.GenericFlag{
					Name:  "group",
					Usage: "the name of the group to delete",
					Value: new(app.GroupNameFlag),
				},
				forceFlag,
			},
			Action: app.Action(args.Optional("group"), with.RequiredFlags("group"), with.Group("group"), deleteGroup),
		}, {
			Name:        "key",
			Usage:       "deletes the specified key",
			UsageText:   "bytemark delete key [--user <user>] <key>",
			Description: "Keys may be specified as just the comment part or as the whole key. If there are multiple keys with the comment given, an error will be returned",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "Which user to delete the key from. Defaults to the username you log in as.",
				},
				cli.StringFlag{
					Name:  "public-key",
					Usage: "The public key to delete. Can be the comment part or the whole public key",
				},
			},
			Action: app.Action(args.Join("public-key"), with.RequiredFlags("public-key"), with.Auth, func(c *app.Context) (err error) {
				user := c.String("user")
				if user == "" {
					user = c.Config().GetIgnoreErr("user")
				}

				key := c.String("public-key")
				if key == "" {
					return c.Help("You must specify a key to delete.\r\n")
				}

				err = c.Client().DeleteUserAuthorizedKey(user, key)
				if err == nil {
					log.Log("Key deleted successfully")
				}
				return
			}),
		}, {
			Name:        "server",
			Usage:       "delete the given server",
			UsageText:   `bytemark delete server [--purge] <server name>`,
			Description: "Deletes the given server. Deleted servers still exist and can be restored. To ensure a server is fully deleted, use the --purge flag.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "purge",
					Usage: "If set, the server will be irrevocably deleted.",
				},
				forceFlag,
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to delete",
					Value: new(app.VirtualMachineNameFlag),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.VirtualMachine("server"), deleteServer),
		}, {
			Name:        "backup",
			Usage:       "delete the given backup",
			UsageText:   `bytemark delete backup <server name> <disc label> <backup label>`,
			Description: "Deletes the given backup. Backups cannot be recovered after deletion.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "disc",
					Usage: "the disc to delete a backup of",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to delete a backup from",
					Value: new(app.VirtualMachineNameFlag),
				},
				cli.StringFlag{
					Name:  "backup",
					Usage: "the name or ID of the backup to delete",
				},
			},
			Action: app.Action(args.Optional("server", "disc", "backup"), with.RequiredFlags("server", "disc", "backup"), with.Auth, deleteBackup),
		}},
	})
}

func deleteServer(c *app.Context) (err error) {
	purge := c.Bool("purge")
	vm := c.VirtualMachine

	if vm.Deleted && !purge {
		log.Errorf("Server %s has already been deleted.\r\nIf you wish to permanently delete it, add --purge\r\n", vm.Hostname)
		// we don't return an error because we want a 0 exit code - the deletion request has happened, just not now.
		return
	}
	fstr := fmt.Sprintf("Are you certain you wish to delete %s?", vm.Hostname)
	if purge {
		fstr = fmt.Sprintf("Are you certain you wish to permanently delete %s? You will not be able to un-delete it.", vm.Hostname)

	}

	if !c.Bool("force") && !util.PromptYesNo(fstr) {
		err = util.UserRequestedExit{}
		return
	}

	vmName := c.VirtualMachineName("server")
	err = c.Client().DeleteVirtualMachine(vmName, purge)
	if err != nil {
		return
	}
	if purge {
		log.Logf("Server %s purged successfully.\r\n", vm.Hostname)
	} else {
		log.Logf("Server %s deleted successfully.\r\n", vm.Hostname)
	}
	return
}

func countRunning(group *brain.Group) (running int) {
	for _, vm := range group.VirtualMachines {
		if vm.PowerOn {
			running++
		}
	}
	return running
}

func deleteGroup(c *app.Context) (err error) {
	recursive := c.Bool("recursive")
	groupName := c.GroupName("group")
	if len(c.Group.VirtualMachines) > 0 && recursive {
		prompt := fmt.Sprintf("The group '%s' has %d servers in it which will be irrevocably deleted", c.Group.Name, len(c.Group.VirtualMachines))
		running := countRunning(c.Group)
		if running != 0 {
			stopped := len(c.Group.VirtualMachines) - running
			andStopped := ""
			if stopped > 0 {
				andStopped = fmt.Sprintf("and %d stopped ", stopped)
			}
			prompt = fmt.Sprintf("The group '%s' has %d currently-running %sservers in it which will be forcibly stopped and irrevocably deleted", c.Group.Name, running, andStopped)
		}

		if !c.Bool("force") && !util.PromptYesNo(prompt+" - are you sure you wish to delete this group?") {
			return util.UserRequestedExit{}
		}
		err = recursiveDeleteGroup(c, &groupName, c.Group)
		if err != nil {
			return
		}
	} else if !recursive {
		err = &util.WontDeleteNonEmptyGroupError{Group: &groupName}
		return
	}
	err = c.Client().DeleteGroup(groupName)
	if err == nil {
		log.Logf("Group %s deleted successfully.\r\n", groupName.String())
	}
	return
}

func recursiveDeleteGroup(c *app.Context, name *lib.GroupName, group *brain.Group) error {
	log.Log("", "")
	vmn := lib.VirtualMachineName{Group: name.Group, Account: name.Account}
	for _, vm := range group.VirtualMachines {
		vmn.VirtualMachine = vm.Name
		err := c.Client().DeleteVirtualMachine(vmn, true)
		if err != nil {
			return err
		}
		log.Logf("%s\r\n", vm.Name)

	}
	log.Log()
	return nil
}

/*log.Log("usage: bytemark delete account <account>")
	log.Log("       bytemark delete disc <server> <label>")
	log.Log("       bytemark delete group [--recursive] <group>")
	//log.Log("       bytemark delete user <user>")
	log.Log("       bytemark delete key [--user=<user>] <public key identifier>")
	log.Log("       bytemark delete server [--force] [---purge] <server>")
	log.Log("       bytemark undelete server <server>")
}*/
func deleteBackup(c *app.Context) (err error) {
	err = c.Client().DeleteBackup(c.VirtualMachineName("server"), c.String("disc"), c.String("backup"))
	if err != nil {
		return
	}
	log.Logf("Backup '%s' deleted successfully", c.String("backup"))
	return
}
