package add

import (
	"fmt"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/image"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:    "vm default",
		Aliases: []string{"vm-default"},
		Usage:   "adds a new VM Default",
		UsageText: "--admin add vm default <name> <public> [<cores>[<memory>[<disc-specs>]...]]",
		Description: `adds a new VM Default to the current account, which can be specified as either public or private.
  					  the server settings can be specified for the vm default with flags`,
		// TODO (tom): add to description

		Flags: append(app.OutputFlags("vmdefault", "object"),
			cli.StringFlag{
				Name:  "name",
				Usage: "The name of the VM Default to add",
			},
			cli.BoolFlag{
				Name:  "public",
				Usage: "If the VM Default should be made public or not",
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
			cli.GenericFlag{
				Name:  "disc",
				Usage: "One or more disc specifications. Defaults to a single 25GiB sata-grade disc",
				Value: new(util.DiscSpecFlag),
			},
			cli.GenericFlag{
				Name:  "vm-name",
				Usage: "The name of the VM",
				Value: new(app.VirtualMachineNameFlag),
			},
			cli.StringFlag{
				Name:  "image",
				Usage: "Image to install on the server. See `bytemark images` for the list of available images.",
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
		),
		Action: app.Action(args.Optional("name", "public", "cores", "memory", "disc"), with.RequiredFlags("name", "public"), with.Auth, createVMDefault),
	})
}


// createVmDefault creates a server object to be created by the brain and sends it.
func createVMDefault(c *app.Context) (err error) {
	name := c.String("name")
	public := c.Bool("public")

	if name == "" {
		name = "vm default"
	}

	serverSettings, err := createVMDPrepSpec(c)
	if err != nil {
		return
	}

	// add pretty print

	err = c.Client().CreateVMDefault(name, public, serverSettings)
	if err != nil {
		return err
	}

	// add pretty print for created vm default
	return
}

// createServerPrepSpec sets up the server spec by reading in all the flags.
func createVMDPrepSpec(c *app.Context) (spec brain.VmDefaultSpec, err error) {
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

	name := c.VirtualMachineName("vm-name")

	spec = brain.VmDefaultSpec{
		VmDefault: brain.VMDefault{
			Name:                  name.VirtualMachine,
			Cores:                 cores,
			Memory:                memory,
			ZoneName:              c.String("zone"),
			CdromURL:              c.String("cdrom"),
			HardwareProfile:       c.String("hwprofile"),
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
// starting from midnight tonight.
func defaultBackupSchedule() brain.BackupSchedule {
	tomorrow := time.Now().Add(24 * time.Hour)
	y, m, d := tomorrow.Date()
	midnightTonight := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	defaultStartDate := midnightTonight.Format("2006-01-02 15:04:05 MST")
	return brain.BackupSchedule{
		StartDate: defaultStartDate,
		Interval:  7 * 86400,
		Capacity:  1,
	}
}