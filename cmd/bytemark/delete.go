package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
)

func init() {
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
			Flags:       []cli.Flag{forceFlag},
			Aliases:     []string{"disk"},
			Action: With(VirtualMachineNameProvider, DiscLabelProvider, AuthProvider, func(c *Context) (err error) {
				if !c.Bool("force") && !util.PromptYesNo("Are you sure you wish to delete this disc? It is impossible to recover.") {
					return util.UserRequestedExit{}
				}

				return global.Client.DeleteDisc(c.VirtualMachineName, *c.DiscLabel)
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
				forceFlag,
			},
			Action: With(GroupProvider, deleteGroup),
		}, {
			Name:        "key",
			Usage:       "deletes the specified key",
			UsageText:   "bytemark delete key [--user <user>] <key>",
			Description: "Keys may be specified as just the comment part or as the whole key. If there are multiple keys with the comment given, an error will be returned",
			Action: With(func(c *Context) (err error) {
				user := global.Config.GetIgnoreErr("user")

				key := strings.Join(c.Args(), " ")
				if key == "" {
					return c.Help("You must specify a key to delete.\r\n")
				}

				err = EnsureAuth()
				if err != nil {
					return
				}

				err = global.Client.DeleteUserAuthorizedKey(user, key)
				if err == nil {
					log.Log("Key deleted successfullly")
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
			},
			Action: With(VirtualMachineProvider, deleteServer),
		}},
	})
}

func deleteServer(c *Context) (err error) {
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

	err = global.Client.DeleteVirtualMachine(c.VirtualMachineName, purge)
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
	running := 0
	for _, vm := range servers {
		if vm.PowerOn {
			running++
		}
	}
	return running
}

func deleteGroup(c *Context) (err error) {
	recursive := c.Bool("recursive")
	if len(c.Group.VirtualMachines) > 0 && recursive {
		prompt := fmt.Sprintf("The group '%s' has %d servers in it, these servers will be irrevocably deleted", c.GroupName.Group, len(c.Group.VirtualMachines))
		if countRunning(c.Group) != 0 {
			prompt = fmt.Sprintf("The group '%s' has %d currently-running servers in it, these servers will be forcibly stopped and irrevocably deleted", c.GroupName.Group, running)
		}

		if !c.Bool("force") && !util.PromptYesNo(prompt+" - are you sure you wish to delete this group?") {
			return util.UserRequestedExit{}
		}
		err = recursiveDeleteGroup(c.GroupName, c.Group)
		if err != nil {
			return
		}
	} else if !recursive {
		err = &util.WontDeleteNonEmptyGroupError{Group: c.GroupName}
		return
	}
	err = global.Client.DeleteGroup(c.GroupName)
	if err != nil {
		log.Logf("Group %s deleted successfully.\r\n", c.GroupName.String())
	}
	return
}

func recursiveDeleteGroup(name *lib.GroupName, group *brain.Group) error {
	log.Log("", "")
	vmn := lib.VirtualMachineName{Group: name.Group, Account: name.Account}
	for _, vm := range group.VirtualMachines {
		vmn.VirtualMachine = vm.Name
		err := global.Client.DeleteVirtualMachine(&vmn, true)
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
