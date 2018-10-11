package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "storage",
		Usage:       "show available storage grades for cloud servers",
		UsageText:   "show storage",
		Description: "This outputs the available storage grades for cloud servers.",
		Flags:       app.OutputFlags("storage grades", "array"),
		Action: app.Action(with.Definitions, func(c *app.Context) error {
			return c.OutputInDesiredForm(c.Definitions.StorageGradeDefinitions(), output.List)
		}),
	})
}
