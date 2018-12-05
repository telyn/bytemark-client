package main

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {

	commands = append(commands, cli.Command{
		Name:      "delete",
		Usage:     "delete a given server, disc, group, account or key",
		UsageText: "delete account|disc|group|key|server",
		Description: `delete a given server, disc, group, account or key

   Only empty groups and accounts can be deleted.

   The restore server command may be used to restore a deleted (but not purged) server to its state prior to deletion.`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "disc",
			Usage:       "delete the given disc",
			UsageText:   "delete disc [--server <virtual machine name> --label <disc label>] | [--id <disc ID>]",
			Description: "Deletes the given disc. To find out a disc's label you can use the `bytemark show server` command or `bytemark list discs` command.",
			Flags: []cli.Flag{
				forceFlag,
				cli.StringFlag{
					Name:  "label",
					Usage: "the disc to delete, must provide a server too",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server whose disc you wish to delete, must provide a label too",
					Value: new(flags.VirtualMachineName),
				},
				cli.StringFlag{
					Name:  "id",
					Usage: "the ID of the disc to delete",
				},
			},
			Aliases: []string{"disk"},
			Action: app.Action(args.Optional("server", "label", "id"), with.Auth, func(c *app.Context) (err error) {
				if !c.Bool("force") && !util.PromptYesNo(c.Prompter(), "Are you sure you wish to delete this disc? It is impossible to recover.") {
					return util.UserRequestedExit{}
				}
				vmName := c.VirtualMachineName("server")
				discLabel := c.String("label")
				discID := c.String("id")

				if discID != "" {
					return brainRequests.DeleteDiscByID(c.Client(), discID)
				} else if vmName.String() != "" && discLabel != "" {
					return c.Client().DeleteDisc(vmName, discLabel)
				} else {
					return fmt.Errorf("Please include both --server and --label flags or provide --id")
				}
			}),
		}, {
			Name:        "key",
			Usage:       "deletes the specified key",
			UsageText:   "delete key [--user <user>] <key>",
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

				err = brainRequests.DeleteUserAuthorizedKey(c.Client(), user, key)
				if err == nil {
					log.Log("Key deleted successfully")
				}
				return
			}),
		}, {
			Name:        "server",
			Usage:       "delete the given server",
			UsageText:   `delete server [--purge] <server name>`,
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
					Value: new(flags.VirtualMachineName),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.VirtualMachine("server"), deleteServer),
		}, {
			Name:        "backup",
			Usage:       "delete the given backup",
			UsageText:   `delete backup <server name> <disc label> <backup label>`,
			Description: "Deletes the given backup. Backups cannot be recovered after deletion.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "disc",
					Usage: "the disc to delete a backup of",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to delete a backup from",
					Value: new(flags.VirtualMachineName),
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

	if !c.Bool("force") && !util.PromptYesNo(c.Prompter(), fstr) {
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

func deleteBackup(c *app.Context) (err error) {
	err = c.Client().DeleteBackup(c.VirtualMachineName("server"), c.String("disc"), c.String("backup"))
	if err != nil {
		return
	}
	log.Logf("Backup '%s' deleted successfully", c.String("backup"))
	return
}
