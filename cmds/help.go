package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
	"strings"
)

func (cmds *CommandSet) HelpForTopic(topic string) util.ExitCode {
	topics := map[string]string{
		"profiles": `Having multiple configurations with Bytemark client is useful if you regularly log in as two different users,
or to different instances of the Bytemark API. One can set up and use different configurations with the --config-dir global flag.

For example, to set up your default configuration to log in by default as 'alice', and one configuration where you log in as 'bob'
with a yubikey, run the following commands (windows users, note that you'll need to use --config-dir="%HOME%\.bob"):

    bytemark config set user alice
    bytemark --config-dir="$HOME/.bob" set user bob
    bytemark --config-dir="$HOME/.bob" set yubikey

At this point you can set up an alias to use your 'bob' configuration. Say you use bash/zsh, add the following to your bashrc/zshrc:
    alias bytemark-bob='bytemark --config-dir="$HOME/.bob"'

Now you can run 'bytemark-bob list vms' to list all the vms in bob's default account and 'bytemark list vms' to do the same for alice.

Sorted.
`,
		"scripting": `The Bytemark client has been programmed from the beginning to attempt to make it easy for users to script with it.

Some particularly relavent notes: 
    * The 'list' command is entirely designed for scripting - it outputs to stdout, one item per line.
    * The reimage and create server commands only print the root password to stdout - all other output is sent to stderr.
    * All exit codes are documented - see the help topic exit codes.
    * If you're a fan of jq or you want to use bytemark-client within a OO scripting language, you can get json output from several commands like show using --json.

Here are just a couple of tricks I've been able to come up with.

To output the uptime for all your machines in the "critical" group:
  for i in $(bytemark list vms critical); do echo "${i%%.*}:"; ssh $i uptime; done

To add 10GB of space to each archive grade disk in your "storage" vm:
  for disc in $(bytemark list discs storage | grep "archive grade"); do bytemark resize disc --size +10G $machine $(awk '{print $2}'); done

To list all my VMs that have a disc bigger than the default 25GiB:
    bytemark show account --json telyn | jq '[.groups[].virtual_machines[] | select(.discs[] | .size > 25600) | .hostname ]' | uniq
`,
	}
	log.Log(topics[topic])
	return util.E_USAGE_DISPLAYED
}

// HelpForHelp shows overall usage information for the BigV client, including a list of available commands.
func (cmds *CommandSet) HelpForHelp() util.ExitCode {
	log.Log("usage: bytemark [flags] <command> [flags] [args]")
	log.Log()
	log.Log("See `bytemark help <command | topic>` for help specific to a command or topic")
	log.Log()
	log.Log("Help topics available: profiles, scripting, exit codes")
	log.Log()
	log.Log("Commands available")
	log.Log("   Config commands: config, config get, config set, config unset")
	log.Log("   Common commands: create, delete, list, show, undelete")
	log.Log("   Server commands: console, lock hwprofile, reimage, request ip, set, shutdown, start")
	log.Log("   Cloud disk commands: resize")
	log.Log("   User commands: add key, remove key")
	log.Log("   Information commands: hwprofiles, images, storage grades, version, zones")

	/*log.Log()
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

	// users
	log.Log("  User management commands:")
	log.Log("    add key [--key-user=<user>] <public key> - add an SSH public key to the given user, or you by default")
	//log.Log("    grant <user> <privilege> <object>")
	log.Log("    list keys [--key-user=<user>] - list the SSH keys authorised for management by the given user. Defaults to showing you your keys")
	log.Log("    remove key [--key-user=<user>] <public key identifier> - remove the given key from the given user, defaulting to you.")
	//log.Log("    revoke <user> <privilege>")
	log.Log("    show user <name> - shows details about the given user, including their authorised keys") //and any privileges you have granted them.")
	log.Log()

	log.Log("  Informative commands:")
	log.Log("    images - lists the available operating system images that can be passed to create vm and reimage")
	log.Log("    storage grades - lists the available storage grades, along with a description.")
	log.Log("    privileges - lists the privileges that can possibly be granted")
	*/

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
	case "console":
		return cmds.HelpForConsole()
	case "create":
		if len(args) > 1 {
			switch args[1] {
			case "vm", "server":
				return cmds.HelpForCreateVM()
			}
		}
		return cmds.HelpForCreate()
	case "debug":
		return cmds.HelpForDebug()
	case "delete", "undelete":
		return cmds.HelpForDelete()
	case "exit":
		return util.HelpForExitCodes()
	case "hwprofiles":
		return cmds.HardwareProfiles([]string{})
	case "images", "distributions":
		return cmds.Distributions([]string{})
	case "list":
		return cmds.HelpForList()
	case "lock", "unlock":
		return cmds.HelpForLocks()
	case "profiles":
		return cmds.HelpForTopic("profiles")
	case "reimage":
		return cmds.HelpForReimage()
	case "restart", "reset", "shutdown", "stop", "power":
		return cmds.HelpForPower()
	case "resize":
		return cmds.HelpForResize()
	case "scripting":
		return cmds.HelpForTopic("scripting")
	case "set":
		return cmds.HelpForSet()
	case "show":
		return cmds.HelpForShow()
	case "storage":
		return cmds.StorageGrades([]string{})
	case "zones":
		return cmds.Zones([]string{})

	}
	return cmds.HelpForHelp()

}
