package with

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/urfave/cli"
)

// preflight runs Auth and Preprocess (if needed).
// it's useful because every other *Provider needs to make sure these are run (if needed)
func preflight(c *app.Context) (err error) {
	err = Auth(c)
	if err != nil {
		return
	}
	err = c.Preprocess()
	return
}

func isIn(needle string, haystack []string) bool {
	for _, str := range haystack {
		if needle == str {
			return true
		}
	}
	return false
}

func flagValueIsOK(c *app.Context, flag cli.Flag) bool {
	switch realFlag := flag.(type) {
	case cli.GenericFlag:
		switch value := realFlag.Value.(type) {
		case *flags.VirtualMachineNameFlag:
			return value.VirtualMachineName.VirtualMachine != ""
		case *flags.GroupNameFlag:
			return value.GroupName.Group != ""
		case *flags.AccountNameFlag:
			return value.AccountName != ""
		case *flags.SizeSpecFlag:
			return *value != 0
		case *flags.PrivilegeFlag:
			return value.Username != "" && value.Level != ""
		}
	case cli.StringFlag:
		return c.String(realFlag.Name) != ""
	case cli.IntFlag:
		return c.Int(realFlag.Name) != 0
	}
	return true
}

// RequiredFlags makes sure that the named flags are not their zero-values.
// (or that VirtualMachineName / GroupName flags have the full complement of values needed)
func RequiredFlags(flagNames ...string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		err = c.Preprocess()
		if err != nil {
			return
		}
		for _, flag := range c.Command().Flags {
			if isIn(flag.GetName(), flagNames) && !flagValueIsOK(c, flag) {
				return fmt.Errorf("--%s not set (or should not be blank/zero)", flag.GetName())
			}
		}
		return nil
	}
}

// Auth makes sure authentication has been successfully completed, attempting it if necessary.
func Auth(c *app.Context) (err error) {
	if !c.Authed {
		err = EnsureAuth(c.Client(), c.Config())
		if err != nil {
			return
		}
	}
	c.Authed = true
	return
}

// Disc gets a VirtualMachineName from a flag and a disc from another, then gets the named Disc from the brain and attaches it to the app.Context.
func Disc(vmFlagName, discFlagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		if c.Group != nil {
			return
		}
		err = preflight(c)
		if err != nil {
			return
		}

		vmName := flags.VirtualMachineName(c, vmFlagName)
		discLabel := c.String(discFlagName)
		disc, err := c.Client().GetDisc(vmName, discLabel)
		if err != nil {
			return
		}
		c.Disc = &disc
		return
	}
}

// Definitions gets the Definitions from the brain and attaches them to the app.Context.
func Definitions(c *app.Context) (err error) {
	if c.Definitions != nil {
		return
	}
	defs, err := c.Client().ReadDefinitions()
	if err != nil {
		return
	}
	c.Definitions = &defs
	return
}

// Group gets a GroupName from a flag, then gets the named Group from the brain and attaches it to the app.Context.
func Group(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		if c.Group != nil {
			return
		}
		err = preflight(c)
		if err != nil {
			return
		}

		groupName := flags.GroupName(c, flagName)
		if groupName.Account == "" {
			groupName.Account = c.Config().GetIgnoreErr("account")
		}
		group, err := c.Client().GetGroup(groupName)
		// this if is a guard against tricky-to-debug nil-pointer errors
		if err != nil {
			return
		}
		c.Group = &group
		return
	}
}
