package add

import "github.com/urfave/cli"

// Commands is a slice which is populated during init() with all the available top-level commands
var Commands = make([]cli.Command, 0)

// forceFlag is common to a bunch of commands and can have a generic Usage.
var forceFlag = cli.BoolFlag{
	Name:  "force",
	Usage: "Do not prompt for confirmation when destroying data or increasing costs.",
}
