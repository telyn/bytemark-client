package util

import (
	"github.com/urfave/cli"
)

// forceFlag is common to a bunch of commands and can have a generic Usage.
var ForceFlag = cli.BoolFlag{
	Name:  "force",
	Usage: "Do not prompt for confirmation when destroying data or increasing costs.",
}
