package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "config",
		Usage:     "show the bytemark client's configuration",
		UsageText: "show config",
		Description: `View the bytemark-client configuration.

    The following variables are displayed:` + config.VarsDescription,
		Flags:  app.OutputFlags("vars", "array"),
		Action: app.Action(viewConfig),
	})
}

func viewConfig(c *app.Context) error {
	vars, err := c.Config().GetAll()
	if err != nil {
		return err
	}
	return c.OutputInDesiredForm(vars, output.List)
}
