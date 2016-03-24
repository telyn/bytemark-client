package main

import (
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
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
	endpoint - the API endpoint to connect to. https://uk0.bigv.io is the default
	billing-endpoint - the billing API endpoint to connect to.
	auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.
	debug-level - the default debug level. Set to 0 unless you like lots of output
	token - the token used for authentication.") // You can get one using bytemark auth.`,
		Subcommands: []cli.Command{
			{
				Name:      "set",
				UsageText: "bytemark config set <variable> <value>",
				Usage:     "Sets a variable by writing to your bytemark config (usually ~/.bytemark)",
				Action: With(func(ctx *Context) error {
					varname, err := ctx.NextArg()
					if err != nil {
						return err
					}
					varname = strings.ToLower(varname)

					oldVar, err := global.Config.GetV(varname)
					if err != nil {
						return err
					}

					// TODO(telyn): consider validating input for the set command
					value, err := ctx.NextArg()
					if err != nil {
						return err
					}

					err = global.Config.SetPersistent(varname, value, "CMD set")
					if err != nil {
						return err
						// TODO(telyn): wrap the error in an error of my own to expose this bhvr
						/*
							if e, ok := err.(*util.ConfigReadError); ok {
								log.Errorf("Couldn't set %s - %v\r\n", e.Name, e.Err)
							} else {
								log.Errorf("Couldn't set %s - %v\r\n", variable, err)
							}
						*/
					}

					if oldVar.Source == "config" {
						log.Logf("%s has been changed.\r\nOld value: %s\r\nNew value: %s\r\n", varname, oldVar.Value, global.Config.GetIgnoreErr(varname))
					} else {
						log.Logf("%s has been set. \r\nNew value: %s\r\n", varname, global.Config.GetIgnoreErr(varname))
					}
					return nil
				}),
			}, {
				Name:      "unset",
				UsageText: "bytemark config unset <variable>",
				Usage:     "Unsets a variable by removing data from bytemark config (usually ~/.bytemark)",
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
		Action: func(ctx *cli.Context) {
			vars, err := global.Config.GetAll()
			if err != nil {
				global.Error = err
				return
			}
			for _, v := range vars {
				log.Logf("%s\t: '%s' (%s)\r\n", v.Name, v.Value, v.Source)
			}
		},
	})
}
