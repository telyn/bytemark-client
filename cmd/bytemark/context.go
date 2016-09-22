package main

import (
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"net"
)

// Context is a wrapper around urfave/cli.Context which provides easy access to
// the next unused argument and can have various bytemarky types attached to it
// in order to keep code DRY
type Context struct {
	Context            *cli.Context
	AccountName        *string
	Account            *lib.Account
	Authed             bool
	Definitions        *lib.Definitions
	DiscLabel          *string
	GroupName          *lib.GroupName
	Group              *brain.Group
	User               *brain.User
	UserName           *string
	VirtualMachine     *brain.VirtualMachine
	VirtualMachineName *lib.VirtualMachineName

	currentArgIndex int
}

// args returns all the args that were passed to the Context (i.e. all the args passed to this (sub)command)
func (c *Context) args() cli.Args {
	return c.Context.Args()
}

// Args returns all the unused arguments
func (c *Context) Args() []string {
	return c.args()[c.currentArgIndex:]
}

// NextArg returns the next unused argument, and marks it as used.
func (c *Context) NextArg() (string, error) {
	if len(c.args()) <= c.currentArgIndex {
		return "", util.NotEnoughArgumentsError{}
	}
	arg := c.args()[c.currentArgIndex]
	c.currentArgIndex++
	return arg, nil
}

// Help returns the Help for this Context (i.e. command or subcommand) with the given string prepended with a couple of newlines
func (c *Context) Help(whatsyourproblem string) (err error) {
	log.Output(whatsyourproblem, "")
	err = cli.ShowSubcommandHelp(c.Context)
	if err != nil {
		return
	}
	return util.UsageDisplayedError{TheProblem: whatsyourproblem}
}

// flags below

// Bool returns the value of the named flag as a bool
func (c *Context) Bool(flagname string) bool {
	return c.Context.Bool(flagname)
}

// Discs returns the discs passed along as the named flag.
// I can't imagine why I'd ever name a disc flag anything other than --disc, but the flexibility is there just in case.
func (c *Context) Discs(flagname string) []brain.Disc {
	disc := c.Context.Generic(flagname)
	if disc, ok := disc.(*util.DiscSpecFlag); ok {
		return []brain.Disc(*disc)
	}
	return []brain.Disc{}
}

// FileName returns the name of the file given by the named flag
func (c *Context) FileName(flagname string) string {
	file := c.Context.Generic(flagname)
	if file, ok := file.(*util.FileFlag); ok {
		return file.FileName
	}
	return ""
}

// FileContents returns the contents of the file given by the named flag.
func (c *Context) FileContents(flagname string) string {
	file := c.Context.Generic(flagname)
	if file, ok := file.(*util.FileFlag); ok {
		return file.Value
	}
	return ""
}

// Int returns the value of the named flag as an int
func (c *Context) Int(flagname string) int {
	return c.Context.Int(flagname)
}

// IPs returns the ips passed along as the named flag.
func (c *Context) IPs(flagname string) []net.IP {
	ips := c.Context.Generic(flagname)
	if ips, ok := ips.(*util.IPFlag); ok {
		return []net.IP(*ips)
	}
	return []net.IP{}
}

// String returns the value of the named flag as a string
func (c *Context) String(flagname string) string {
	return c.Context.String(flagname)
}

// Size returns the value of the named SizeSpecFlag as an int in megabytes
func (c *Context) Size(flagname string) int {
	size := c.Context.Generic(flagname)
	if size, ok := size.(*util.SizeSpecFlag); ok {
		return int(*size)
	}
	return 0
}

// IfNotMarshalJSON checks to see if the json flag was set, and outputs obj as a JSON object if so.
// if not, runs fn
func (c *Context) IfNotMarshalJSON(obj interface{}, fn func() error) error {
	if c.Bool("json") {
		js, err := json.MarshalIndent(obj, "", "    ")
		if err != nil {
			return err
		}
		log.Output(string(js))
		return nil
	}
	return fn()
}
