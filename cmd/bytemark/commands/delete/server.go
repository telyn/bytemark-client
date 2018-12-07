package delete

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flagsets"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "server",
		Usage:       "delete the given server",
		UsageText:   `delete server [--purge] <server name>`,
		Description: "Deletes the given server. Deleted servers still exist and can be restored. To ensure a server is fully deleted, use the --purge flag.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "purge",
				Usage: "If set, the server will be irrevocably deleted.",
			},
			flagsets.Force,
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to delete",
				Value: new(flags.VirtualMachineNameFlag),
			},
		},
		Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.VirtualMachine("server"), func(c *app.Context) (err error) {
			purge := c.Bool("purge")
			vm := c.VirtualMachine

			if vm.Deleted && !purge {
				c.LogErr("Server %s has already been deleted.\r\nIf you wish to permanently delete it, add --purge", vm.Hostname)
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

			vmName := flags.VirtualMachineName(c, "server")
			err = c.Client().DeleteVirtualMachine(vmName, purge)
			if err != nil {
				return
			}
			if purge {
				c.Log("Server %s purged successfully.", vm.Hostname)
			} else {
				c.Log("Server %s deleted successfully.", vm.Hostname)
			}
			return
		}),
	})
}
