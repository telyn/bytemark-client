package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util/sizespec"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "resize",
		Usage:     "resize a cloud server's disc",
		UsageText: "bytemark resize disc <server> <disc label> <size>",
		Description: `resize a cloud server's disc

Resizes the given disc to the given size. Sizes may be specified with a + in front, in which case they are interpreted as relative. For example, '+2GB' is parsed as 'increase the disc size by 2GiB', where '2GB' is parsed as 'set the size of the disc to 2GiB'`,
		Action: cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{{
			Name:        "disc",
			Aliases:     []string{"disk"},
			Usage:       "resize a cloud server's disc",
			UsageText:   "bytemark resize disc <server> <disc label> <size>",
			Description: "Resizes the given server's disc to the given size. Sizes may be specified with a + in front, in which case they are interpreted as relative. For example, '+2GB' is parsed as 'increase the disc size by 2GiB', where '2GB' is parsed as 'set the size of the disc to 2GiB'",
			Flags: []cli.Flag{
				forceFlag,
				cli.StringFlag{
					Name:  "disc",
					Usage: "the disc to resize",
				},
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server that the disc is attached to",
					Value: new(VirtualMachineNameFlag),
				},
			},
			Action: With(OptionalArgs("server", "disc"), AuthProvider, func(c *Context) (err error) {
				// TODO(telyn): replace all this with a flag
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

				size, err := sizespec.Parse(sizeStr)
				if err != nil {
					return err
				}

				vmName := c.VirtualMachineName("server")

				oldDisc, err := global.Client.GetDisc(&vmName, c.String("disc"))
				if err != nil {
					return err
				}

				if mode == INCREASE {
					size = oldDisc.Size + size
				}
				log.Logf("Resizing %s from %dGiB to %dGiB...", oldDisc.Label, oldDisc.Size/1024, size/1024)

				if !c.Bool("force") && !util.PromptYesNo(fmt.Sprintf("Are you certain you wish to perform this resize?")) {
					return util.UserRequestedExit{}
				}

				err = global.Client.ResizeDisc(&vmName, c.String("disc"), size)
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
