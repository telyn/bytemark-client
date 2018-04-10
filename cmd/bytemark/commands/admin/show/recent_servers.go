package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "recent servers",
		Usage:     "shows a list of recently created servers",
		UsageText: "--admin show recent servers [--json | --table] [--table-fields <fields> | --table-fields help]",
		Flags:     app.OutputFlags("servers", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			vms, err := c.Client().GetRecentVMs()
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(vms, output.Table)
		}),
	})
}
