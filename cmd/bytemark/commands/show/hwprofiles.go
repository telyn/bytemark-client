package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "hwprofiles",
		Usage:     "show hardware profiles available for cloud servers",
		UsageText: "show hwprofiles [--json]",
		Description: `Hardware profiles are used by cloud servers and choosing between them can be thought of as 'which virtual motherboard should I use?'.
Generally bytemark provide two - virtio and compatibility. The virtio one has better performance but may not work with obscure operating systems, or without drivers (particularly important if you are installing windows from CD rather than our images`,
		Flags: app.OutputFlags("hardware profiles", "array"),
		Action: app.Action(with.Definitions, func(c *app.Context) error {
			return c.OutputInDesiredForm(c.Definitions.HardwareProfileDefinitions(), output.List)
		}),
	})
}
