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
		Aliases:   []string{"vm-default"},
		Usage:     "adds a new VM Default",
		UsageText: "--admin add vm default <name>",
		Description: `adds a new VM Default to the current account, which can be specified as either public or private.
  					  the server settings can be specified for the vm default with aditional flags

--name is an identifier for the default, not a default name for servers created based upon it.

A disc spec looks like the following: grade:size. The grade field is optional and will default to sata.
Multiple --disc flags can be used to add multiple discs to the VM Default

If --backup is set then a backup of the first disk will be taken at the
frequency specified - never, daily, weekly or monthly. If not specified the backup will default to weekly.`,
		Flags: append(vmdefaultFlags,
			cli.StringFlag{
				Name:  "name",
				Usage: "The name of the VM Default to add",
			},
			cli.BoolFlag{
				Name:  "public",
				Usage: "If the VM Default should be made public or not",
			},
		),
		Action: app.Action(args.Optional("name", "public", "cores", "memory", "disc"), with.RequiredFlags("name"), with.Auth, createVMDefault),
	})
}

// createVMDefault creates a server object to be created by the brain and sends it.
func createVMDefault(c *app.Context) (err error) {
	name := c.String("name")
	public := c.Bool("public")

	if name == "" {
		name = "vm-default"
	}
	spec, err := flags.PrepareServerSpec(c)
	if err != nil {
		return
	}

	vmd := brain.VirtualMachineDefault{
		Name:           name,
		Public:         public,
		ServerSettings: spec,
	}

	err = brainRequests.CreateVMDefault(c.Client(), vmd)
	if err != nil {
		return err
	}
	return
}
