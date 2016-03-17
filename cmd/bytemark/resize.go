package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
	"strings"
)

func (cmds *CommandSet) HelpForResize() util.ExitCode {
	log.Log("bytemark resize")
	log.Log("")
	log.Log("usage: bytemark resize disc <server> <disc> <size>")
	return util.E_USAGE_DISPLAYED
}

func init() {
	commands = append(commands, cli.Command{
		Name: "resize",
		Subcommands: []cli.Command{{
			Name:    "disc",
			Aliases: []string{"disk"},
			Action: With(VirtualMachineProvider, DiscLabelProvider, func(c *Context) (err error) {
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
