package add

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
	vmdefaultFlags := append(app.OutputFlags("vmdefault", "object"),
		flags.ImageInstallFlags...)
	vmdefaultFlags = append(vmdefaultFlags, flags.ServerSpecFlags...)
	Commands = append(Commands, cli.Command{
		Name:      "vm default",
		Aliases:   []string{"vm-default"},
		Usage:     "adds a new VM Default",
		UsageText: "--admin add vm default <name> <public> [<cores>[<memory>[<disc-specs>]...]]",
		Description: `adds a new VM Default to the current account, which can be specified as either public or private.
  					  the server settings can be specified for the vm default with aditional flags

The VM Default name will be a name for the whole VM Default spcification. This is what the VMDefault will be referred as.

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
			cli.StringFlag{
				Name:  "cdrom",
				Usage: "URL pointing to an ISO which will be attached to the cloud server as a CD",
			},
			cli.IntFlag{
				Name:  "cores",
				Value: 1,
				Usage: "Number of CPU cores",
			},
			cli.GenericFlag{
				Name:  "memory",
				Value: new(util.SizeSpecFlag),
				Usage: "How much memory the server will have available, specified in GiB or with GiB/MiB units. Defaults to 1GiB.",
			},
			cli.StringFlag{
				Name:  "vm-name",
				Usage: "The name of the VM",
			},
			flags.Force,
			cli.StringFlag{
				Name:  "image",
				Usage: "Image to install on the server. See `bytemark images` for the list of available images.",
			},
			cli.StringFlag{
				Name:  "hwprofile",
				Usage: "The hardware profile to use. Defaults to the current modern profile. See `bytemark profiles` for a list of hardware profiles available.",
			},
			cli.StringFlag{
				Name:  "backup",
				Usage: "Add a backup schedule for the first disk at the given frequency (daily, weekly, monthly, or never)",
				Value: "weekly",
			},
			cli.StringFlag{
				Name:  "zone",
				Usage: "Which zone the server will be created in. See `bytemark zones` for the choices.",
			},
			cli.GenericFlag{
				Name:  "disc",
				Usage: "One or more disc specifications. Defaults to a single 25GiB sata-grade disc",
				Value: new(util.DiscSpecFlag),
			},
			cli.StringFlag{
				Name:  "firstboot-script",
				Usage: "The firstboot script the VM Default should use.",
			},
		),
		Action: app.Action(args.Optional("name", "public", "cores", "memory", "disc"), with.RequiredFlags("name", "public"), with.Auth, createVMDefault),
	})
}

// createVMDefault creates a server object to be created by the brain and sends it.
func createVMDefault(c *app.Context) (err error) {
	name := c.String("name")
	public := c.Bool("public")

	if name == "" {
		name = "vm default"
	}

	spec, err := flags.PrepareServerSpec(c)
	if err != nil {
		return
	}

	err = brainRequests.CreateVMDefault(c.Client(), name, public, spec)
	if err != nil {
		return err
	}

	return c.OutputInDesiredForm(CreatedVMDefault{Spec: spec})
}

// CreatedVMDefault is a struct containing the vmd object returned by the VM Default after creation,
// and the spec that went into creating it.
type CreatedVMDefault struct {
	Spec      brain.VMDefaultSpec `json:"server_settings"`
	VMDefault brain.VMDefault     `json:"virtual_machine"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (c CreatedVMDefault) DefaultFields(f output.Format) string {
	return "Spec, VMDefault"
}

// PrettyPrint outputs this created vm default in a vaguely nice format to the given writer. detail is ignored.
func (c CreatedVMDefault) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	err = c.VMDefault.PrettyPrint(wr, prettyprint.Full)
	if err != nil {
		return
	}
	return
}
