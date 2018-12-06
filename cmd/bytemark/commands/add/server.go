package add

import (
	"fmt"
	"io"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flagsets"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/cliutil"
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

The server name can be used to specify which group and account the server should be created in, for example myserver.group1.myaccount.
    
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
		Flags: cliutil.ConcatFlags(app.OutputFlags("server", "object"),
			flagsets.ServerSpecFlags, flagsets.ImageInstallFlags, flagsets.ImageInstallAuthFlags,
			[]cli.Flag{
				cli.GenericFlag{
					Name:  "name",
					Usage: "The new server's name",
					Value: new(app.VirtualMachineNameFlag),
				},
				cli.GenericFlag{
					Name:  "ip",
					Value: new(util.IPFlag),
					Usage: "Specify an IPv4 or IPv6 address to use. This will only be useful if you are creating the machine in a private VLAN.",
				},
			}),
		Action: app.Action(args.Optional("name", "cores", "memory", "disc"), with.RequiredFlags("name"), with.Auth, createServer),
	}
	Commands = append(Commands, createServerCmd)
}

// createServer creates a server objec to be created by the brain and sends it.
func createServer(c *app.Context) (err error) {
	name := c.VirtualMachineName("name")
	spec, err := flagsets.PrepareServerSpec(c, true)
	if err != nil {
		return
	}
	spec.VirtualMachine.Name = name.VirtualMachine
	// add a fake Hostname so that FullName() works for the prompt below
	spec.VirtualMachine.Hostname = name.String()

	ipspec, err := createServerReadIPs(c)
	if err != nil {
		return
	}
	spec.IPs = ipspec

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

	// Clear hostname - it was there for the pre-flight check
	spec.VirtualMachine.Hostname = ""

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
