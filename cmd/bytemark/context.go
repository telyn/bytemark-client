package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"github.com/codegangsta/cli"
	"net"
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

func (c *Context) Args() []string {
	return c.args()[c.currentArgIndex:]
}

func (c *Context) NextArg() (string, error) {
	if len(c.args()) <= c.currentArgIndex {
		return "", util.NotEnoughArgumentsError{}
	}
	arg := c.args()[c.currentArgIndex]
	c.currentArgIndex++
	return arg, nil
}

func (c *Context) Bool(flagname string) bool {
	return c.Context.Bool(flagname)
}

func (c *Context) Discs(flagname string) []lib.Disc {
	disc := c.Context.Generic(flagname)
	if disc, ok := disc.(*util.DiscSpecFlag); ok {
		return []lib.Disc(*disc)
	}
	return []lib.Disc{}
}

func (c *Context) FileName(flagname string) string {
	file := c.Context.Generic(flagname)
	if file, ok := file.(*util.FileFlag); ok {
		return file.FileName
	}
	return ""
}

func (c *Context) FileContents(flagname string) string {
	file := c.Context.Generic(flagname)
	if file, ok := file.(*util.FileFlag); ok {
		return file.Value
	}
	return ""
}

func (c *Context) Int(flagname string) int {
	return c.Context.Int(flagname)
}

func (c *Context) IPs(flagname string) []net.IP {
	ips := c.Context.Generic(flagname)
	if ips, ok := ips.(*util.IPFlag); ok {
		return []net.IP(*ips)
	}
	return []net.IP{}
}

func (c *Context) String(flagname string) string {
	return c.Context.String(flagname)
}

func (c *Context) Size(flagname string) int {
	size := c.Context.Generic(flagname)
	if size, ok := size.(*util.SizeSpecFlag); ok {
		return int(*size)
	}
	return 0
}
