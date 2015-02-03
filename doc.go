// BUG(telyn): Needs more create-vm flags. Also boot-script isn't here yet.
// BUG(telyn): Flesh out the list of commands
// BUG(telyn): Can you remove a disk from a running VM?
// BUG(telyn): Not default to "wheezy" and 25GiB, default to whatever the API says are the defaults? https://projects.bytemark.co.uk/issues/9378
/*
BigV API Client

Basic Usage:

	bigv [flags] <command> [command-flags] [command args]

	Common Flags:

		Common flags may be placed anywhere

		--force        - Runs without prompting, except purges.
		--purge        - Runs purges without prompting.
		--yubikey      - Will prompt for a yubikey one-time pass.
		--no-yubikey   - Will not use BIGV_YUBIKEY. See below for more information on environment variables.
		--yubikey-otp  - Your yubikey one-time pass. Defaults to nothing.
		--endpoint     - The API endpoint of the BigV service you're trying to access. Without a URL scheme, assumes https. Defaults to uk0.bigv.io.
		--user         - Your BigV username.
		--help         - Show this hel

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

Commands available:

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

	resize-disk <vm> <disk> [size]

		Resizes the disk named <disk> attached to the specified VM.
		If size begins with a +, will increase the disk's size by the amount specified.

		You'll be prompted to confirm the new size unless the --force flag is given.

		Disks cannot be resized to a lower size than they currently have.

	delete-disk <vm> <disk>
	remove-disk <vm> <disk>

		Removes the named disk from the specified VM.
		This operation can probably only be performed on stopped VMs.

		You will be prompted for confirmation unless the --force flag is given.

*/
package main
