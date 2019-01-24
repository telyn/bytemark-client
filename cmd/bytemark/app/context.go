package app

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
)

// Context is a wrapper around urfave/cli.Context which provides easy access to
// the next unused argument and can have various bytemarky types attached to it
// in order to keep code DRY
type Context struct {
	Context        innerContext
	Account        *lib.Account
	Authed         bool
	Definitions    *lib.Definitions
	Disc           *brain.Disc
	Group          *brain.Group
	Privilege      brain.Privilege
	User           *brain.User
	VirtualMachine *brain.VirtualMachine

	currentArgIndex  int
	preprocessHasRun bool
}

// Reset replaces the Context with a blank one (keeping the cli.Context)
func (ctx *Context) Reset() {
	*ctx = Context{
		Context: ctx.Context,
	}
}

// App returns the cli.App that this context is part of. Usually this will be the same as global.App, but it's nice to depend less on globals.
func (ctx *Context) App() *cli.App {
	return ctx.Context.App()
}

// args returns all the args that were passed to the Context (i.e. all the args passed to this (sub)command)
func (ctx *Context) args() cli.Args {
	return ctx.Context.Args()
}

// Args returns all the unused arguments
func (ctx *Context) Args() []string {
	return ctx.args()[ctx.currentArgIndex:]
}

// Writer returns the app writer. just a convenience method for c.App().Writer
func (ctx *Context) Writer() io.Writer {
	return ctx.App().Writer
}

// ErrWriter returns the app writer. just a convenience method for c.App().ErrWriter
func (ctx *Context) ErrWriter() io.Writer {
	return ctx.App().ErrWriter
}

// Prompter returns the prompter which is used by this Context for prompting the user for input
func (ctx *Context) Prompter() util.Prompter {
	if prompter, ok := ctx.App().Metadata["prompter"].(util.Prompter); ok {
		return prompter
	}
	return nil
}

// Command returns the cli.Command this context is for
func (ctx *Context) Command() cli.Command {
	return ctx.Context.Command()
}

// Config returns the config attached to the App this Context is for
func (ctx *Context) Config() config.Manager {
	if config, ok := ctx.App().Metadata["config"].(config.Manager); ok {
		return config
	}
	return nil
}

// Client returns the API client attached to the App this Context is for
func (ctx *Context) Client() lib.Client {
	if client, ok := ctx.App().Metadata["client"].(lib.Client); ok {
		return client
	}
	return nil
}

// IsTest returns whether this app is being run as part of a test
// It uses the "buf" on the App's Metadata - which is added by
// app_test.BaseTestSetup and used to capture output for later assertions
func (ctx *Context) IsTest() bool {
	if _, ok := ctx.App().Metadata["buf"]; ok {
		return true
	}
	return false
}

// NextArg returns the next unused argument, and marks it as used.
func (ctx *Context) NextArg() (string, error) {
	if len(ctx.args()) <= ctx.currentArgIndex {
		return "", ctx.Help("not enough arguments were specified")
	}
	arg := ctx.args()[ctx.currentArgIndex]
	ctx.currentArgIndex++
	return arg, nil
}

// Help creates a UsageDisplayedError that will output the issue and a message to consult the documentation
func (ctx *Context) Help(whatsyourproblem string) (err error) {
	return util.UsageDisplayedError{TheProblem: whatsyourproblem, Command: ctx.Command().FullName()}
}

// IsSet returns true if the specified flag has been set.
func (ctx *Context) IsSet(flagName string) bool {
	return ctx.Context.IsSet(flagName)
}

// Flags provided by urfave/cli below.
// For all the flags in the flags package, see flags/accessors.go

// Bool returns the value of the named flag as a bool
func (ctx *Context) Bool(flagname string) bool {
	return ctx.Context.Bool(flagname)
}

// Int returns the value of the named flag as an int
func (ctx *Context) Int(flagname string) int {
	return ctx.Context.Int(flagname)
}

// Int64 returns the value of the named flag as an int64
func (ctx *Context) Int64(flagname string) int64 {
	return ctx.Context.Int64(flagname)
}

// String returns the value of the named flag as a string
func (ctx *Context) String(flagname string) string {
	if ctx.Context.IsSet(flagname) || ctx.Context.String(flagname) != "" {
		ctx.Debug("IsSet || String() != nil")
		return ctx.Context.String(flagname)
	}
	return ctx.Context.GlobalString(flagname)
}

// StringSlice returns the values of the named flag as a []string
func (ctx *Context) StringSlice(flagname string) []string {
	if ctx.Context.IsSet(flagname) || ctx.Context.String(flagname) != "" {
		return ctx.Context.StringSlice(flagname)
	}
	return ctx.Context.GlobalStringSlice(flagname)
}
