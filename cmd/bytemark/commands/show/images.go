package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "images",
		Aliases:     []string{"distributions", "distros"},
		Usage:       "show images available for installation on all servers",
		UsageText:   "show images",
		Description: "This command lists all the images that are available for installation on Bytemark servers.",
		Flags:       app.OutputFlags("images", "array"),
		Action: app.Action(with.Definitions, func(c *app.Context) error {
			return c.OutputInDesiredForm(c.Definitions.DistributionDefinitions(), output.List)
		}),
	})
}
