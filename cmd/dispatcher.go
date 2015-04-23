package cmd

import (
	client "bigv.io/client/lib"
	"flag"
	"fmt"
	"os"
	"strings"
)

// type der is used to create API requests and direct output to views,
// except probably when those API requests don't require authorisation (e.g. /definitions, new user)
type Dispatcher struct {
	//Config *Config
	Flags      *flag.FlagSet
	cmds       Commands
	debugLevel int
}

// NewDispatcher creates a new Dispatcher given a config.
func NewDispatcher(config ConfigManager) (d *Dispatcher) {
	d = new(Dispatcher)
	bigv, err := client.New(config.Get("endpoint"))
	if err != nil {
		exit(err)
	}

	d.debugLevel = config.GetDebugLevel()

	d.cmds = NewCommandSet(config, bigv)
	return d
}

// NewderWithCommands is for writing tests with mock CommandSets
func NewDispatcherWithCommands(config ConfigManager, commands Commands) *Dispatcher {
	d := NewDispatcher(config)
	d.cmds = commands
	return d
}

// EnsureAuth makes sure a valid token is stored in config.
// This should be called by anything that needs auth.

// TODO(telyn): Write a test for Do. Somehow.

// Do takes the command line arguments and figures out what to do
func (d *Dispatcher) Do(args []string) {
	//	help := d.Flags.Lookup("help")
	///	fmt.Printf("%+v", help)
	if d.debugLevel >= 1 {
		fmt.Fprintf(os.Stderr, "Args passed to Do: %#v\n", args)
	}

	if /*help == true || */ len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Printf("No command specified.\n\n")
		d.cmds.Help(args)
		return
	}

	// short-circuit commands that don't take arguments
	switch strings.ToLower(args[0]) {
	case "config":
		d.cmds.Config(args[1:])
		return
	case "help":
		d.cmds.Help(args[1:])
		return
	}

	// do this
	if len(args) == 1 {
		d.cmds.Help(args)
		return
	}

	switch strings.ToLower(args[0]) {
	case "debug":
		d.cmds.Debug(args[1:])
		return
	case "show-account":
		d.cmds.ShowAccount(args[1:])
		return
	case "show-vm":
		d.cmds.ShowVM(args[1:])
		return
	}
	fmt.Fprintf(os.Stderr, "Unrecognised command '%s'\r\n", args[0])
	d.cmds.Help(args)
}
