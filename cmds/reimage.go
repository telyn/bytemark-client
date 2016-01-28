package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"flag"
	"io/ioutil"
)

func (cmds *CommandSet) HelpForReimage() util.ExitCode {
	log.Log("usage: bytemark reimage (--image <image>) [flags] <server>")
	log.Log("")
	log.Log("flags available")
	log.Log("    --firstboot-script-file <file name> - a script that will be run on first boot")
	log.Log("    --force - disables the confirmation prompt")
	log.Log("    --image <image name> - specify what to image the server with. Default is 'symbiosis'. See `bytemark images` for a list of available images")
	log.Log("    --public-keys <keys> (newline seperated)")
	log.Log("    --public-keys-file <file> (will be read & appended to --public-keys)")
	log.Log("    --root-password <password> (if not set, will be randomly generated)")
	log.Log()
	log.Log("Image the given server with the specified image, prompting for confirmation.")
	log.Log("If the --image flag is not specified, will prompt with a list.")
	log.Log("Specify --force to prevent prompting.")
	log.Log("")
	log.Log("The root password will be the only thing output on stdout - good for scripts!")

	return util.E_USAGE_DISPLAYED
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

func (cmds *CommandSet) Reimage(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	addImageInstallFlags(flags)
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForReimage()
		return util.E_PEBKAC
	}

	name, err := cmds.bigv.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {
		return util.ProcessError(err)
	}
	imageInstall, defaulted, err := prepareImageInstall(flags)
	if err != nil {
		return util.ProcessError(err)
	}
	if defaulted {
		log.Log("No image was specified")

		cmds.HelpForReimage()
		return util.E_PEBKAC
	}

	cmds.EnsureAuth()

	err = cmds.bigv.ReimageVirtualMachine(name, imageInstall)
	return util.ProcessError(err)
}
