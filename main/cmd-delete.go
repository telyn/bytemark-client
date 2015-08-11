package main

import (
	bigv "bigv.io/client/lib"
	"fmt"
	"os"
)

// HelpForDelete outputs usage information for the delete command
func (cmds *CommandSet) HelpForDelete() {
	fmt.Println("go-bigv delete")
	fmt.Println()
	fmt.Println("usage: bigv delete [--force] [--purge] <name>")
	fmt.Println("       bigv delete vm [--force] [---purge] <virtual machine>")
	fmt.Println("       bigv delete group <group>")
	fmt.Println("       bigv delete account <account>")
	fmt.Println("       bigv delete user <auser>")
	fmt.Println("       bigv undelete vm <virtual machine>")
	fmt.Println()
	fmt.Println("Deletes the given virtual machine, group, account or user. Only empty groups and accounts can be deleted.")
	fmt.Println("If the --purge flag is given and the target is a virtual machine, will permanently delete the VM. Billing will cease and you will be unable to recover the VM.")
	fmt.Println("If the --force flag is given, you will not be prompted to confirm deletion.")
	fmt.Println()
	fmt.Println("The undelete vm command may be used to restore a deleted (but not purged) vm to its state prior to deletion.")
	fmt.Println()
}

// DeleteVM implements the delete-vm command, which is used to delete and purge BigV VMs. See HelpForDelete for usage information.
func (cmds *CommandSet) DeleteVM(args []string) ExitCode {
	flags := MakeCommonFlagSet()

	purge := *flags.Bool("purge", false, "Whether or not to purge the VM. If yes, will delete all your data.")

	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name := cmds.bigv.ParseVirtualMachineName(args[0])
	cmds.EnsureAuth()

	vm, err := cmds.bigv.GetVirtualMachine(name)
	if err != nil {
		return processError(err)
	}
	if vm.Deleted && !purge {
		fmt.Printf("Virtual machine %s has already been deleted.\r\nIf you wish to permanently delete it, add --purge", vm.Hostname)
		return E_SUCCESS
	}

	if !cmds.config.Force() {
		fstr := fmt.Sprintf("Are you certain you wish to delete %s?", vm.Hostname)
		if purge {
			fstr = fmt.Sprintf("Are you certain you wish to permanently delete %s? You will not be able to un-delete it.", vm.Hostname)

		}
		if !PromptYesNo(fstr) {
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
	return E_SUCCESS
}

func (cmds *CommandSet) DeleteGroup(args []string) ExitCode {
	flags := MakeCommonFlagSet()

	recursive := flags.Bool("recursive", false, "")

	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	name := cmds.bigv.ParseGroupName(flags.Args()[0])

	err := cmds.EnsureAuth()
	if err != nil {
		return processError(err)
	}

	group, err := cmds.bigv.GetGroup(name)
	if err != nil {
		return processError(err)
	}

	if len(group.VirtualMachines) > 0 {
		if *recursive {

			fmt.Fprintln(os.Stderr, "WARNING: The following VMs will be permanently deleted, without any way to recover or un-delete them:")
			for _, vm := range group.VirtualMachines {
				fmt.Fprintf(os.Stderr, "\t%s\r\n", vm.Name)
			}
			fmt.Fprint(os.Stderr, "\r\n\r\n")
			if PromptYesNo("Are you sure you want to continue?") {
				vmn := bigv.VirtualMachineName{Group: name.Group, Account: name.Account}
				for _, vm := range group.VirtualMachines {
					vmn.VirtualMachine = vm.Name
					err := cmds.bigv.DeleteVirtualMachine(vmn, true)
					if err != nil {
						return processError(err)
					} else {
						fmt.Fprintf(os.Stderr, "Virtual machine %s purged successfully.\r\n", name)
					}

				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Group %s contains virtual machines, will not be deleted without --recursive\r\n", name.Group)
			return E_WONT_DELETE_NONEMPTY
		}
	}
	return processError(cmds.bigv.DeleteGroup(name))

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

	return E_SUCCESS
}
