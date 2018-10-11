package update

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

type configVariable struct {
	name        string
	configName  string
	description string
	validate    func(*app.Context, string) error
	needAuth    bool
}

type configVars []configVariable

var configVariables = configVars{
	{
		name:        "endpoint",
		description: "brain endpoint to connect to",
		validate:    validateEndpointForConfigFunc("endpoint"),
	},
	{
		name:        "api-endpoint",
		description: "endpoint for domains",
		validate:    validateEndpointForConfigFunc("api-endpoint"),
	},
	{
		name:        "billing-endpoint",
		description: "billing API endpoint to connect to",
		validate:    validateEndpointForConfigFunc("billing-endpoint"),
	},
	{
		name:        "spp-endpoint",
		description: "SPP endpoint to use",
		validate:    validateEndpointForConfigFunc("spp-endpoint"),
	},
	{
		name:        "auth-endpoint",
		description: "endpoint to authenticate to",
		validate:    validateEndpointForConfigFunc("auth-endpoint"),
	},
	{
		name:        "default-debug-level",
		configName:  "debug_level",
		description: "default debug level",
		validate:    validateIntForConfigFunc("default-debug-level"),
	},
	{
		name:        "token",
		description: "token used for authentication",
	},
	{
		name:        "user",
		description: "user that you log in as by default",
	},
	{
		name:        "account",
		description: "default account",
		validate:    validateAccountForConfig,
		needAuth:    true,
	},
	{
		name:        "group",
		description: "default group",
		validate:    validateGroupForConfig,
		needAuth:    true,
	},
}

func (variable configVariable) confName() string {
	if variable.configName != "" {
		return variable.configName
	}
	return variable.name
}

func (variable configVariable) getFlags(c *app.Context) (string, bool) {
	return c.String(variable.name), c.Bool("unset-" + variable.name)
}

func (variable configVariable) present(c *app.Context) bool {
	set, unset := variable.getFlags(c)
	return set != "" || unset
}

func (variables configVars) present(c *app.Context) (out configVars) {
	out = configVars{}
	for _, variable := range variables {
		if variable.present(c) {
			out = append(out, variable)
		}
	}
	return
}

func (variables configVars) configFlags() (flags []cli.Flag) {
	flags = make([]cli.Flag, len(variables)*2)
	for i, variable := range variables {
		flags[i*2] = cli.StringFlag{
			Name:  variable.name,
			Usage: "Sets the " + variable.description,
		}
		flags[i*2+1] = cli.BoolFlag{
			Name:  "unset-" + variable.name,
			Usage: "Unsets the " + variable.description,
		}
	}
	return
}

func validateAccountForConfig(c *app.Context, name string) (err error) {
	_, err = c.Client().GetAccount(name)
	if err != nil {
		switch err.(type) {
		case lib.NotFoundError:
			return fmt.Errorf("No such account %s - check your typing and specify --yubikey if necessary", name)
		case lib.BillingAccountNotFound:
			return nil
		}
	}
	return
}

func validateGroupForConfig(c *app.Context, name string) (err error) {
	groupName := lib.ParseGroupName(name, c.Config().GetGroup())
	_, err = c.Client().GetGroup(groupName)
	if err != nil {
		if _, ok := err.(lib.NotFoundError); ok {
			return fmt.Errorf("No such group %v - check your typing and specify --yubikey if necessary", groupName)
		}
		return err
	}
	return
}

func validateEndpointForConfigFunc(variable string) func(*app.Context, string) error {
	return func(c *app.Context, endpoint string) error {
		url, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		if url.Scheme != "http" && url.Scheme != "https" {
			return errors.New(variable + " URL should start with http:// or https://")
		}
		if url.Host == "" {
			return errors.New(variable + " URL should have a hostname")
		}
		return nil
	}
}

func validateIntForConfigFunc(variable string) func(*app.Context, string) error {
	return func(c *app.Context, value string) error {
		_, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return errors.New(variable + " must be an integer")
		}
		return nil
	}
}

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "config",
		Usage:     "update the bytemark client's configuration",
		UsageText: "update config [flags]",
		Description: `Manipulate the bytemark-client configuration

    Available variables:` + config.VarsDescription,
		Flags:  append(configVariables.configFlags(), flags.Force),
		Action: app.Action(configVariables.updateConfig),
	})
}

func (variables configVars) updateConfig(c *app.Context) error {
	presentVariables := variables.present(c)
	withAuth := false
	// first pass, validate
	if len(presentVariables) == 0 {
		return c.Help("missing arguments")
	}
	for _, variable := range presentVariables {
		set, unset := variable.getFlags(c)
		if set != "" && unset {
			return c.Help("cannot set and unset " + variable.name)
		}
		if set != "" && !flags.Forced(c) && variable.validate != nil {
			if variable.needAuth && !withAuth {
				err := with.Auth(c)
				if err != nil {
					return err
				}
				withAuth = true
			}
			err := variable.validate(c, set)
			if err != nil {
				return err
			}
		}
	}
	// second pass, apply
	for _, variable := range presentVariables {
		set, unset := variable.getFlags(c)
		if unset {
			err := c.Config().Unset(variable.confName())
			if err != nil {
				return err
			}
			log.Logf("%s has been unset. \r\n", variable.name)
		} else {
			oldVar, err := c.Config().GetV(variable.confName())
			if err != nil {
				return err
			}
			err = c.Config().SetPersistent(variable.confName(), set, "CMD set")
			if err != nil {
				return err
			}

			if oldVar.Source == "config" {
				log.Logf("%s has been changed.\r\nOld value: %s\r\nNew value: %s\r\n", variable.confName(), oldVar.Value, c.Config().GetIgnoreErr(variable.confName()))
			} else {
				log.Logf("%s has been set. \r\nNew value: %s\r\n", variable.confName(), c.Config().GetIgnoreErr(variable.confName()))
			}
		}
	}
	return nil
}
