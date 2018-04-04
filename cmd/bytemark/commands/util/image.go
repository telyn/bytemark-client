package util

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
)

// ImageInstallFlags is common to a few commands and contains additonial flags
// that get stitched the commands that use it.
var ImageInstallFlags = []cli.Flag{
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

// PrepareImageInstall is a funcion that prepares an image to be imaged on a server.
func PrepareImageInstall(c *app.Context) (imageInstall brain.ImageInstall, defaulted bool, err error) {
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
