package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "migrating servers",
		Usage:     "shows a list of migrating servers",
		UsageText: "--admin show migrating servers [--json]",
		Flags:     app.OutputFlags("migrating servers", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			vms, err := c.Client().GetMigratingVMs()
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(vms, output.Table)
		}),
	})
}
