package add

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	commandsUtil "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/util"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:    "discs",
		Aliases: []string{"disc", "disk", "disks"},
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "disc",
				Usage: "A disc to add. You can specify as many discs as you like by adding more --disc flags.",
				Value: new(util.DiscSpecFlag),
			},
			flags.ForceFlag,
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to add the disc to",
				Value: new(app.VirtualMachineNameFlag),
			},
		},
		Usage:     "add virtual discs attached to one of your cloud servers",
		UsageText: "add discs [--disc <disc spec>]... <cloud server>",
		Description: `A disc spec looks like the following: label:grade:size
The label and grade fields are optional. If grade is empty, defaults to sata.
If there are two fields, they are assumed to be grade and size.
Multiple --disc flags can be used to add multiple discs`,
		Action: app.Action(args.Optional("server", "cores", "memory", "disc"), with.Auth, createDiscs),
	})
}

// createDiscs adds the disc(s) to the speicified server
func createDiscs(c *app.Context) (err error) {
	discs := c.Discs("disc")

	for i := range discs {
		d, err := discs[i].Validate()
		if err != nil {
			return err
		}
		discs[i] = *d
	}
	vmName := c.VirtualMachineName("server")

	log.Logf("Adding %d discs to %s:\r\n", len(discs), vmName)
	for _, d := range discs {
		log.Logf("    %dGiB %s...", d.Size/1024, d.StorageGrade)
		err := c.Client().CreateDisc(vmName, d)
		if err != nil {
			log.Errorf("failure! %v\r\n", err.Error())
		} else {
			log.Log("success!")
		}
	}
	return
}
