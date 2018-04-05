package add

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "ip range",
		Usage:     "add a new IP range in a VLAN",
		UsageText: "--admin add ip range <ip-range> <vlan-num>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "ip-range",
				Usage: "the IP range to add",
			},
			cli.IntFlag{
				Name:  "vlan-num",
				Usage: "The VLAN number to add the IP range to",
			},
		},
		Action: app.Action(args.Optional("ip-range", "vlan-num"), with.RequiredFlags("ip-range", "vlan-num"), with.Auth, func(c *app.Context) error {
			if err := c.Client().CreateIPRange(c.String("ip-range"), c.Int("vlan-num")); err != nil {
				return err
			}
			log.Logf("IP range added\r\n")
			return nil
		}),
	})
}
