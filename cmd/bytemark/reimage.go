package main

import (
	"fmt"
	"os"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flagsets"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/cliutil"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	isatty "github.com/mattn/go-isatty"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "reimage",
		Usage:     "install a fresh operating system on a server from bytemark's images",
		UsageText: "reimage server [flags] <server>",
		Description: `Image the given server with the specified image, prompting for confirmation.
Specify --force to prevent prompting.

The root password will be output on stdout if the imaging succeeded, otherwise nothing will (and the exit code will be nonzero)`,
		Subcommands: []cli.Command{
			{
				Name:      "server",
				Usage:     "install a fresh operating system on a server from bytemark's images",
				UsageText: "reimage server [flags] <server>",
				Description: `Image the given server with the specified image, prompting for confirmation.
Specify --force to prevent prompting.

The root password will be output on stdout if the imaging succeeded, otherwise nothing will (and the exit code will be nonzero)`,
				Flags: cliutil.ConcatFlags(flagsets.ImageInstallFlags, flagsets.ImageInstallAuthFlags,
					[]cli.Flag{
						forceFlag, cli.GenericFlag{
							Name:  "server",
							Usage: "the server to reimage",
							Value: new(app.VirtualMachineNameFlag),
						},
					}),
				Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
					vmName := c.VirtualMachineName("server")
					imageInstall, defaulted, err := flagsets.PrepareImageInstall(c, true)
					if err != nil {
						return
					}

					if defaulted {
						return c.Help("No image was specified")
					}

					log.Logf("%s will be reimaged with the following. Note that this will wipe all data on the main disc:\r\n\r\n", vmName)
					// don't use ctx.App().ErrWriter / ctx.LogErr here to avoid outputting a password to debug.log
					err = imageInstall.PrettyPrint(os.Stderr, prettyprint.Full)
					if err != nil {
						return
					}

					if !c.Bool("force") && !util.PromptYesNo(c.Prompter(), "Are you certain you wish to continue?") {
						log.Error("Exiting")
						return util.UserRequestedExit{}
					}

					err = c.Client().ReimageVirtualMachine(vmName, imageInstall)
					if err != nil && !isatty.IsTerminal(os.Stdout.Fd()) {
						// by default everything gets output to stdout + debug.log - don't want to output the password to debug log
						_, _ = fmt.Fprintf(os.Stdout, imageInstall.RootPassword)
					}
					return
				}),
			},
		},
	})
}
