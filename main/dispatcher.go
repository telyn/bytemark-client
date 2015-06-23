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
func NewDispatcher(config ConfigManager) (d *Dispatcher, err error) {
	d = new(Dispatcher)
	endpoint, err := config.Get("endpoint")
	if err != nil {
		return nil, err
	}
	bigv, err := client.New(endpoint)
	if err != nil {
		return nil, err
	}

	d.debugLevel = config.GetDebugLevel()
	bigv.SetDebugLevel(d.debugLevel)

	d.cmds = NewCommandSet(config, bigv)
	return d, nil
}

// NewDispatcherWithCommands is for writing tests with mock CommandSets
func NewDispatcherWithCommands(config ConfigManager, commands Commands) (*Dispatcher, error) {
	d, err := NewDispatcher(config)
	if err != nil {
		return nil, err
	}
	d.cmds = commands
	return d, nil
}

// EnsureAuth makes sure a valid token is stored in config.
// This should be called by anything that needs auth.

// Do takes the command line arguments and figures out what to do
func (d *Dispatcher) Do(args []string) ExitCode {
	if d.debugLevel >= 1 {
		fmt.Fprintf(os.Stderr, "Args passed to Do: %#v\n", args)
	}

	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Printf("No command specified.\n\n")
		d.cmds.Help(args)
		return E_SUCCESS
	}

	// short-circuit commands that don't require arguments
	switch strings.ToLower(args[0]) {
	case "config":
		return d.cmds.Config(args[1:])
	case "create-group":
		return d.cmds.CreateGroup(args[1:])
	case "create-vm":
		return d.cmds.CreateVM(args[1:])
	case "help":
		d.cmds.Help(args[1:])
		return E_USAGE_DISPLAYED
	}

	// do this
	if len(args) == 1 {
		d.cmds.Help(args)
		return E_USAGE_DISPLAYED
	}

	switch strings.ToLower(args[0]) {
	case "debug":
		return d.cmds.Debug(args[1:])
	case "delete-vm":
		return d.cmds.DeleteVM(args[1:])
	case "restart":
		return d.cmds.Restart(args[1:])
	case "reset":
		return d.cmds.ResetVM(args[1:])
	case "show-account":
		return d.cmds.ShowAccount(args[1:])
	case "show-group":
		return d.cmds.ShowGroup(args[1:])
	case "show-vm":
		return d.cmds.ShowVM(args[1:])
	case "shutdown":
		return d.cmds.Shutdown(args[1:])
	case "start":
		return d.cmds.Start(args[1:])
	case "stop":
		return d.cmds.Stop(args[1:])
	case "undelete-vm":
		return d.cmds.UndeleteVM(args[1:])
	}
	fmt.Fprintf(os.Stderr, "Unrecognised command '%s'\r\n", args[0])
	d.cmds.Help(args)
	return E_USAGE_DISPLAYED
}
