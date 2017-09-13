package with

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
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
		case *app.VirtualMachineNameFlag:
			return value.VirtualMachineName != nil
		case *app.GroupNameFlag:
			return value.GroupName != nil
		case *app.AccountNameFlag:
			return value.AccountName != ""
		case *util.SizeSpecFlag:
			return *value != 0
		case *app.PrivilegeFlag:
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

// Account gets an account name from a flag, then the account details from the API, then stitches it to the context
func Account(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		err = preflight(c)
		if err != nil {
			return
		}
		accName := c.String(flagName)
		c.Debug("flagName: %s accName: %s\n", flagName, accName)
		if accName == "" {
			accName = c.Config().GetIgnoreErr("account")
		}
		c.Debug("flagName: %s a4tName: %s\n", flagName, accName)

		acc, err := c.Client().GetAccount(accName)
		if err != nil {
			return
		}
		c.Account = &acc
		return
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

		vmName := c.VirtualMachineName(vmFlagName)
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

		groupName := c.GroupName(flagName)
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

// normalisePrivilegeLevel makes sure the level provided is actually a valid PrivilegeLevel and provides a couple of aliases.
func normalisePrivilegeLevel(l brain.PrivilegeLevel) (level brain.PrivilegeLevel, ok bool) {
	level = brain.PrivilegeLevel(strings.ToLower(string(l)))
	switch level {
	case "cluster_admin", "account_admin", "group_admin", "vm_admin", "vm_console":
		ok = true
	case "server_admin", "server_console":
		level = brain.PrivilegeLevel(strings.Replace(string(level), "server", "vm", 1))
		ok = true
	case "console":
		level = "vm_console"
		ok = true
	}
	return
}

// Privilege gets the named PrivilegeFlag from the context, then resolves its target to an ID if needed to create a brain.Privilege, then attaches that to the context
func Privilege(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		pf := c.PrivilegeFlag(flagName)
		level, ok := normalisePrivilegeLevel(pf.Level)
		if !ok && !c.Bool("force") {
			return fmt.Errorf("Unexpected privilege level '%s' - expecting account_admin, group_admin, vm_admin or vm_console", pf.Level)
		}
		c.Privilege = brain.Privilege{
			Username: pf.Username,
			Level:    level,
		}
		err = preflight(c)
		if err != nil {
			return
		}
		switch c.Privilege.TargetType() {
		case brain.PrivilegeTargetTypeVM:
			var vm brain.VirtualMachine
			vm, err = c.Client().GetVirtualMachine(*pf.VirtualMachineName)
			c.Privilege.VirtualMachineID = vm.ID
		case brain.PrivilegeTargetTypeGroup:
			var group brain.Group
			group, err = c.Client().GetGroup(*pf.GroupName)
			c.Privilege.GroupID = group.ID
		case brain.PrivilegeTargetTypeAccount:
			var acc lib.Account
			acc, err = c.Client().GetAccount(pf.AccountName)
			c.Privilege.AccountID = acc.BrainID
		}
		return
	}
}

// User gets a username from the given flag, then gets the corresponding User from the brain, and attaches it to the app.Context.
func User(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		if c.User != nil {
			return
		}
		if err = preflight(c); err != nil {
			return
		}
		username := c.String(flagName)
		if username == "" {
			username = c.Client().GetSessionUser()
		}

		user, err := c.Client().GetUser(username)
		if err != nil {
			return
		}
		c.User = &user
		return
	}
}

// VirtualMachine gets a VirtualMachineName from a flag, then gets the named VirtualMachine from the brain and attaches it to the app.Context.
func VirtualMachine(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		if c.VirtualMachine != nil {
			return
		}
		err = preflight(c)
		if err != nil {
			return
		}
		vmName := c.VirtualMachineName(flagName)
		vm, err := c.Client().GetVirtualMachine(vmName)
		if err != nil {
			return
		}
		c.VirtualMachine = &vm
		return
	}
}
