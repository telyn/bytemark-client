package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// HelpForDelete outputs usage information for the delete command
func (cmds *CommandSet) HelpForDelete() {
	fmt.Println("go-bigv delete")
	fmt.Println()
	fmt.Println("usage: bigv delete [-p | --purge] <name>")
	fmt.Println("       bigv delete vm [-p | ---purge] <virtual machine>")
	fmt.Println("       bigv delete group <group>")
	fmt.Println("       bigv delete account <account>")
	fmt.Println("       bigv delete user <auser>")
	fmt.Println("       bigv undelete vm <virtual machine>")
	fmt.Println()
	fmt.Println("Deletes the given virtual machine, group, account or user.")
	fmt.Println("If the --purge flag is given and the target is a virtual machine, will permanently delete the VM. Billing will cease and you will be unable to recover the VM.")
	fmt.Println()
	fmt.Println("The undelete vm command may be used to restore a deleted (but not purged) vm to its state prior to deletion.")
	fmt.Println()
}

// DeleteVM implements the delete-vm command, which is used to delete and purge BigV VMs. See HelpForDelete for usage information.
func (cmds *CommandSet) DeleteVM(args []string) ExitCode {

	// TODO(telyn): rewrite flag code here.
	flags := flag.NewFlagSet("DeleteVM", flag.ExitOnError)

	var purge = *flags.Bool("purge", false, "Whether or not to purge the VM. If yes, will delete all your data.")
	var force = *flags.Bool("force", false, "Don't confirm deletion. Be careful when using with --purge!")

	flags.Parse(args)

	name := cmds.bigv.ParseVirtualMachineName(flags.Args()[0])
	cmds.EnsureAuth()

	vm, err := cmds.bigv.GetVirtualMachine(name)
	if err != nil {
		return processError(err)
	}
	if vm.Deleted && !purge {
		fmt.Printf("Virtual machine %s has already been deleted.\r\nIf you wish to permanently delete it, add --purge", vm.Hostname)
		return E_SUCCESS
	}

	if !force {
		// TODO(telyn): use PromptYN
		buf := bufio.NewReader(os.Stdin)
		fstr := "Are you certain you wish to delete %s? (y/n) "
		if purge {
			fstr = "Are you certain you wish to PERMANENTLY delete %s? (y/n) "

		}
		fmt.Fprintf(os.Stderr, fstr, vm.Hostname)
		chr, err := buf.ReadByte()
		fmt.Println()
		if err != nil {
			return processError(err)
		} else if chr != 'y' {

			return processError(&UserRequestedExit{})

		}
	}

	err = cmds.bigv.DeleteVirtualMachine(name, purge)

	if err != nil {
		return processError(err)
	}

	if purge {
		fmt.Printf("Virtual machine %s purged successfully.\r\n", name)
	} else {
		fmt.Printf("Virtual machine %s deleted successfully.\r\n", name)
	}
	return 0
}

// UndeleteVM implements the undelete-vm command, which is used to remove the deleted flag from BigV VMs, allowing them to be reactivated.
func (cmds *CommandSet) UndeleteVM(args []string) ExitCode {

	name := cmds.bigv.ParseVirtualMachineName(args[0])
	cmds.EnsureAuth()

	vm, err := cmds.bigv.GetVirtualMachine(name)
	if err != nil {
		return processError(err)
	}

	if !vm.Deleted {
		fmt.Printf("Virtual machine %s was already undeleted", vm.Hostname)
		return E_SUCCESS
	}

	err = cmds.bigv.UndeleteVirtualMachine(name)

	if err != nil {
		return processError(err)
	}
	fmt.Printf("Successfully restored virtual machine %s\r\n", vm.Hostname)

	return 0
}
