package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
	"strings"
)

func (cmds *CommandSet) HelpForResize() util.ExitCode {
	log.Log("bytemark resize")
	log.Log("")
	log.Log("usage: bytemark resize disc <virtual machine> <disc> <size>")
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) ResizeDisc(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	const (
		SET = iota
		INCREASE
	)
	mode := SET

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		log.Error("No virtual machine specified")
		cmds.HelpForResize()
		return util.E_PEBKAC
	}

	name, err := cmds.client.ParseVirtualMachineName(nameStr, cmds.config.GetVirtualMachine())
	if err != nil {

	}

	discId, ok := util.ShiftArgument(&args, "disc id")
	if !ok {
		log.Error("No disc specified")
		cmds.HelpForResize()
		return util.E_PEBKAC
	}

	sizeStr, ok := util.ShiftArgument(&args, "")
	if !ok {
		log.Error("No size specified")
		cmds.HelpForResize()
		return util.E_PEBKAC
	}

	if strings.HasPrefix(sizeStr, "+") {
		sizeStr = sizeStr[1:]
		mode = INCREASE
	}

	size, err := util.ParseSize(sizeStr)
	if err != nil {
		return util.ProcessError(err)
	}

	err = cmds.EnsureAuth()

	if err != nil {
		return util.ProcessError(err)
	}

	oldDisc, err := cmds.client.GetDisc(name, discId)
	if err != nil {
		return util.ProcessError(err)
	}

	if mode == INCREASE {
		size = oldDisc.Size + size
	}

	log.Logf("Resizing %s from %dGiB to %dGiB...", oldDisc.Label, oldDisc.Size/1024, size/1024)

	err = cmds.client.ResizeDisc(name, discId, size)
	if err != nil {
		log.Logf("Failed!\r\n")
		return util.ProcessError(err)
	} else {
		log.Logf("Completed.\r\n")
		return util.E_SUCCESS
	}
}
