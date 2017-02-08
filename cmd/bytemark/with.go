package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"net"
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
				return nil
			}
			return c.Context.Set(name, value)
		}
		return nil
	}
}

// AccountProvider uses AccountNameProvider and then gets the named account details from the API, then stitches it to the context
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

// GroupNameProvider reads the NextArg, parses it as a GroupName and attaches it to the Context
func GroupNameProvider(c *Context) (err error) {
	if c.GroupName != nil {
		return
	}

	if err = AuthProvider(c); err != nil {
		return
	}

	name, err := c.NextArg()
	if err != nil {
		return err
	}
	c.GroupName = global.Client.ParseGroupName(name, global.Config.GetGroup())
	return
}

// GroupProvider calls GroupNameProvider then gets the named Group from the brain and attaches it to the Context.
func GroupProvider(c *Context) (err error) {
	if c.Group != nil {
		return
	}
	err = GroupNameProvider(c)
	if err != nil {
		return
	}

	c.Group, err = global.Client.GetGroup(c.GroupName)
	return
}

// UserNameProvider reads the NextArg and attaches it to the Context as UserName
func UserNameProvider(c *Context) (err error) {
	if c.UserName != nil {
		return
	}
	var username string
	username, err = c.NextArg()
	if err != nil {
		username, err = global.Config.Get("user")
		if username == "" {
			username = global.Client.GetSessionUser()
			err = nil
		}
	}
	c.UserName = &username
	return

}

// UserProvider calls UserNameProvider, gets the User from the brain, and attaches it to the Context.
func UserProvider(c *Context) (err error) {
	if c.User != nil {
		return
	}
	err = UserNameProvider(c)
	if err != nil {
		return
	}
	if err = AuthProvider(c); err != nil {
		return
	}
	c.User, err = global.Client.GetUser(*c.UserName)
	return
}

// VirtualMachineNameProvider reads the NextArg, parses it as a VirtualMachineName and attaches it to the Context
func VirtualMachineNameProvider(c *Context) (err error) {
	if err = AuthProvider(c); err != nil {
		return
	}

	if c.VirtualMachineName != nil {
		log.Log("VMNameProvider: VirtualMachineName already defined")
		return
	}
	name, err := c.NextArg()
	if err != nil {
		return err
	}

	c.VirtualMachineName, err = global.Client.ParseVirtualMachineName(name, global.Config.GetVirtualMachine())
	return
}

// VirtualMachineProvider calls VirtualMachineNameProvider then gets the named VirtualMachine from the brain and attaches it to the Context.
func VirtualMachineProvider(c *Context) (err error) {
	if c.VirtualMachine != nil {
		return
	}
	err = VirtualMachineNameProvider(c)
	if err != nil {
		return
	}
	c.VirtualMachine, err = global.Client.GetVirtualMachine(c.VirtualMachineName)
	return
}
