package flagsets

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/urfave/cli"
)

const (
	force = "force"
)

// Force is a flag common to a bunch of commands and can have a generic Usage.
var Force = cli.BoolFlag{
	Name:  force,
	Usage: "Do not prompt for confirmation when destroying data or increasing costs.",
}

// Forced indicates whether the Force flag was set in this context
func Forced(c *app.Context) bool {
	return c.Bool(force)
}
