package cliutil

import (
	"github.com/urfave/cli"
)

// AssembleCommands slices of different command sets, merges and returns them as one full list
func AssembleCommands(commandSlices ...[]cli.Command) (allCommands []cli.Command) {
	for _, commands := range commandSlices {
		allCommands = MergeCommands(allCommands, commands)
	}
	return
}
