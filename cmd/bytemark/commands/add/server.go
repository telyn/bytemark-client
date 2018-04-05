package add

import (
	"fmt"
	"io"

	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/image"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	createServerCmd := cli.Command{
		Name:      "server",
		Usage:     `add a new server with bytemark`,
		UsageText: "add server [flags] <name> [<cores> [<memory [<disc specs>]...]]",
		Description: `Adds a Cloud Server with the given specification, defaulting to a basic server with Symbiosis installed and weekly backups of the first disc.
    
A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to add multiple discs

If --backup is set then a backup of the first disk will be taken at the
frequency specified - never, daily, weekly or monthly. This backup will be free if
it's below a certain threshold of size. By default, a backup is taken every week.
This may cost money if your first disk is larger than the default.
See the price list for more details at http://www.bytemark.co.uk/prices

If --hwprofile-locked is set then the cloud server's virtual hardware won't be changed over time.`,
		Flags: append(app.OutputFlags("server", "object"),
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
				Value: new(util.DiscSpecFlag),
			},
			flags.Force,
			cli.StringFlag{
				Name:  "hwprofile",
				Usage: "The hardware profile to use. Defaults to the current modern profile. See `bytemark profiles` for a list of hardware profiles available.",
			},
			cli.BoolFlag{
				Name:  "hwprofile-locked",
				Usage: "If set, the hardware profile will be 'locked', meaning that when Bytemark updates the hardware profiles your VM will keep its current one.",
			},
			cli.GenericFlag{
				Name:  "ip",
				Value: new(util.IPFlag),
				Usage: "Specify an IPv4 or IPv6 address to use. This will only be useful if you are creating the machine in a private VLAN.",
			},
			cli.GenericFlag{
				Name:  "memory",
				Value: new(util.SizeSpecFlag),
				Usage: "How much memory the server will have available, specified in GiB or with GiB/MiB units. Defaults to 1GiB.",
			},
			cli.GenericFlag{
				Name:  "name",
				Usage: "The new server's name",
				Value: new(app.VirtualMachineNameFlag),
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
		),
		Action: app.Action(args.Optional("name", "cores", "memory", "disc"), with.RequiredFlags("name"), with.Auth, createServer),
	}
	createServerCmd.Flags = append(createServerCmd.Flags, image.ImageInstallFlags...)
	Commands = append(Commands, createServerCmd)
}

// createServer creates a server objec to be created by the brain and sends it.
func createServer(c *app.Context) (err error) {
	name := c.VirtualMachineName("name")
	spec, err := createServerPrepSpec(c)
	if err != nil {
		return
	}

	groupName := name.GroupName()
	err = c.Client().EnsureGroupName(&groupName)
	if err != nil {
		return
	}

	log.Logf("The following server will be created in %s:\r\n", groupName)
	err = spec.PrettyPrint(c.App().Writer, prettyprint.Full)
	if err != nil {
		return err
	}

	// If we're not forcing, prompt. If the prompt comes back false, exit.
	if !c.Bool("force") && !util.PromptYesNo(c.Prompter(), "Are you certain you wish to continue?") {
		log.Error("Exiting.")
		return util.UserRequestedExit{}
	}

	_, err = c.Client().CreateVirtualMachine(groupName, spec)
	if err != nil {
		return err
	}
	vm, err := c.Client().GetVirtualMachine(name)
	if err != nil {
		return
	}
	return c.OutputInDesiredForm(CreatedVirtualMachine{Spec: spec, VirtualMachine: vm})
}

// createServerPrepSpec sets up the server spec by reading in all the flags.
func createServerPrepSpec(c *app.Context) (spec brain.VirtualMachineSpec, err error) {
	noImage := c.Bool("no-image")
	backupFrequency := c.String("backup")

	discs, cores, memory, err := createServerReadArgs(c)
	if err != nil {
		return
	}

	discs, err = createServerPrepDiscs(backupFrequency, discs)
	if err != nil {
		return
	}

	ipspec, err := createServerReadIPs(c)
	if err != nil {
		return
	}

	imageInstall, _, err := image.PrepareImageInstall(c)
	if err != nil {
		return
	}

	stopped := c.Bool("stopped")
	cdrom := c.String("cdrom")

	// if stopped isn't set and a CDROM or image are present, start the server
	autoreboot := !stopped && (!noImage || cdrom != "")

	spec = brain.VirtualMachineSpec{
		VirtualMachine: brain.VirtualMachine{
			Name:                  c.VirtualMachineName("name").VirtualMachine,
			Autoreboot:            autoreboot,
			Cores:                 cores,
			Memory:                memory,
			ZoneName:              c.String("zone"),
			CdromURL:              c.String("cdrom"),
			HardwareProfile:       c.String("hwprofile"),
			HardwareProfileLocked: c.Bool("hwprofile-locked"),
		},
		Discs:   discs,
		IPs:     ipspec,
		Reimage: &imageInstall,
	}
	if noImage {
		spec.Reimage = nil
	}
	return
}

// createServerPrepDiscs checks to see if discs are valid and sets up a backup schedule (if any).
func createServerPrepDiscs(backupFrequency string, discs []brain.Disc) ([]brain.Disc, error) {
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
func createServerReadArgs(c *app.Context) (discs []brain.Disc, cores, memory int, err error) {
	discs = c.Discs("disc")
	cores = c.Int("cores")
	memory = c.Size("memory")
	if memory == 0 {
		memory = 1024
	}
	return
}

// createServerReadIPs reads the IP flags and creates an IPSpec
func createServerReadIPs(c *app.Context) (ipspec *brain.IPSpec, err error) {
	ips := c.IPs("ip")

	if len(ips) > 2 {
		err = c.Help("A maximum of one IPv4 and one IPv6 address may be specified")
		return
	}

	if len(ips) > 0 {
		ipspec = &brain.IPSpec{}

		for _, ip := range ips {
			if ip.To4() != nil {
				if ipspec.IPv4 != "" {
					err = c.Help("A maximum of one IPv4 and one IPv6 address may be specified")
					return
				}
				ipspec.IPv4 = ip.To4().String()
			} else {
				if ipspec.IPv6 != "" {
					err = c.Help("A maximum of one IPv4 and one IPv6 address may be specified")
					return

				}
				ipspec.IPv6 = ip.String()
			}
		}
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

// CreatedVirtualMachine is a struct containing the vm object returned by the VM after creation, and the spec that went into creating it.
// TODO(telyn): move this type into lib/brain?
type CreatedVirtualMachine struct {
	Spec           brain.VirtualMachineSpec `json:"spec"`
	VirtualMachine brain.VirtualMachine     `json:"virtual_machine"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (cvm CreatedVirtualMachine) DefaultFields(f output.Format) string {
	return "Spec, VirtualMachine"
}

// PrettyPrint outputs this created virtual machine in a vaguely nice format to the given writer. detail is ignored.
func (cvm CreatedVirtualMachine) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	_, err = fmt.Fprintf(wr, "cloud server created successfully\r\n")
	if err != nil {
		return
	}

	err = cvm.VirtualMachine.PrettyPrint(wr, prettyprint.Full)
	if err != nil {
		return
	}
	if cvm.Spec.Reimage != nil {
		_, err = fmt.Fprintf(wr, "\r\nRoot password: %s\r\n", cvm.Spec.Reimage.RootPassword)
	} else {
		_, err = fmt.Fprintf(wr, "Machine was not imaged\r\n")
	}
	return
}
