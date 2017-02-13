package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
	"net"
	"strings"
)

// ProviderFunc is the function type that can be passed to With()
type ProviderFunc func(*Context) error

// With is a convenience function for making cli.Command.Actions that sets up a Context, runs all the providers, cleans up afterward and returns errors from the actions if there is one
func With(providers ...ProviderFunc) func(c *cli.Context) error {
	return func(cliContext *cli.Context) error {
		c := Context{Context: cliContext}
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
			arg, err := c.NextArg()
			if err != nil {
				// don't return the error - just means we ran out of arguments to slurp
				break
			}
			value = append(value, arg)
		}
		err = c.Context.Set(flagName, strings.Join(value, " "))
		return

	}
}

// AccountProvider gets an account name from a flag, then the account details from the API, then stitches it to the context
func AccountProvider(flagName string) ProviderFunc {
	return func(c *Context) (err error) {
		c.Account, err = global.Client.GetAccount(c.String(flagName))
		return
	}
}

// AuthProvider makes sure authentication has been successfully completed, attempting it if necessary.
func AuthProvider(c *Context) (err error) {
	if !c.Authed {
		err = EnsureAuth()
		if err != nil {
			return
		}
	}
	c.Authed = true
	return
}

// DefinitionsProvider gets the Definitions from the brain and attaches them to the Context.
func DefinitionsProvider(c *Context) (err error) {
	if c.Definitions != nil {
		return
	}
	c.Definitions, err = global.Client.ReadDefinitions()
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
		c.Group, err = global.Client.GetGroup(&groupName)
		return
	}
}

// UserProvider calls UserNameProvider, gets the User from the brain, and attaches it to the Context.
func UserProvider(flagName string) ProviderFunc {
	return func(c *Context) (err error) {
		if c.User != nil {
			return
		}
		user := c.String(flagName)
		if user != "" {
			user = global.Config.GetIgnoreErr("user")
		}
		if err = AuthProvider(c); err != nil {
			return
		}
		c.User, err = global.Client.GetUser(user)
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
		c.VirtualMachine, err = global.Client.GetVirtualMachine(&vmName)
		return
	}
}
