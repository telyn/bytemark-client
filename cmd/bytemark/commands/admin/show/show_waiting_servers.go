package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "stopped waiting servers",
		Usage:     "shows a list of stopped VMs that should be running",
		UsageText: "--admin show waiting servers [--json]",
		Flags:     app.OutputFlags("servers", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			vms, err := c.Client().GetStoppedEligibleVMs()
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(vms, output.Table)
		}),
	})
}
