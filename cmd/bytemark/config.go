package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
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
				Action:    config_set, // defined below - it's just a bit long
			},
			{
				Name:      "unset",
				UsageText: "bytemark config unset <variable>",
				Usage:     "Unsets a variable by removing data from bytemark config (usually ~/.bytemark)",
				Action: func(ctx *cli.Context) {
					variable := strings.ToLower(ctx.Args()[0])
					err := global.Config.Unset(variable)
					global.Error = err
				},
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

func config_set(ctx *cli.Context) {
	variable := strings.ToLower(ctx.Args().First())

	oldVar, err := global.Config.GetV(variable)
	if err != nil {
		global.Error = err
	}

	if len(ctx.Args()) < 2 {
		global.Error = util.PEBKACError{}
		return
	}

	// TODO(telyn): consider validating input for the set command
	err = global.Config.SetPersistent(variable, ctx.Args()[1], "CMD set")
	if err != nil {
		global.Error = err
		return
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
		log.Logf("%s has been changed.\r\nOld value: %s\r\nNew value: %s\r\n", variable, oldVar.Value, global.Config.GetIgnoreErr(variable))
	} else {
		log.Logf("%s has been set. \r\nNew value: %s\r\n", variable, global.Config.GetIgnoreErr(variable))
	}
}
