package admin

import "github.com/urfave/cli"

// Commands is a slice which is populated during init() with all the available admin commands
var Commands = make([]cli.Command, 0)
