package flagsets

import (
	"fmt"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
)

// ServerSpecFlags are a set of flags common to commands which set up a server
// spec.  (most notably add vm default and add server)
var ServerSpecFlags = []cli.Flag{
	cli.IntFlag{
		Name:  "cores",
		Value: 1,
		Usage: "Number of CPU cores",
	},
	cli.StringFlag{
		Name:  "cdrom",
		Usage: "URL pointing to an ISO which will be attached to the cloud server as a CD",
	},
	cli.GenericFlag{
		Name:  "disc",
		Usage: "One or more disc specifications. Defaults to a single 25GiB sata-grade disc",
		Value: new(flags.DiscSpecFlag),
	},
	Force,
	cli.StringFlag{
		Name:  "hwprofile",
		Usage: "The hardware profile to use. Defaults to the current modern profile. See `bytemark profiles` for a list of hardware profiles available.",
	},
	cli.BoolFlag{
		Name:  "hwprofile-locked",
		Usage: "If set, the hardware profile will be 'locked', meaning that when Bytemark updates the hardware profiles your VM will keep its current one.",
	},
	cli.GenericFlag{
		Name:  "memory",
		Value: new(flags.SizeSpecFlag),
		Usage: "How much memory the server will have available, specified in GiB or with GiB/MiB units. Defaults to 1GiB.",
	},
	cli.BoolFlag{
		Name:  "no-image",
		Usage: "Specifies that the server should not be imaged.",
	},
	cli.StringFlag{
		Name:  "backup",
		Usage: "Add a backup schedule for the first disk at the given frequency (daily, weekly, monthly, or never)",
		Value: "weekly",
	},
	cli.BoolFlag{
		Name:  "stopped",
		Usage: "If set, the server will not be started, even to image it.",
	},
	cli.StringFlag{
		Name:  "zone",
		Usage: "Which zone the server will be created in. See `bytemark zones` for the choices.",
	},
}

// PrepareServerSpec sets up the server spec by reading in all the flags.
// set authentication to true when you wanna read in authorized-keys/-file and root-password too.
func PrepareServerSpec(c *app.Context, authentication bool) (spec brain.VirtualMachineSpec, err error) {
	noImage := c.Bool("no-image")
	backupFrequency := c.String("backup")

	discs, cores, memory, err := prepareServerReadArgs(c)
	if err != nil {
		return
	}

	discs, err = prepareDiscs(backupFrequency, discs)
	if err != nil {
		return
	}

	imageInstall, _, err := PrepareImageInstall(c, authentication)
	if err != nil {
		return
	}

	stopped := c.Bool("stopped")
	cdrom := c.String("cdrom")

	// if stopped isn't set and a CDROM or image are present, start the server
	autoreboot := !stopped && (!noImage || cdrom != "")

	spec = brain.VirtualMachineSpec{
		VirtualMachine: brain.VirtualMachine{
			Autoreboot:            autoreboot,
			Cores:                 cores,
			Memory:                memory,
			ZoneName:              c.String("zone"),
			CdromURL:              c.String("cdrom"),
			HardwareProfile:       c.String("hwprofile"),
			HardwareProfileLocked: c.Bool("hwprofile-locked"),
		},
		Discs:   discs,
		Reimage: &imageInstall,
	}
	if noImage {
		spec.Reimage = nil
	}
	return
}

// prepareDiscs checks to see if discs are valid and sets up a backup schedule (if any).
func prepareDiscs(backupFrequency string, discs []brain.Disc) ([]brain.Disc, error) {
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

	interval, err := BackupScheduleIntervalFromWords(backupFrequency)
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

// prepareServerReadArgs sets up the initial defaults, reads in the --disc, --cores and --memory flags
func prepareServerReadArgs(c *app.Context) (discs []brain.Disc, cores, memory int, err error) {
	discs = flags.Discs(c, "disc")
	cores = c.Int("cores")
	memory = flags.Size(c, "memory")
	if memory == 0 {
		memory = 1024
	}
	return
}

// BackupScheduleIntervalFromWords deteremines the backup interval
func BackupScheduleIntervalFromWords(words string) (freq int, err error) {
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
