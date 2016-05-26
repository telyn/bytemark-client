package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "config",
		Usage:     "Manage the bytemark client's configuration",
		UsageText: "Outputs the current values of all variables and what source they were derived from.",
		Description: `When invoked with no subcommand, outputs the current values of all variables and what source they were derived from.
The set and unset subcommands can be used to set and unset such variables.
		
Available variables:
	account - the default account, used when you do not explicitly state an account - defaults to the same as your user name
	token - the token used for authentication
	user - the user that you log in as by default
	group - the default group, used when you do not explicitly state a group (defaults to 'default')

	debug-level - the default debug level. Set to 0 unless you like lots of output.
	auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.
	endpoint - the brain endpoint to connect to. https://uk0.bigv.io is the default.
	billing-endpoint - the billing API endpoint to connect to. https://bmbilling.bytemark.co.uk is the default.
	spp-endpoint - the SPP endpoint to use. https://spp-submissions.bytemark.co.uk is the default.`,
		Subcommands: []cli.Command{
			{
				Name:        "set",
				UsageText:   "bytemark config set <variable> <value>",
				Usage:       "Sets a variable by writing to your bytemark config (usually ~/.bytemark)",
				Description: "Sets the named variable to the given value. See `bytemark help config` for which variables are available",
				Action: With(func(ctx *Context) error {
					varname, err := ctx.NextArg()
					if err != nil {
						return err
					}
					varname = strings.ToLower(varname)

					if !util.IsConfigVar(varname) {
						return ctx.Help(fmt.Sprintf("%s is not a valid variable name", varname))
					}

					oldVar, err := global.Config.GetV(varname)
					if err != nil {
						return err
					}

					value, err := ctx.NextArg()
					if err != nil {
						return err
					}

					err = global.Config.SetPersistent(varname, value, "CMD set")
					if err != nil {
						return err
					}

					if oldVar.Source == "config" {
						log.Logf("%s has been changed.\r\nOld value: %s\r\nNew value: %s\r\n", varname, oldVar.Value, global.Config.GetIgnoreErr(varname))
					} else {
						log.Logf("%s has been set. \r\nNew value: %s\r\n", varname, global.Config.GetIgnoreErr(varname))
					}
					return nil
				}),
			}, {
				Name:        "unset",
				UsageText:   "bytemark config unset <variable>",
				Usage:       "Unsets a variable by removing data from bytemark config (usually ~/.bytemark)",
				Description: "Unsets the named variable.",
				Action: With(func(ctx *Context) error {
					varname, err := ctx.NextArg()
					if err != nil {
						return err
					}
					varname = strings.ToLower(varname)
					return global.Config.Unset(varname)
				}),
			},
		},
		Action: func(ctx *cli.Context) error {
			vars, err := global.Config.GetAll()
			if err != nil {
				return err
			}
			for _, v := range vars {
				log.Logf("%s\t: '%s' (%s)\r\n", v.Name, v.Value, v.Source)
			}
			return nil
		},
	})
}
