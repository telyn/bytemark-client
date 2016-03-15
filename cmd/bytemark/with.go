package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"github.com/codegangsta/cli"
)

func WithVirtualMachineName(fn func(*cli.Context, *lib.VirtualMachineName)) func(c *cli.Context) {
	return func(c *cli.Context) {
		args := c.Args()
		if !args.Present() {
			global.Error = &util.PEBKACError{}
			return
		}
		name, err := global.Client.ParseVirtualMachineName(args.First())
		if err != nil {
			global.Error = err
			return
		}
		fn(c, &name)
	}
}

func WithGroupName(fn func(*cli.Context, *lib.GroupName)) func(c *cli.Context) {
	return func(c *cli.Context) {
		args := c.Args()
		if !args.Present() {
			global.Error = util.PEBKACError{}
		}

		name := global.Client.ParseGroupName(args.First())
		fn(c, &name)
	}
}
