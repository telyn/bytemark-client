package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "resize",
		Usage:       "Resize a cloud server's disc",
		UsageText:   "bytemark resize disc <server> <disc label> <size>",
		Description: "Resizes the given server's disc to the given size. Sizes may be specified with a + in front, in which case they are interpreted as relative. For example, '+2GB' is parsed as 'increase the disc size by 2GiB', where '2GB' is parsed as 'set the size of the disc to 2GiB'",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "disc",
			Aliases:     []string{"disk"},
			Usage:       "Resize a cloud server's disc",
			UsageText:   "bytemark resize disc <server> <disc label> <size>",
			Description: "Resizes the given server's disc to the given size. Sizes may be specified with a + in front, in which case they are interpreted as relative. For example, '+2GB' is parsed as 'increase the disc size by 2GiB', where '2GB' is parsed as 'set the size of the disc to 2GiB'",
			Action: With(VirtualMachineNameProvider, DiscLabelProvider, AuthProvider, func(c *Context) (err error) {
				const (
					SET = iota
					INCREASE
				)
				mode := SET
				sizeStr, err := c.NextArg()
				if err != nil {
					log.Error("No size specified")
					return err
				}
				if strings.HasPrefix(sizeStr, "+") {
					sizeStr = sizeStr[1:]
					mode = INCREASE
				}

				size, err := util.ParseSize(sizeStr)
				if err != nil {
					return err
				}

				oldDisc, err := global.Client.GetDisc(c.VirtualMachineName, *c.DiscLabel)
				if err != nil {
					return err
				}

				if mode == INCREASE {
					size = oldDisc.Size + size
				}

				log.Logf("Resizing %s from %dGiB to %dGiB...", oldDisc.Label, oldDisc.Size/1024, size/1024)

				err = global.Client.ResizeDisc(c.VirtualMachineName, *c.DiscLabel, size)
				if err != nil {
					log.Logf("Failed!\r\n")
					return
				}
				log.Logf("Completed.\r\n")
				return
			}),
		}},
	})
}
