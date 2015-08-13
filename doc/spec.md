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

`id`
---------------

A non-zero positive integer

`name`
---------------

A valid DNS name containing no `` `.` `` characters.

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
`bigv create disc [--account <account>] [--group <group>] [--size <size>] [--grade <storage grade>] <virtual machine>`
`bigv create group [--account <account>] <name>`
`bigv create disc[s] <virtual machine> <disc specs>`
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
	    --memory <num> (default 1, unit GB)
	    --public-keys <keys> (newline seperated)
	    --public-keys-file <file> (will be read & appended to --public-keys)
	    --root-password <password>
	    --stopped (if set, machine won't be started)
	    --zone <name> (default york)

`bigv list vms <group> - lists the vms in the given group, one per line
`bigv list groups <account> - lists the groups in the given account, one per line
`bigv list accounts - lists the accounts you can see, one per line
`bigv lock hwprofile <virtual machine>`
`bigv serial [--connect] <virtual machine>` - alias to `bigv console --serial`
`bigv set hwprofile <virtual machine> <hardware profile>`
`bigv show vm [--json] <virtual machine> - shows an overview of the given VM. Its discs, IPs, and such.
`bigv show group [--json] <group> - shows an overview of the given group, a list of VMs in them w/ size information
`bigv show account [--json] <account> - shows an overview of the given account, a list of groups and vms within them
`bigv show user <name>` - shows details about the given user - their authorised keys and any privileges you have granted them
`bigv unlock hwprofile <virtual machine>`
`bigv vnc [--connect | --panel] <virtual machine>` - alias for `bigv console --vnc`
