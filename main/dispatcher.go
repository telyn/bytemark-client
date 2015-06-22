package main

import (
	client "bigv.io/client/lib"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Dispatcher is used to determine what functions to run for the command-line arguments provided
type Dispatcher struct {
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
	bigv.SetDebugLevel(d.debugLevel)

	d.cmds = NewCommandSet(config, bigv)
	return d
}

// NewDispatcherWithCommands is for writing tests with mock CommandSets
func NewDispatcherWithCommands(config ConfigManager, commands Commands) *Dispatcher {
	d := NewDispatcher(config)
	d.cmds = commands
	return d
}

// EnsureAuth makes sure a valid token is stored in config.
// This should be called by anything that needs auth.

// Do takes the command line arguments and figures out what to do
func (d *Dispatcher) Do(args []string) {
	help := d.Flags.Lookup("help")
	fmt.Printf("%+v", help)
	if d.debugLevel >= 1 {
		fmt.Fprintf(os.Stderr, "Args passed to Do: %#v\n", args)
	}

	if help.Value.String() == "true" || len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Printf("No command specified.\n\n")
		d.cmds.Help(args)
		return
	}

	// short-circuit commands that don't require arguments
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
	case "create-group":
		d.cmds.CreateGroup(args[1:])
		return
	case "create-vm":
		d.cmds.CreateVM(args[1:])
		return
	case "debug":
		d.cmds.Debug(args[1:])
		return
	case "delete-vm":
		d.cmds.DeleteVM(args[1:])
		return
	case "restart":
		d.cmds.Restart(args[1:])
		return
	case "reset":
		d.cmds.ResetVM(args[1:])
		return
	case "show-account":
		d.cmds.ShowAccount(args[1:])
		return
	case "show-group":
		d.cmds.ShowGroup(args[1:])
		return
	case "show-vm":
		d.cmds.ShowVM(args[1:])
		return
	case "shutdown":
		d.cmds.Shutdown(args[1:])
		return
	case "start":
		d.cmds.Start(args[1:])
		return
	case "stop":
		d.cmds.Stop(args[1:])
		return
	case "undelete-vm":
		d.cmds.UndeleteVM(args[1:])
		return
	}
	fmt.Fprintf(os.Stderr, "Unrecognised command '%s'\r\n", args[0])
	d.cmds.Help(args)
}
