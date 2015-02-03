// BUG(telyn): Unsure of the default hwprofile
// BUG(telyn): Needs more create-vm flags. Also boot-script isn't here yet.
// BUG(telyn): Flesh out the list of commands
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

Commands available:

	create-vm [flags] <name> [image] [disc-specs]
	create [flags] <name> [image] [disc-specs]
	new-vm [flags] <name> [image] [disc-specs]
	new [flags] <name> [image] [disc-specs]

		Creates a VM with the provided name, image and disc-specs.

		If image and disc-specs are not provided, will interactively prompt for them, or default to "wheezy" and a 25GiB SATA SSD.

		Flags available:

			--boot-script - A filename of a boot-script to upload from the local machine. This script will be run the first time the virtual machine starts up.
			--hwprofile   - Sets the hardware profile. Defaults to virtio2013 probably.
			--lock        - Locks the hardware profile, preventing automatic upgrades. Not set by default





*/
package main
