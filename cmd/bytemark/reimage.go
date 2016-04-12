package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "reimage",
		Usage:     "Install a fresh operating system on a server from bytemark's images",
		UsageText: "bytemark reimage [flags] --image <image name> <server>",
		Description: `Image the given server with the specified image, prompting for confirmation.
If the --image flag is not specified, will prompt with a list.
Specify --force to prevent prompting.

The root password will be the only thing output on stdout - good for scripts!
	    `,
		Flags: []cli.Flag{
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
				Name:  "public-keys",
				Usage: "Public keys that will be authorised to log in as root, separated by newlines.",
			},
			cli.GenericFlag{
				Name:  "public-keys-file",
				Usage: "Local file to read the public keys from",
				Value: new(util.FileFlag),
			},
			cli.StringFlag{
				Name:  "root-password",
				Usage: "Password for the root/Administrator user. If unset, will be randomly generated.",
			},
		},
		Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
			imageInstall, defaulted, err := prepareImageInstall(c)
			if err != nil {
				return err
			}
			if defaulted {
				return c.Help("No image was specified")
			}

			return global.Client.ReimageVirtualMachine(c.VirtualMachineName, imageInstall)
		}),
	})
}

func prepareImageInstall(c *Context) (imageInstall *lib.ImageInstall, defaulted bool, err error) {
	image := c.String("image")
	firstbootScript := c.String("firstboot-script")
	firstbootScriptFile := c.FileContents("firstboot-script-file")
	pubkeys := c.String("public-keys")
	pubkeysFile := c.FileContents("public-keys-file")
	rootPassword := c.String("root-password")

	if image == "" {
		image = "symbiosis"
		defaulted = true
	}
	if pubkeysFile != "" {
		pubkeys += "\r\n" + pubkeysFile
	}

	if firstbootScript == "" {
		firstbootScript = firstbootScriptFile
	}

	if rootPassword == "" {
		rootPassword = util.GeneratePassword()
	}

	return &lib.ImageInstall{
		Distribution:    image,
		FirstbootScript: firstbootScript,
		PublicKeys:      pubkeys,
		RootPassword:    rootPassword,
	}, defaulted, nil
}
