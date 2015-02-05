// BUG(telyn): Needs more create-vm flags. Also boot-script isn't here yet.

// BUG(telyn): Flesh out the list of commands

// BUG(telyn): Can you remove a disk from a running VM?

// BUG(telyn): Not default to "wheezy" and 25GiB, default to whatever the API says are the defaults? https://projects.bytemark.co.uk/issues/9378

// BUG(telyn): List of privileges?

/*
BigV API Client


Basic Usage:

	bigv [flags] <command> [command-flags] [command args]


Common Flags:

Common flags may be placed anywhere

	--force         - Runs without prompting, except purges.
	--purge         - Runs purges without prompting.
	--yubikey       - Will prompt for a yubikey one-time pass.
	--no-yubikey    - Will not use BIGV_YUBIKEY. See below for more information on environment variables.
	--yubikey-otp   - Your yubikey one-time pass. Defaults to nothing.
	--endpoint      - The API endpoint of the BigV service you're trying to access. Without a URL scheme, assumes https. Defaults to uk0.bigv.io.
	--show-requests - Outputs requests made on stderr, useful for learning the BigV API.
	--json          - Silences usual output and outputs JSON instead, useful for using bigv from other applications or learning the BigV API.
	--user          - Your BigV username.
	--help          - Show this help

Both the `--flag=value` and `--flag value` forms are supported.

The BigV command will generally prompt you for your username and password.
Set the BIGV_USER and BIGV_PASS (and optionally BIGV_YUBIKEY) environment variables,
or the --user and --yubikey-otp flags in order to not receive prompting for these.

If BIGV_YUBIKEY or --yubikey-otp is set, it will be used unless --no-yubikey is specified.

Command flags trump environment variables, so if BIGV_USER and --user are specified, the value passed in --user will be used.

Where a VM, group or account name can be entered, you may enter the entire domain name of the machine, for example:
bigv.is.awesome.uk0.bigv.io for the VM "bigv" in group "is" of account "awesome". If a VM is in the default group of your
primary account, you may specify just the VM name, for example for a user with the default account "example", specifying
"awesomevm" is the same as specifying awesomevm.default.example or awesomevm.default.example.uk0.bigv.io
(if you haven't provided an --endpoint)

Dashes in commands may be replaced with spaces, for example: create vm is an alias for create-vm.

New is always an alias for create, for example: new-vm is an alias for create-vm.
You may also misspell disk as disc, at no penalty ;-)

The show-X and list-X commands perform the same API calls, but the show commands are designed to be detailed and human readable,
while the list-X commands are terse and designed to be machine readable.


Commands available:


COMMAND: create-vm

	create-vm [flags] <name> [image] [disk-specs]
	new-vm [flags] <name> [image] [disk-specs]

Creates a VM with the provided name, image and disk-specs.

Disk specs must match the following format: [storage-grade:]size[MgGtT] and are comma or space separated. The available storage grades can be listed with the list-storage-grades command. See the create-disk documentation below for the size suffixes.

If image and disk-specs are not provided, will interactively prompt for them, or default to "wheezy" and a 25GiB SATA SSD if the --force flag is present.

The VM's full specification will be given and you will be prompted to confirm it unless the --force flag is present.

Flags available:

	--boot-script - A filename of a boot-script to upload from the local machine. This script will be run the first time the virtual machine starts up.
	--hwprofile   - Sets the hardware profile. See the output of list-hwprofiles.
	--lock        - Locks the hardware profile, preventing automatic upgrades. Not set by default

COMMAND: delete

	delete <vm>

Deletes the specified VM. After some period the machine will purged, at which point it will no longer be possible to undelete.

This command will prompt you to confirm deletion unless --force is set, and will in either case output the time you have available before it is deleted.

COMMAND: undelete

	undelete <vm>

COMMAND: purge

	purge <vm>

Irreversibly deletes the given VM. Use with care!

This command will prompt you to confirm deletion unless --purge is set.

COMMAND: set-memory

COMMAND: set-cores

COMMAND: set-hwprofile

	set-hwprofile [--lock | --unlock] <profile>

Sets the hardware profile (and optionally un/locks it)

COMMAND: send-signal

	send-signal <vm> <signal>
	signal <vm> <signal>

COMMAND: create-disk

	create-disk [flags] <vm> <disk spec> [name]
	new-disk [flags] <vm> <disk spec> [name]

Creates a disk attached to the VM provided, with the given spec

Disk specs must match the following format: [storage-grade:]size[MgGtT]. The available storage grades can be listed with the list-storage-grades command.

Size must be a positive integer, and the suffixes are as follows:
	M - Megabytes - default
	g - GB  (1000 megabytes)
	G - GiB (1024 megabytes)
	t - TB  (1000 GB)
	T - TiB (1024 GiB)

The disk's specification will be given and you will be prompted to confirm it unless the --force flag is present.

Flags available:

	--grade - Specifies what storage grade to use. Overridden if the storage-grade is part of the disk-spec.


COMMAND: resize-disk

	resize-disk <vm> <disk> [size]

Resizes the disk named <disk> attached to the specified VM.
If size begins with a +, will increase the disk's size by the amount specified.

You'll be prompted to confirm the new size unless the --force flag is given.

Disks cannot be resized to a lower size than they currently have.


COMMAND: delete-disk

	delete-disk <vm> <disk>
	remove-disk <vm> <disk>

Removes the named disk from the specified VM.
This operation can probably only be performed on stopped VMs.

You will be prompted for confirmation unless the --force flag is given.


COMMAND: request-ip

	request-ip [reason]
	new-ip [reason]

Requests a new IP assignment. Will prompt for a reason if not specified.

COMMAND: delete-ip

	delete-ip

COMMAND: start / stop / restart /shutdown

	start <vm>
	shutdown <vm>
	stop <vm>
	restart <vm>

They all do what you expect. Stop and restart are both forceful, shutdown sends an ACPI shutdown message.


COMMAND: reimage

	reimage <vm> [image]

Reimages the given VM with the specified image, prompting if not specified. To see the list of images
available, run bigv show-images.


COMMAND: grant

	grant <user> <privilege> [object]

Grants the specified privilege (on the given object, if relevant) to that user.
A list of privileges available can be found... somewhere?


COMMAND: revoke

	revoke <user> <privilege> [object]

Grants the specified privilege (on the given object, if relevant) to that user.
A list of privileges available can be found... somewhere?

COMMAND: show-account
COMMAND: show-group
COMMAND: show-vm
COMMAND: show-ip
COMMAND: show-grants

COMMAND: list-accounts
COMMAND: list-groups
COMMAND: list-vms
COMMAND: list-ips
COMMAND: list-grants

COMMAND: privileges
COMMAND: storage-grades
COMMAND: images
COMMAND: hardware-profiles

COMMAND: connect

*/
package main
