package add

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
	vmdefaultFlags := append(app.OutputFlags("vm default", "object"),
		flags.ImageInstallFlags...)
	vmdefaultFlags = append(vmdefaultFlags, flags.ServerSpecFlags...)
	Commands = append(Commands, cli.Command{
		Name:      "vm default",
		Usage:     "adds a new VM default",
		UsageText: "--admin add vm default <default name>",
		Description: `adds a new VM Default to the current account, which can be specified as either public or private.
The server settings can be specified for the vm default with aditional flags

--default-name (and the <default name> positional argument) is an identifier for the default, not a default name for servers created based upon it.

A disc spec looks like the following: grade:size. The grade field is optional and will default to sata.
Multiple --disc flags can be used to add multiple discs to the VM Default

If --backup is set then a backup of the first disk will be taken at the
frequency specified - never, daily, weekly or monthly. If not specified the backup will default to weekly.`,
		Flags: append(vmdefaultFlags,
			cli.StringFlag{
				Name:  "default-name",
				Usage: "The name of the VM default to add",
			},
			cli.BoolFlag{
				Name:  "public",
				Usage: "If the VM default should be made public or not",
			},
			cli.GenericFlag{
				Name:  "account",
				Usage: "the account to add the default to (will use 'bytemark' if unset)",
				Value: new(app.AccountNameFlag),
			},
		),
		Action: app.Action(args.Optional("default-name"), with.RequiredFlags("default-name"), with.Auth, func(c *app.Context) (err error) {
			accountName := c.String("account")
			if !c.IsSet("account") {
				accountName = "bytemark"
			}
			account, err := c.Client().GetAccount(accountName)
			if err != nil {
				return
			}
			spec, err := flags.PrepareServerSpec(c)
			if err != nil {
				return
			}

			vmd := brain.VirtualMachineDefault{
				AccountID:      account.BrainID,
				Name:           c.String("default-name"),
				Public:         c.Bool("public"),
				ServerSettings: spec,
			}

			err = brainRequests.CreateVMDefault(c.Client(), vmd)
			if err != nil {
				return
			}
			return
		}),
	})
}
