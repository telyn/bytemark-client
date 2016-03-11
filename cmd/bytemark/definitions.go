package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
)

func (cmds *CommandSet) readDefinitions(args []string) (*lib.Definitions, error) {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	defs, err := cmds.client.ReadDefinitions()
	return defs, err

}

func (cmds *CommandSet) Distributions(args []string) util.ExitCode {
	defs, err := cmds.readDefinitions(args)
	if err != nil {
		return util.ProcessError(err)
	}
	for distro, description := range defs.DistributionDescriptions {
		log.Logf("'%s': %s\r\n", distro, description)
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) HardwareProfiles(args []string) util.ExitCode {
	defs, err := cmds.readDefinitions(args)
	if err != nil {
		return util.ProcessError(err)
	}
	for _, profile := range defs.HardwareProfiles {
		log.Log(profile)
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) StorageGrades(args []string) util.ExitCode {
	defs, err := cmds.readDefinitions(args)
	if err != nil {
		return util.ProcessError(err)
	}

	for grade, description := range defs.StorageGradeDescriptions {
		log.Logf("'%s': %s\r\n", grade, description)
	}
	return util.E_SUCCESS
}

func (cmds *CommandSet) Zones(args []string) util.ExitCode {
	defs, err := cmds.readDefinitions(args)
	if err != nil {
		return util.ProcessError(err)
	}
	for _, zone := range defs.ZoneNames {
		log.Log(zone)
	}
	return util.E_SUCCESS
}
