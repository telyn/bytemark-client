package main

import (
	"errors"
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"net/url"
	"strconv"
	"strings"
)

func validateEndpointForConfig(endpoint string) error {
	url, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return errors.New("The endpoint URL should start with http:// or https://")
	}
	if url.Host == "" {
		return errors.New("The endpoint URL should have a hostname")
	}
	return nil
}

func validateAccountForConfig(c *Context, name string) (err error) {
	_, err = global.Client.GetAccount(name)
	if err != nil {
		if _, ok := err.(lib.NotFoundError); ok {
			return fmt.Errorf("No such account %s - check your typing and specify --yubikey if necessary", name)
		}
		return err
	}
	return
}

func validateGroupForConfig(c *Context, name string) (err error) {
	// we can't just use GroupProvider because it expects NextArg() to be the account name - there's no way to pass one in.
	groupName := global.Client.ParseGroupName(name, global.Config.GetGroup())
	_, err = global.Client.GetGroup(groupName)
	if err != nil {
		if _, ok := err.(lib.NotFoundError); ok {
			return fmt.Errorf("No such group %v - check your typing and specify --yubikey if necessary", groupName)
		}
		return err
	}
	return
}

func validateConfigValue(c *Context, varname string, value string) error {
	if c.Bool("force") {
		return nil
	}
	switch varname {
	case "endpoint", "api-endpoint", "billing-endpoint", "spp-endpoint", "auth-endpoint":
		return validateEndpointForConfig(value)
	case "account":
		return validateAccountForConfig(c, value)
	case "group":
		return validateGroupForConfig(c, value)
	case "debug-level":
		_, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return errors.New("debug-level must be an integer")
		}
	}
	return nil
}

func init() {
	commands = append(commands, cli.Command{
		Name:      "config",
		Usage:     "manage the bytemark client's configuration",
		UsageText: "config [ set | unset ]",
		Description: `view and manipulate the bytemark-client configuration
		
When invoked with no subcommand, outputs the current values of all variables and what source they were derived from.
The set and unset subcommands can be used to set and unset such variables.
		
    Available variables:
        account - the default account, used when you do not explicitly state an account - defaults to the same as your user name
        token - the token used for authentication
        user - the user that you log in as by default
        group - the default group, used when you do not explicitly state a group (defaults to 'default')

        debug-level - the default debug level. Set to 0 unless you like lots of output.
	api-endpoint - the endpoint for domains (among other things?)
        auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.
        endpoint - the brain endpoint to connect to. https://uk0.bigv.io is the default.
        billing-endpoint - the billing API endpoint to connect to. https://bmbilling.bytemark.co.uk is the default.
        spp-endpoint - the SPP endpoint to use. https://spp-submissions.bytemark.co.uk is the default.`,
		Subcommands: []cli.Command{
			{
				Name:        "set",
				UsageText:   "bytemark config set <variable> <value>",
				Usage:       "sets a bytemark client configuration request",
				Description: "Sets the named variable to the given value. See `bytemark help config` for which variables are available",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "force",
						Usage: "Don't run any validation checks against the value",
					},
				},
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

					if varname == "account" || varname == "group" {
						err = AuthProvider(ctx)
						if err != nil {
							return err
						}
					}

					err = validateConfigValue(ctx, varname, value)
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
				Usage:       "unsets a bytemark client configuration option",
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
		Action: func(ctx *cli.Context) (err error) {
			if ctx.Bool("help") {
				err = cli.ShowSubcommandHelp(ctx)
				return
			}
			vars, err := global.Config.GetAll()
			if err != nil {
				return
			}
			for _, v := range vars {
				log.Logf("%s\t: '%s' (%s)\r\n", v.Name, v.Value, v.Source)
			}
			return
		},
	})
}
