<style>
code {
    font-size: 11pt
    }
</style>

Notation
========

In this document, words in angle brackets `<like so>` denote a positional argument that must be there.
Words in square brackets `[like so]` denote optional arguments / parts.
For example, in the disc spec in the below section, `[:<storage grade>]` means an optional colon, which if present must be followed by a storage grade
A pipe character `|` represents a choice of options, of which only one may be selected.

Argument specs
==============

`account`
---------------

	    <name | id>

`disc spec`
---------------

	    <size>[:<storage grade>]

`disc specs`
---------------

	    <disc spec>[,<disc specs>]

`group`
---------------

	    <name | id>[.<account>]

n.b. account will default to the account with the same name as the user you log in as

`hardware profile`
---------------

	    <name>

they are usually of the form virtioYYYY or compatibilityYYYY where YYYY is a 4-digit year.

host name
---------

A valid DNS name 
TODO(telyn): find rfc reference

`id`
---------------

A non-zero positive integer

`name`
---------------

A valid DNS name containing no `` `.` `` characters.

`resize spec`
---------------

    [+-]<size>

Will set the size to `<size>`, or if `+` is specified, grow the thing by `<size>`.
If `-` is specified, will shrink the thing by `<size>`. Or at least ask to.

`size`
---------------

	<int>[m[b]|g[b]]

Where int is a non-zero positive integer.
n.b. case insensitive.
The default unit is `g`. `g` is gibibytes, `m` is mebibytes.

`storage grade`
---------------

	    sata | archive

`virtual machine`
---------------

	    <name | id>[.<group>]

n.b. group will default to "default"



n.b. storage grades can be found out by asking the api for /definitions.json

List of commands
================

`bigv config` - outputs the current config
`bigv config set <variable> <value>`  persistently sets a variable
`bigv config unset <variable>` - persistently unsets a variable
`bigv console [--serial | --vnc] [--connect | --panel] <virtual machine>`
`bigv create group [--account <account>] <name>`
`bigv create dis<c|k>[s] [--account <account>] [--group <group>] [--size <size>] [--grade <storage grade>] <virtual machine> [<disc specs>]` - if ambiguous, berate user
`bigv create ip [--reason reason] <virtual machine>`
`bigv create vm [flags] <name> [<cores> [<memory> [<disc specs>]]]` - creates a vm with the given name and the following flags

	    --account <name>
	    --cores <num> (default 1)
	    --cdrom <url>
	    --discs <disc specs> (default 25)
	    --force
	    --group <name>
	    --hwprofile <profile>
	    --hwprofile-locked (if specified, will lock the hwprofile)
	    --image <image name> 
	    --memory <size> (default 1, default unit GB)
	    --public-keys <keys> (newline seperated)
	    --public-keys-file <file> (will be read & appended to --public-keys)
	    --root-password <password>
	    --stopped (if set, machine won't be started)
	    --zone <name> (default york)
`bigv delete [--force] [--purge] <name>
`bigv delete account <account>
`bigv delete dis<c|k> [--force] [---purge] <virtual machine> <disc label>`
`bigv delete group <group>
`bigv delete ip <ip>` - _actually no you can't delete an ip_
`bigv delete nic <virtual machine> <nic id>`
`bigv delete user <user>
`bigv delete vm [--force] [---purge] <virtual machine>`
`bigv debug [--junk-token] [--auth] <method> <path>` - Make an HTTP request to the given path on the current endpoint.
`bigv debug config` - output the current config as json to debug Config's internal state
`bigv grant <user> <privilege> <object>`
`bigv help [command | topic]` - output the help for bigv or for the given command or topic
`bigv lock hwprofile <virtual machine>`
`bigv list accounts` - lists the accounts you can see, one per line
`bigv list images` - lists the available operating system images that can be passed to create vm and reimage
`bigv list (grades | storage-grades)` - lists the available storage grades, along with a description. One per line.
`bigv list privileges` - lists the privileges that can possibly be granted
`bigv list groups <account>` - lists the groups in the given account, one per line
`bigv list vms <group>` - lists the vms in the given group, one per line
`bigv reimage [--image <image>] <virtual machine> [<image>]`
`bigv request ip <virtual machine> [<nic id>]` - requests an IP on the given NIC, or the default NIC if not specified
`bigv reset <virtual machine>` - Need to discuss whether this is useful
`bigv resize dis<c|k> [--size <size>] <virtual machine> [<resize spec>]` - resize to `size`. if ambiguous, berate user.
`bigv revoke <user> <privilege>`
`bigv serial [--connect] <virtual machine>` - alias to `bigv console --serial`
`bigv set cores <virtual machine> <num>`
`bigv set hwprofile <virtual machine> <hardware profile>`
`bigv set memory <virtual machine> <size>`
`bigv set rdns <ip> <host name>`
`bigv show account [--json] <account>` - shows an overview of the given account, a list of groups and vms within them
`bigv show group [--json] <group>` - shows an overview of the given group, a list of VMs in them w/ size information
`bigv show user <name>` - shows details about the given user - their authorised keys and any privileges you have granted them
`bigv show vm [--json] [--nics] <virtual machine>` - shows an overview of the given VM. Its discs, IPs, and such.
`bigv shutdown <virtual machine>`
`bigv start <virtual machine>`
`bigv stop <virtual machine>`
`bigv undelete vm <virtual machine>`
`bigv unlock hwprofile <virtual machine>`
`bigv vnc [--connect | --panel] <virtual machine>` - alias for `bigv console --vnc`
