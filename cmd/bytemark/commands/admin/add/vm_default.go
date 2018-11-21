package add

import (
	"fmt"
	"io"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/image"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
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
		Flags: append(app.OutputFlags("vmdefault", "object"),
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

	spec, err := createVMDPrepSpec(c)
	if err != nil {
		return
	}

	err = brainRequests.CreateVMDefault(c.Client(), name, public, spec)
	if err != nil {
		return err
	}

	return c.OutputInDesiredForm(CreatedVMDefault{Spec: spec})
}

// createServerPrepSpec sets up the server spec by reading in all the flags.
func createVMDPrepSpec(c *app.Context) (spec brain.VMDefaultSpec, err error) {
	backupFrequency := c.String("backup")

	discs, cores, memory, err := createVMDReadArgs(c)
	if err != nil {
		return
	}

	discs, err = createVMDPrepDiscs(backupFrequency, discs)
	if err != nil {
		return
	}

	imageInstall, _, err := image.PrepareImageInstall(c)
	if err != nil {
		return
	}

	spec = brain.VMDefaultSpec{
		VMDefault: brain.VMDefault{
			Name:            c.String("vm-name"),
			Cores:           cores,
			Memory:          memory,
			ZoneName:        c.String("zone"),
			CdromURL:        c.String("cdrom"),
			HardwareProfile: c.String("hwprofile"),
		},
		Discs:   discs,
		Reimage: &imageInstall,
	}
	return
}

// createServerPrepDiscs checks to see if discs are valid and sets up a backup schedule (if any).
func createVMDPrepDiscs(backupFrequency string, discs []brain.Disc) ([]brain.Disc, error) {
	if len(discs) == 0 {
		discs = append(discs, brain.Disc{Size: 25600})
	}

	for i := range discs {
		d, discErr := discs[i].Validate()
		if discErr != nil {
			return discs, discErr
		}
		discs[i] = *d
	}

	interval, err := backupScheduleIntervalFromWords(backupFrequency)
	if err != nil {
		return discs, err
	}

	if interval > 0 {
		if len(discs) > 0 {
			bs := defaultBackupSchedule()
			bs.Interval = interval
			discs[0].BackupSchedules = brain.BackupSchedules{bs}
		}
	}
	return discs, nil
}

// createServerReadArgs sets up the initial defaults, reads in the --disc, --cores and --memory flags
func createVMDReadArgs(c *app.Context) (discs []brain.Disc, cores, memory int, err error) {
	discs = c.Discs("disc")
	cores = c.Int("cores")
	memory = c.Size("memory")
	if memory == 0 {
		memory = 1024
	}
	return
}

// backupScheduleIntervalFromWords deteremines the backup interval
func backupScheduleIntervalFromWords(words string) (freq int, err error) {
	switch words {
	case "daily":
		freq = 86400
	case "weekly":
		freq = 7 * 86400
	case "never":
		// the brain will reject a -1 - so even if the frequency accidentally
		// makes it to the brain the schedule won't be made
		freq = -1
	default:
		err = fmt.Errorf("invalid backup frequency '%s'", words)
	}
	return

}

// defaultBackupSchedule returns a schedule that will backup every week (well - every 604800 seconds)
func defaultBackupSchedule() brain.BackupSchedule {
	return brain.BackupSchedule{
		StartDate: "",
		Interval:  7 * 86400,
		Capacity:  1,
	}
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
