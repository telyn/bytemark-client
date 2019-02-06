package delete

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flagsets"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "group",
		Usage:     "deletes the given group",
		UsageText: "delete group [--force] [--recursive] <group name>",
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
				Value: new(flags.GroupNameFlag),
			},
			flagsets.Force,
		},
		Action: app.Action(args.Optional("group"), with.RequiredFlags("group"), with.Group("group"), func(ctx *app.Context) (err error) {
			recursive := ctx.Bool("recursive")
			groupName := flags.GroupName(ctx, "group")
			if len(ctx.Group.VirtualMachines) > 0 {
				if !recursive {
					err = util.WontDeleteGroupWithVMsError{Group: groupName}
					return
				}

				err = promptForRecursiveDeleteGroup(ctx)
				if err != nil {
					return
				}

				err = deleteVmsInGroup(ctx, groupName, ctx.Group)
				if err != nil {
					return
				}
			} else if !ctx.Bool("force") && !util.PromptYesNo(ctx.Prompter(), fmt.Sprintf("Are you sure you wish to delete the %s group?", groupName)) {
				return util.UserRequestedExit{}
			}
			err = ctx.Client().DeleteGroup(groupName)
			if err == nil {
				ctx.Log("\nGroup %s deleted successfully.", groupName.String())
			}
			return
		}),
	})
}

func promptForRecursiveDeleteGroup(ctx *app.Context) error {
	prompt := fmt.Sprintf("The group '%s' has %d servers in it which will be irrevocably deleted", ctx.Group.Name, len(ctx.Group.VirtualMachines))
	running := countRunning(ctx.Group)

	if running != 0 {
		stopped := len(ctx.Group.VirtualMachines) - running
		andStopped := ""
		if stopped > 0 {
			andStopped = fmt.Sprintf("and %d stopped ", stopped)
		}
		prompt = fmt.Sprintf("The group '%s' has %d currently-running %sservers in it which will be forcibly stopped and irrevocably deleted", ctx.Group.Name, running, andStopped)
	}

	if !ctx.Bool("force") && !util.PromptYesNo(ctx.Prompter(), prompt+" - are you sure you wish to delete this group?") {
		return util.UserRequestedExit{}
	}
	return nil
}

func deleteVmsInGroup(ctx *app.Context, name pathers.GroupName, group *brain.Group) error {
	ctx.Log("\nPurging all VMs in %s...", name)

	recurseErr := util.RecursiveDeleteGroupError{Group: name, Errors: map[string]error{}}

	vmn := pathers.VirtualMachineName{GroupName: name}
	for _, vm := range group.VirtualMachines {
		vmn.VirtualMachine = vm.Name
		ctx.Logf("%s...", vm.Name)

		err := ctx.Client().DeleteVirtualMachine(vmn, true)
		if err != nil {
			ctx.Log("failed")
			recurseErr.Errors[vm.Name] = err
		}
		ctx.Log("deleted")
	}
	if len(recurseErr.Errors) > 0 {
		return recurseErr
	}
	return nil
}

func countRunning(group *brain.Group) (running int) {
	for _, vm := range group.VirtualMachines {
		if vm.PowerOn {
			running++
		}
	}
	return running
}
