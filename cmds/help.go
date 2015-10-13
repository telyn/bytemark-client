package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
	"strings"
)

// HelpForHelp shows overall usage information for the BigV client, including a list of available commands.
func (cmds *CommandSet) HelpForHelp() util.ExitCode {
	log.Log("Bytemark command-line client")
	log.Log()
	log.Log("Usage")
	log.Log()
	log.Log("    bytemark [flags] <command> [flags] [args]")
	log.Log()
	log.Log("See `bytemark help <command>` for help specific to a command")
	log.Log()
	log.Log("Commands available")
	log.Log()
	log.Log("    help [command | topic] - output the help for the client or for the given command or topic")
	log.Log()

	// ALL of this should be in a sweet datastructure
	// config
	log.Log("  Config commands:")
	log.Log("    config  output all info about the current config")
	log.Log("    config get <variable>  output the value & source of the given variable")
	log.Log("    config set <variable> <value>  persistently sets a bigv-client variable")
	log.Log("    config unset <variable> - persistently unsets a bigv-client variable")
	log.Log()

	// machines generally
	log.Log("  Common server commands:")
	log.Log("    console [--serial | --vnc] [--connect | --panel] <virtual machine>")
	log.Log("    request ip <virtual machine> <reason>")
	log.Log("    set rdns <ip> <host name>")
	log.Log("    shutdown [--force] <virtual machine> - if force given, immediately stop the VM.")
	log.Log("    start <virtual machine>")
	log.Log()

	// virtual machine
	log.Log("  Virtual machine commands:")
	log.Log("    create disc[s] [--account <account>] [--group <group>] [--size <size>] [--grade <storage grade>] <virtual machine> [<disc specs>]")
	log.Log("    create vm [flags] <name> [<cores> [<memory> [<disc specs>]]] - creates a vm. See `bytemark help create` for detail on the flags")
	log.Log("    delete disc [--force] [---purge] <virtual machine> <disc label>")
	log.Log("    delete vm [--force] [---purge] <virtual machine>")
	log.Log("    list discs <virtual machine> - lists the discs in the given VM, with their size and labels")
	log.Log("    list vms <group> - lists the vms in the given group, one per line")
	log.Log("    lock hwprofile <virtual machine>")
	log.Log("    reimage [--image <image>] <virtual machine> [<image>]")
	log.Log("    resize disc [--size <size>] <virtual machine> [<resize spec>] - resize to size. if ambiguous, berate user.")
	log.Log("    set cores <virtual machine> <num>")
	log.Log("    set hwprofile <virtual machine> <hardware profile>")
	log.Log("    set memory <virtual machine> <size>")
	log.Log("    show vm [--json] [--nics] <virtual machine> - shows an overview of the given VM. Its discs, IPs, and such.")
	log.Log("    undelete vm <virtual machine> - Bring an unpurged machine back from deletion")
	log.Log("    unlock hwprofile <virtual machine>")
	log.Log()

	//log.Log("  Dedicated host commands:")
	//log.Log("    None yet!")

	log.Log("  Account management commands:")
	log.Log("    create group [--account <account>] <name>")
	log.Log("    delete account <account>")
	log.Log("    delete group <group>")
	log.Log("    list accounts - lists the accounts you can see, one per line")
	log.Log("    list groups <account> - lists the groups in the given account, one per line")
	log.Log("    show account [--json] <account> - shows an overview of the given account, a list of groups and vms within them")
	log.Log("    show group [--json] <group> - shows an overview of the given group, a list of VMs in them w/ size information")
	log.Log()

	log.Log("  Informative commands:")
	log.Log("    list images - lists the available operating system images that can be passed to create vm and reimage")
	log.Log("    list (grades | storage-grades) - lists the available storage grades, along with a description.")
	log.Log("    list privileges - lists the privileges that can possibly be granted")

	return util.E_USAGE_DISPLAYED
}

// Help implements the help command, which gives usage information specific to each command. Usage: bytemark help [command]
func (cmds *CommandSet) Help(args []string) util.ExitCode {
	if len(args) == 0 {
		return cmds.HelpForHelp()
	}

	// please try and keep these in alphabetical order
	switch strings.ToLower(args[0]) {
	case "config":
		return cmds.HelpForConfig()
	case "create":
		return cmds.HelpForCreate()
	case "debug":
		return cmds.HelpForDebug()
	case "delete":
		return cmds.HelpForDelete()
	case "exit":
		return util.HelpForExitCodes()
	case "exit-codes":
		return util.HelpForExitCodes()
	case "show":
		return cmds.HelpForShow()
	}
	return cmds.HelpForHelp()

}
