package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
)

// ProviderFunc is the function type that can be passed to With()
type ProviderFunc func(*Context) error

// With is a convenience function for making cli.Command.Actions that sets up a Context, runs all the providers, cleans up afterward and returns errors from the actions if there is one
func With(providers ...ProviderFunc) func(c *cli.Context) error {
	return func(cliContext *cli.Context) error {
		c := Context{Context: cliContextWrapper{cliContext}}
		err := foldProviders(&c, providers...)
		cleanup(&c)
		return err
	}
}

// cleanup resets the value of special flags between invocations of global.App.Run so that the tests pass.
// This is needed because the init() functions are only executed once during the testing cycle.
// Outside of the tests, global.App.Run is only called once before the program closes.
func cleanup(c *Context) {
	ips, ok := c.Context.Generic("ip").(*util.IPFlag)
	if ok {
		*ips = make([]net.IP, 0)
	}
	disc, ok := c.Context.Generic("disc").(*util.DiscSpecFlag)
	if ok {
		*disc = make([]brain.Disc, 0)
	}
	size, ok := c.Context.Generic("memory").(*util.SizeSpecFlag)
	if ok {
		*size = 0
	}
	server, ok := c.Context.Generic("server").(*VirtualMachineNameFlag)
	if ok {
		*server = VirtualMachineNameFlag{}
	}
	server, ok = c.Context.Generic("from").(*VirtualMachineNameFlag)
	if ok {
		*server = VirtualMachineNameFlag{}
	}
	server, ok = c.Context.Generic("to").(*VirtualMachineNameFlag)
	if ok {
		*server = VirtualMachineNameFlag{}
	}
	group, ok := c.Context.Generic("group").(*GroupNameFlag)
	if ok {
		*group = GroupNameFlag{}
	}
}

// foldProviders runs all the providers with the given context, stopping if there's an error
func foldProviders(c *Context, providers ...ProviderFunc) (err error) {
	for _, provider := range providers {
		err = provider(c)
		if err != nil {
			return
		}
	}
	return
}

// OptionalArgs takes a list of flag names. For each flag name it attempts to read the next arg and set the flag with the corresponding name.
// for instance:
// OptionalArgs("server", "disc", "size")
// will attempt to read 3 arguments, setting the "server" flag to the first, "disc" to the 2nd, "size" to the third.
func OptionalArgs(args ...string) ProviderFunc {
	return func(c *Context) error {
		for _, name := range args {
			value, err := c.NextArg()
			if err != nil {
				// if c.NextArg errors that means there aren't more arguments
				// so we just return nil - returning an error would stop the execution of the action.
				return nil
			}
			err = c.Context.Set(name, value)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// JoinArgs is like OptionalArgs, but reads up to n arguments joined with spaces and sets the one named flag.
// if n is not set, reads all the remaining arguments.
func JoinArgs(flagName string, n ...int) ProviderFunc {
	return func(c *Context) (err error) {
		toRead := len(c.Args())
		if len(n) > 0 {
			toRead = n[0]
		}

		value := make([]string, 0, toRead)
		for i := 0; i < toRead; i++ {
			arg, argErr := c.NextArg()
			if argErr != nil {
				// don't return the error - just means we ran out of arguments to slurp
				break
			}
			value = append(value, arg)
		}
		err = c.Context.Set(flagName, strings.Join(value, " "))
		return

	}
}

func isIn(needle string, haystack []string) bool {
	for _, str := range haystack {
		if needle == str {
			return true
		}
	}
	return false
}

func flagValueIsOK(c *Context, flag cli.Flag) bool {
	switch realFlag := flag.(type) {
	case cli.GenericFlag:
		switch value := realFlag.Value.(type) {
		case *VirtualMachineNameFlag:
			return value.VirtualMachine != ""
		case *GroupNameFlag:
			return value.Group != ""
		case *AccountNameFlag:
			return *value != ""
		case *util.SizeSpecFlag:
			return *value != 0
		case *PrivilegeFlag:
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
func RequiredFlags(flagNames ...string) ProviderFunc {
	return func(c *Context) (err error) {
		for _, flag := range c.Command().Flags {
			if isIn(flag.GetName(), flagNames) && !flagValueIsOK(c, flag) {
				return fmt.Errorf("--%s not set (or should not be blank/zero)", flag.GetName())
			}
		}
		return nil
	}
}

// AccountProvider gets an account name from a flag, then the account details from the API, then stitches it to the context
func AccountProvider(flagName string) ProviderFunc {
	return func(c *Context) (err error) {
		err = AuthProvider(c)
		if err != nil {
			return
		}
		accName := c.String(flagName)
		if accName == "" {
			accName = c.Config().GetIgnoreErr("account")
		}

		acc, err := c.Client().GetAccount(accName)
		if err != nil {
			return
		}
		c.Account = &acc
		return
	}
}

// AuthProvider makes sure authentication has been successfully completed, attempting it if necessary.
func AuthProvider(c *Context) (err error) {
	if !c.Authed {
		err = EnsureAuth(c.Client(), c.Config())
		if err != nil {
			return
		}
	}
	c.Authed = true
	return
}

// DiscProvider gets a VirtualMachineName from a flag and a disc from another, then gets the named Disc from the brain and attaches it to the Context.
func DiscProvider(vmFlagName, discFlagName string) ProviderFunc {
	return func(c *Context) (err error) {
		if c.Group != nil {
			return
		}
		err = AuthProvider(c)
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

// DefinitionsProvider gets the Definitions from the brain and attaches them to the Context.
func DefinitionsProvider(c *Context) (err error) {
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

// GroupProvider gets a GroupName from a flag, then gets the named Group from the brain and attaches it to the Context.
func GroupProvider(flagName string) ProviderFunc {
	return func(c *Context) (err error) {
		if c.Group != nil {
			return
		}
		err = AuthProvider(c)
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

// PrivilegeProvider gets the named PrivilegeFlag from the context, then resolves its target to an ID if needed to create a brain.Privilege, then attaches that to the context
func PrivilegeProvider(flagName string) ProviderFunc {
	return func(c *Context) (err error) {
		pf := c.PrivilegeFlag(flagName)
		level, ok := normalisePrivilegeLevel(pf.Level)
		if !ok && !c.Bool("force") {
			return fmt.Errorf("Unexpected privilege level '%s' - expecting account_admin, group_admin, vm_admin or vm_console", pf.Level)
		}
		c.Privilege = brain.Privilege{
			Username: pf.Username,
			Level:    level,
		}
		err = AuthProvider(c)
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

// UserProvider gets a username from the given flag, then gets the corresponding User from the brain, and attaches it to the Context.
func UserProvider(flagName string) ProviderFunc {
	return func(c *Context) (err error) {
		if c.User != nil {
			return
		}
		if err = AuthProvider(c); err != nil {
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

// VirtualMachineProvider gets a VirtualMachineName from a flag, then gets the named VirtualMachine from the brain and attaches it to the Context.
func VirtualMachineProvider(flagName string) ProviderFunc {
	return func(c *Context) (err error) {
		if c.VirtualMachine != nil {
			return
		}
		err = AuthProvider(c)
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
