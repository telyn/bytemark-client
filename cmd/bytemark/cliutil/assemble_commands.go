package cliutil

import (
	"github.com/urfave/cli"
)

func AssembleCommands(commandSlices ...[]cli.Command) (allCommands []cli.Command) {
	for _, commands := range commandSlices {
		allCommands = MergeCommands(allCommands, commands)
	}
	return
}
