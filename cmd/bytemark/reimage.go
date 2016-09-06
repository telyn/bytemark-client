package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"os"
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
If the --image flag is not specified, will prompt with a list.
Specify --force to prevent prompting.

The root password will be the only thing output on stdout - good for scripts!
	    `,
		Flags: append(imageInstallFlags, forceFlag),
		Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
			imageInstall, defaulted, err := prepareImageInstall(c)
			if err != nil {
				return err
			}

			if defaulted {
				return c.Help("No image was specified")
			}

			log.Logf("%s will be reimaged with the following. Note that this will wipe all data on the main disc:\r\n\r\n", c.VirtualMachineName.String())
			lib.FormatImageInstall(os.Stderr, imageInstall, "imageinstall")

			if !c.Bool("force") && !util.PromptYesNo("Are you certain you wish to continue?") {
				log.Error("Exiting")
				return util.UserRequestedExit{}
			}

			return global.Client.ReimageVirtualMachine(c.VirtualMachineName, imageInstall)
		}),
	})
}

func prepareImageInstall(c *Context) (imageInstall *brain.ImageInstall, defaulted bool, err error) {
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

	return &brain.ImageInstall{
		Distribution:    image,
		FirstbootScript: firstbootScript,
		PublicKeys:      pubkeys,
		RootPassword:    rootPassword,
	}, defaulted, nil
}
