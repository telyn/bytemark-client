package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"flag"
	"github.com/codegangsta/cli"
	"io/ioutil"
)

func init() {
	commands = append(commands, cli.Command{
		Name: "reimage",
		Description: `flags available
    --firstboot-script-file <file name> - a script that will be run on first boot
    --force - disables the confirmation prompt
    --image <image name> - specify what to image the server with. Default is 'symbiosis'. See bytemark images for a list of available images
    --public-keys <keys> (newline seperated)
    --public-keys-file <file> (will be read & appended to --public-keys)
    --root-password <password> (if not set, will be randomly generated)

Image the given server with the specified image, prompting for confirmation.
If the --image flag is not specified, will prompt with a list.
Specify --force to prevent prompting.

The root password will be the only thing output on stdout - good for scripts!
	    `,
		Action: With(VirtualMachineNameProvider, AuthProvider, func(c *Context) error {
			imageInstall, defaulted, err := prepareImageInstall(nil) // TODO(telyn): this is going to take a bit of a rewrite.
			if err != nil {
				return err
			}
			if defaulted {
				log.Log("No image was specified")
				return &util.PEBKACError{}
			}

			return global.Client.ReimageVirtualMachine(c.VirtualMachineName, imageInstall)
		}),
	})
}

func addImageInstallFlags(flags *flag.FlagSet) {
	flags.String("firstboot-script-file", "", "")
	flags.String("image", "", "")
	flags.String("public-keys", "", "")
	flags.String("public-keys-file", "", "")
	flags.String("root-password", "", "")
}

func prepareImageInstall(flags *flag.FlagSet) (imageInstall *lib.ImageInstall, defaulted bool, err error) {
	if flags == nil || !flags.Parsed() {
		return nil, false, nil
	}
	firstbootScriptFileF := flags.Lookup("firstboot-script-file")
	imageF := flags.Lookup("image")
	publicKeysF := flags.Lookup("public-keys")
	publicKeysFileF := flags.Lookup("public-keys-file")
	rootPasswordF := flags.Lookup("root-password")

	var image, pubkeys, rootpass string

	if imageF != nil && imageF.Value.String() != "" {
		image = imageF.Value.String()
	} else {
		image = "symbiosis"
		defaulted = true
	}
	firstbootScript := ""
	if firstbootScriptFileF != nil && firstbootScriptFileF.Value.String() != "" {
		firstbootScriptContents, err := ioutil.ReadFile(firstbootScriptFileF.Value.String())
		if err != nil {
			return nil, false, err
		}
		firstbootScript = string(firstbootScriptContents)
	}

	if publicKeysF != nil {
		pubkeys = publicKeysF.Value.String()
	}
	if publicKeysFileF != nil && publicKeysFileF.Value.String() != "" {
		pubkeysContents, err := ioutil.ReadFile(publicKeysFileF.Value.String())
		if err != nil {
			return nil, false, err
		}
		pubkeys = pubkeys + "\n" + string(pubkeysContents)
	}

	if rootPasswordF != nil && rootPasswordF.Value.String() != "" {
		rootpass = rootPasswordF.Value.String()
	} else {
		rootpass = util.GeneratePassword()
	}

	return &lib.ImageInstall{
		Distribution:    image,
		FirstbootScript: string(firstbootScript),
		PublicKeys:      pubkeys,
		RootPassword:    rootpass,
	}, defaulted, nil
}
