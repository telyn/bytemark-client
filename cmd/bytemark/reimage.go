package main

import (
	"fmt"
	"os"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli"
)

var imageInstallFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "authorized-keys",
		Usage: "Public keys that will be authorised to log in as root, separated by newlines.",
	},
	cli.GenericFlag{
		Name:  "authorized-keys-file",
		Usage: "Local file to read the --authorized-keys from",
		Value: new(util.FileFlag),
	},
	cli.StringFlag{
		Name:  "firstboot-script",
		Usage: "Script which runs on the server's first boot after imaging.",
	},
	cli.GenericFlag{
		Name:  "firstboot-script-file",
		Usage: "Local file to read the firstboot script from.",
		Value: new(util.FileFlag),
	},
	cli.StringFlag{
		Name:  "image",
		Usage: "Image to install on the server. See `bytemark images` for the list of available images.",
	},
	cli.StringFlag{
		Name:  "root-password",
		Usage: "Password for the root/Administrator user. If unset, will be randomly generated.",
	},
}

func init() {
	commands = append(commands, cli.Command{
		Name:      "reimage",
		Usage:     "install a fresh operating system on a server from bytemark's images",
		UsageText: "bytemark reimage [flags] <server>",
		Description: `Image the given server with the specified image, prompting for confirmation.
Specify --force to prevent prompting.

The root password will be output on stdout if the imaging succeeded, otherwise nothing will (and the exit code will be nonzero)
	    `,
		Flags: append(imageInstallFlags, forceFlag, cli.GenericFlag{
			Name:  "server",
			Usage: "the server to reimage",
			Value: new(app.VirtualMachineNameFlag),
		}),
		Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
			vmName := c.VirtualMachineName("server")
			imageInstall, defaulted, err := prepareImageInstall(c)
			if err != nil {
				return
			}

			if defaulted {
				return c.Help("No image was specified")
			}

			log.Logf("%s will be reimaged with the following. Note that this will wipe all data on the main disc:\r\n\r\n", vmName)
			err = imageInstall.PrettyPrint(os.Stderr, prettyprint.Full)
			if err != nil {
				return
			}

			if !c.Bool("force") && !util.PromptYesNo("Are you certain you wish to continue?") {
				log.Error("Exiting")
				return util.UserRequestedExit{}
			}

			err = c.Client().ReimageVirtualMachine(vmName, imageInstall)
			if err != nil && !isatty.IsTerminal(os.Stdout.Fd()) {
				fmt.Fprintf(os.Stdout, imageInstall.RootPassword)
			}
			return
		}),
	})
}

func prepareImageInstall(c *app.Context) (imageInstall brain.ImageInstall, defaulted bool, err error) {
	image := c.String("image")
	firstbootScript := c.String("firstboot-script")
	firstbootScriptFile := c.FileContents("firstboot-script-file")
	pubkeys := c.String("authorized-keys")
	pubkeysFile := c.FileContents("authorized-keys-file")
	rootPassword := c.String("root-password")

	if image == "" {
		image = "symbiosis"
		defaulted = true
	}

	if !c.Bool("force") {
		var exists bool
		exists, err = imageExists(c, image)
		if err != nil {
			return
		}
		if !exists {
			err = fmt.Errorf("No visible image '%s' - check your spelling or use --force if certain", image)
		}
	}

	if pubkeysFile != "" {
		if pubkeys != "" {
			pubkeys += "\r\n" + pubkeysFile
		} else {
			pubkeys = pubkeysFile
		}
	}

	if firstbootScript == "" {
		firstbootScript = firstbootScriptFile
	}

	if rootPassword == "" {
		rootPassword = util.GeneratePassword()
	}

	return brain.ImageInstall{
		Distribution:    image,
		FirstbootScript: firstbootScript,
		PublicKeys:      pubkeys,
		RootPassword:    rootPassword,
	}, defaulted, err
}

func imageExists(c *app.Context, name string) (exists bool, err error) {
	defs, err := c.Client().ReadDefinitions()
	if err != nil {
		return
	}
	for _, image := range defs.Distributions {
		if image == name {
			exists = true
			return
		}
	}
	return
}
