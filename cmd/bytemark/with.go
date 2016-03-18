package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"github.com/codegangsta/cli"
)

type Context struct {
	Context            *cli.Context
	AccountName        *string
	Account            *lib.Account
	Authed             bool
	Definitions        *lib.Definitions
	DiscLabel          *string
	GroupName          *lib.GroupName
	Group              *lib.Group
	User               *lib.User
	UserName           *string
	VirtualMachine     *lib.VirtualMachine
	VirtualMachineName *lib.VirtualMachineName

	currentArgIndex int
}

func (c *Context) args() cli.Args {
	return c.Context.Args()
}

func (c *Context) NextArg() (string, error) {
	if len(c.args()) <= c.currentArgIndex {
		return "", util.NotEnoughArgumentsError{}
	}
	arg := c.args()[c.currentArgIndex]
	c.currentArgIndex++
	return arg, nil
}

type ProviderFunc func(*Context) error

func With(providers ...ProviderFunc) func(c *cli.Context) {
	return func(cliContext *cli.Context) {
		c := Context{Context: cliContext}
		global.Error = foldProviders(&c, providers...)
	}
}

func foldProviders(c *Context, providers ...ProviderFunc) (err error) {
	for _, provider := range providers {
		err = provider(c)
		if err != nil {
			return
		}
	}
	return
}

func AccountNameProvider(c *Context) (err error) {
	if c.AccountName != nil {
	}
	name, err := c.NextArg()
	if err != nil {
		return err
	}
	accName := global.Client.ParseAccountName(name)
	c.AccountName = &accName
	return
}

func AccountProvider(c *Context) (err error) {
	err = foldProviders(c, AccountNameProvider, AuthProvider)
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

	name, err := c.NextArg()
	if err != nil {
		return err
	}
	c.GroupName = global.Client.ParseGroupName(name)
	return
}

func GroupProvider(c *Context) (err error) {
	if c.Group != nil {
		return
	}
	err = foldProviders(c, GroupNameProvider, AuthProvider)
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
	if c.VirtualMachineName != nil {
		return
	}
	name, err := c.NextArg()
	if err != nil {
		return err
	}

	c.VirtualMachineName, err = global.Client.ParseVirtualMachineName(name)
	return
}

func VirtualMachineProvider(c *Context) (err error) {
	if c.VirtualMachine != nil {
		return
	}
	err = foldProviders(c, VirtualMachineNameProvider, AuthProvider)
	if err != nil {
		return
	}
	c.VirtualMachine, err = global.Client.GetVirtualMachine(c.VirtualMachineName)
	return
}
