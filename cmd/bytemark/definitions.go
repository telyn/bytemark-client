package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "images",
		Aliases:     []string{"distributions", "distros"},
		Usage:       "list images available for installation on all servers",
		UsageText:   "images",
		Description: "This command lists all the images that are available for installation on Bytemark servers.",
		Flags:       app.OutputFlags("images", "array"),
		Action: app.Action(with.Definitions, func(c *app.Context) error {
			return c.OutputInDesiredForm(c.Definitions.DistributionDefinitions(), output.List)
		}),
	}, cli.Command{
		Name:        "zones",
		Usage:       "list available zones for cloud servers",
		UsageText:   "zones",
		Description: "This outputs the zones available for cloud servers to be stored and started in. Note that it is not currently possible to migrate a server between zones.",
		Flags:       app.OutputFlags("zones", "array"),
		Action: app.Action(with.Definitions, func(c *app.Context) error {
			return c.OutputInDesiredForm(c.Definitions.ZoneDefinitions(), output.List)
		}),
	})
}
