package update

import "github.com/urfave/cli"

// Commands is a slice which is populated during init() with all the available top-level commands
var Commands = make([]cli.Command, 0)
