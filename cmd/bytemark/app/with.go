package app

import (
	"net"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
)

// ProviderFunc is the function type that can be passed to With()
type ProviderFunc func(*Context) error

// With is a convenience function for making cli.Command.Actions that sets up a Context, runs all the providers, cleans up afterward and returns errors from the actions if there is one
func With(providers ...ProviderFunc) func(c *cli.Context) error {
	providers = append(providers, providers[len(providers)-1])
	providers[len(providers)-2] = (*Context).Preprocess
	return func(cliContext *cli.Context) error {
		c := Context{Context: CliContextWrapper{cliContext}}
		defer cleanup(&c)

		err := foldProviders(&c, providers...)
		return err
	}
}

// Preprocess runs the Preprocess methods on all flags that implement Preprocessor
func (c *Context) Preprocess() error {
	if c.preprocessHasRun {
		return nil
	}
	c.Debug("Preprocessing\n")
	for _, flag := range c.Command().Flags {
		if gf, ok := flag.(cli.GenericFlag); ok {
			if pp, ok := gf.Value.(Preprocesser); ok {
				c.Debug("Doing some shit to %s\n", flag.GetName())
				c.Debug("b4: %#v ", gf.Value)
				err := pp.Preprocess(c)
				if err != nil {
					return err
				}
				c.Debug("after: %#v\n", gf.Value)
			}
		}
	}
	c.preprocessHasRun = true
	return nil
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

	account, ok := c.Context.Generic("account").(*AccountNameFlag)
	if ok {
		*account = AccountNameFlag{}
	}
}

// foldProviders runs all the providers with the given context, stopping if there's an error
func foldProviders(c *Context, providers ...ProviderFunc) (err error) {
	for i, provider := range providers {
		c.Debug("Provider #%d (%v)n\n", i, provider)
		err = provider(c)
		if err != nil {
			return
		}
	}
	return
}
