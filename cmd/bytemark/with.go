package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"net"
)

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
	ips := c.Context.Generic("ip")
	if ips, ok := ips.(*util.IPFlag); ok {
		*ips = make([]net.IP, 0)
	}
	disc := c.Context.Generic("disc")
	if disc, ok := disc.(*util.DiscSpecFlag); ok {
		*disc = make([]lib.Disc, 0)
	}
	size := c.Context.Generic("memory")
	if size, ok := size.(*util.SizeSpecFlag); ok {
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

// AccountNameProvider stitches the next argument to the context as AccountName
func AccountNameProvider(c *Context) (err error) {
	if c.AccountName != nil {
		return
	}

	if err = AuthProvider(c); err != nil {
		return
	}
	name, err := c.NextArg()
	if err != nil {
		return err
	}
	accName := global.Client.ParseAccountName(name, global.Config.GetIgnoreErr("account"))
	c.AccountName = &accName
	return
}

// AccountProvider uses AccountNameProvider and then gets the account details from the API, then stitches it to the context
func AccountProvider(c *Context) (err error) {
	err = AccountNameProvider(c)
	if err != nil {
		return
	}
	c.Account, err = global.Client.GetAccount(*c.AccountName)
	return
}

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

func DefinitionsProvider(c *Context) (err error) {
	if c.Definitions != nil {
		return
	}
	c.Definitions, err = global.Client.ReadDefinitions()
	return
}

func DiscLabelProvider(c *Context) (err error) {
	if c.DiscLabel != nil {
		return
	}
	discLabel, err := c.NextArg()
	if err != nil {
		return err
	}
	c.DiscLabel = &discLabel
	return
}

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

func UserProvider(c *Context) (err error) {
	if c.User != nil {
		return
	}
	err = UserNameProvider(c)
	if err != nil {
		return
	}
	c.User, err = global.Client.GetUser(*c.UserName)
	return
}

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
		log.Log(err)
		return err
	}

	c.VirtualMachineName, err = global.Client.ParseVirtualMachineName(name, global.Config.GetVirtualMachine())
	return
}

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
