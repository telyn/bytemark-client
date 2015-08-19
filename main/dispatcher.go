package main

import (
	commands "bigv.io/client/cmds"
	util "bigv.io/client/cmds/util"
	client "bigv.io/client/lib"
	"flag"
	"fmt"
	"os"
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

func (d *Dispatcher) DoCreate(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForCreate()
	}

	switch strings.ToLower(args[0]) {
	case "vm":
		return d.cmds.CreateVM(args[1:])
	case "group":
		return d.cmds.CreateGroup(args[1:])
		//    case "disc", "discs"
		//    return d.cmds.CreateDiscs(args[1:]

	}
	fmt.Fprintf(os.Stderr, "Unrecognised command 'create %s'\r\n", args[0])
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
	case "disc":
		return d.cmds.DeleteDisc(args[1:])
	}
	fmt.Fprintf(os.Stderr, "Unknown command 'delete %s'\r\n", args[0])
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
		fmt.Printf("Leave me alone! I'm grumpy.")
		return 666
		//return ShowUser(args[1:])
	case "group":
		return d.cmds.ShowGroup(args[1:])
	case "key", "keys":
		fmt.Printf("Leave me alone, I'm grumpy!")
		return 666
		//return d.cmds.ShowKeys(args[1:])
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
		// TODO: should also try show-vm sprintf("%s.%s.%s", args[0], "default", config.get("user"))
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
	fmt.Fprintf(os.Stderr, "Unrecognised command 'undelete %s'\r\n", args[0])
	return d.cmds.HelpForDelete()
}

func (d *Dispatcher) DoList(args []string) util.ExitCode {
	if len(args) == 0 {
		return d.cmds.HelpForList()
	}
	switch strings.ToLower(args[0]) {
	case "vms":
		return d.cmds.ListVMs(args[1:])
	case "discs":
		return d.cmds.ListDiscs(args[1:])
	case "groups":
		return d.cmds.ListGroups(args[1:])
	case "accounts":
		return d.cmds.ListAccounts(args[1:])
		//case "keys":
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
	fmt.Fprintf(os.Stderr, "Unrecognised command 'lock %s'\r\n", args[0])
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
	fmt.Fprintf(os.Stderr, "Unrecognised command 'unlock %s'\r\n", args[0])
	return d.cmds.HelpForLocks()
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
	fmt.Fprintf(os.Stderr, "Unrecognised command 'set %s'\r\rn", args[0])
	return d.cmds.HelpForSet()
}

// Do takes the command line arguments and figures out what to do.
func (d *Dispatcher) Do(args []string) util.ExitCode {
	if d.debugLevel >= 1 {
		fmt.Fprintf(os.Stderr, "Args passed to Do: %#v\n", args)
	}

	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Printf("No command specified.\n\n")
		d.cmds.Help(args)
		return util.E_SUCCESS
	}

	commands := map[string]CommandFunc{
		"create":   d.DoCreate,
		"config":   d.cmds.Config,
		"console":  d.cmds.Console,
		"connect":  d.cmds.Console,
		"debug":    d.cmds.Debug,
		"delete":   d.DoDelete,
		"help":     d.cmds.Help,
		"list":     d.DoList,
		"lock":     d.DoLock,
		"restart":  d.cmds.Restart,
		"reset":    d.cmds.ResetVM,
		"serial":   d.cmds.Console,
		"set":      d.DoSet,
		"shutdown": d.cmds.Shutdown,
		"stop":     d.cmds.Stop,
		"start":    d.cmds.Start,
		"show":     d.DoShow,
		"undelete": d.DoUndelete,
		"unlock":   d.DoUnlock,
		"vnc":      d.cmds.Console,
	}

	command := strings.ToLower(args[0])
	fn := commands[command]

	if fn != nil {
		return fn(args[1:])
	} else {
		fmt.Fprintf(os.Stderr, "Unrecognised command '%s'\r\n", command)
		fmt.Println()
		return d.cmds.Help(args)
	}

}
