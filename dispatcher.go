package main

import (
	commands "bytemark.co.uk/client/cmds"
	"bytemark.co.uk/client/cmds/util"
	client "bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"flag"
	"strings"
)

// Dispatcher is used to determine what functions to run for the command-line arguments provided
type Dispatcher struct {
	Flags      *flag.FlagSet
	cmds       commands.CommandManager
	config     util.ConfigManager
	debugLevel int
}

// NewDispatcher creates a new Dispatcher given a config.
func NewDispatcher(config util.ConfigManager) (d *Dispatcher, err error) {
	d = new(Dispatcher)

	d.config = config
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

	d.cmds = commands.NewCommandSet(config, bigv)
	return d, nil
}

// NewDispatcherWithCommandManager is for writing tests with mock CommandManagers
func NewDispatcherWithCommandManager(config util.ConfigManager, commands commands.CommandManager) (*Dispatcher, error) {
	d, err := NewDispatcher(config)
	if err != nil {
		return nil, err
	}
	d.cmds = commands
	return d, nil
}

// CommandFunc is a type which takes an array of arguments and returns an util.ExitCode.
type CommandFunc func([]string) util.ExitCode

func (d *Dispatcher) DoAdd(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForAdd()
	}
	switch strings.ToLower(args[0]) {
	case "key":
		return d.cmds.AddKey(args[1:])
	}
	return d.cmds.HelpForAdd()
}

func (d *Dispatcher) DoCreate(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForCreate()
	}

	switch strings.ToLower(args[0]) {
	case "vm":
		return d.cmds.CreateVM(args[1:])
	case "group":
		return d.cmds.CreateGroup(args[1:])
	case "disc", "discs", "disk", "disks":
		return d.cmds.CreateDiscs(args[1:])
	case "key":
		return d.cmds.AddKey(args[1:])

	}
	log.Errorf("Unrecognised command 'create %s'\r\n", args[0])
	return util.E_PEBKAC
}

func (d *Dispatcher) DoDelete(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForDelete()
	}
	switch strings.ToLower(args[0]) {
	case "vm":
		return d.cmds.DeleteVM(args[1:])
	case "group":
		return d.cmds.DeleteGroup(args[1:])
	case "disc", "disk":
		return d.cmds.DeleteDisc(args[1:])
	case "key":
		return d.cmds.DeleteKey(args[1:])
	}
	log.Errorf("Unknown command 'delete %s'\r\n", args[0])
	return d.cmds.HelpForDelete()

}
func (d *Dispatcher) DoShow(args []string) util.ExitCode {
	// Show implements the show command which is a stupendous badass of a command
	if len(args) == 0 {
		d.cmds.HelpForShow()
		return util.E_USAGE_DISPLAYED
	}

	switch strings.ToLower(args[0]) {
	case "vm":
		return d.cmds.ShowVM(args[1:])
	case "account":
		return d.cmds.ShowAccount(args[1:])
	case "user":
		return d.cmds.ShowUser(args[1:])
	case "group":
		return d.cmds.ShowGroup(args[1:])
	}

	name := strings.TrimSuffix(args[0], d.config.EndpointName())
	dots := strings.Count(name, ".")
	switch dots {
	case 2:
		return d.cmds.ShowVM(args)
	case 1:
		return d.cmds.ShowGroup(args)
	case 0:
		return d.cmds.ShowAccount(args)
	}
	return util.E_SUCCESS
}

func (d *Dispatcher) DoUndelete(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForDelete()
	}
	switch strings.ToLower(args[0]) {
	case "vm":
		return d.cmds.UndeleteVM(args[1:])
	}
	log.Errorf("Unrecognised command 'undelete %s'\r\n", args[0])
	return d.cmds.HelpForDelete()
}

func (d *Dispatcher) DoList(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForList()
	}
	switch strings.ToLower(args[0]) {
	case "vms":
		return d.cmds.ListVMs(args[1:])
	case "discs", "disks":
		return d.cmds.ListDiscs(args[1:])
	case "groups":
		return d.cmds.ListGroups(args[1:])
	case "accounts":
		return d.cmds.ListAccounts(args[1:])
	case "key", "keys":
		return d.cmds.ListKeys(args[1:])
	}
	return d.cmds.HelpForList()
}

func (d *Dispatcher) DoLock(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForLocks()
	}

	switch strings.ToLower(args[0]) {
	case "hwprofile":
		return d.cmds.LockHWProfile(args[1:])
	}
	log.Errorf("Unrecognised command 'lock %s'\r\n", args[0])
	return d.cmds.HelpForLocks()
}

func (d *Dispatcher) DoUnlock(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForLocks()
	}

	switch strings.ToLower(args[0]) {
	case "hwprofile":
		return d.cmds.UnlockHWProfile(args[1:])
	}
	log.Errorf("Unrecognised command 'unlock %s'\r\n", args[0])
	return d.cmds.HelpForLocks()
}

func (d *Dispatcher) DoResize(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForResize()

	}
	switch strings.ToLower(args[0]) {
	case "disc", "disk":
		return d.cmds.ResizeDisc(args[1:])
	}

	return d.cmds.HelpForResize()
}

func (d *Dispatcher) DoSet(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForSet()
	}

	switch strings.ToLower(args[0]) {
	case "hwprofile":
		return d.cmds.SetHWProfile(args[1:])
	case "memory":
		return d.cmds.SetMemory(args[1:])
	case "cores":
		return d.cmds.SetCores(args[1:])
	}
	log.Errorf("Unrecognised command 'set %s'\r\n", args[0])
	return d.cmds.HelpForSet()
}

// Do takes the command line arguments and figures out what to do.
func (d *Dispatcher) Do(args []string) util.ExitCode {
	log.Debugf(1, "bytemark %v\r\n\r\n", client.GetVersion())

	log.Debugf(1, "Args passed to Do: %#v\r\n", args)

	for _, arg := range args {
		if arg == "/?" || arg == "/h" || arg == "--help" || arg == "-h" {
			d.cmds.Help(args)
			return util.E_USAGE_DISPLAYED
		}
	}

	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		log.Errorf("No command specified.\r\n")
		d.cmds.Help(args)
		return util.E_SUCCESS
	}

	commands := map[string]CommandFunc{
		"add":           d.DoAdd,
		"create":        d.DoCreate,
		"config":        d.cmds.Config,
		"console":       d.cmds.Console,
		"connect":       d.cmds.Console,
		"debug":         d.cmds.Debug,
		"delete":        d.DoDelete,
		"distributions": d.cmds.Distributions,
		"help":          d.cmds.Help,
		"hwprofiles":    d.cmds.HardwareProfiles,
		"images":        d.cmds.Distributions,
		"list":          d.DoList,
		"lock":          d.DoLock,
		"resize":        d.DoResize,
		"restart":       d.cmds.Restart,
		"reset":         d.cmds.ResetVM,
		"serial":        d.cmds.Console,
		"set":           d.DoSet,
		"shutdown":      d.cmds.Shutdown,
		"stop":          d.cmds.Stop,
		"storage":       d.cmds.StorageGrades,
		"start":         d.cmds.Start,
		"show":          d.DoShow,
		"undelete":      d.DoUndelete,
		"unlock":        d.DoUnlock,
		"vnc":           d.cmds.Console,
		"zones":         d.cmds.Zones,
	}

	command := strings.ToLower(args[0])
	fn := commands[command]

	if fn != nil {
		return fn(args[1:])
	} else {
		log.Errorf("Unrecognised command '%s'\r\n\r\n", command)

		return d.cmds.Help(args)
	}

}
