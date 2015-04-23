package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func (cmds *CommandSet) HelpForDelete() {
	// TODO(telyn): Replace instances of bigv with $0, however you get $0 in go?
	fmt.Println("bigv delete")
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

func (cmds *CommandSet) DeleteVM(args []string) {

	flags := flag.NewFlagSet("DeleteVM", flag.ExitOnError)

	var purge = *flags.Bool("purge", false, "Whether or not to purge the VM. If yes, will delete all your data.")
	var force = *flags.Bool("force", false, "Don't confirm deletion. Be careful when using with --purge!")

	flags.Parse(args)

	name := ParseVirtualMachineName(args[0])
	cmds.EnsureAuth()

	vm, err := cmds.bigv.GetVirtualMachine(name)
	if err != nil {
		exit(err)
	}
	if vm.Deleted && !purge {
		exit(nil, "Virtual machine %s has already been deleted.\r\nIf you wish to permanently delete it, add --purge", vm.Hostname)
	}

	if !force {
		buf := bufio.NewReader(os.Stdin)
		fstr := "Are you certain you wish to delete %s? (y/n)"
		if purge {
			fstr = "Are you certain you wish to PERMANENTLY delete %s? (y/n)"

		}
		fmt.Fprintf(os.Stderr, fstr, vm.Hostname)
		chr, err := buf.ReadByte()
		if err != nil {
			exit(err)
		} else if chr != 'y' {
			exit(nil, "Aborting.")

		}
	}

	err = cmds.bigv.DeleteVirtualMachine(name, purge)

	if err != nil {
		exit(err)
	}

	fmt.Println(FormatVirtualMachine(vm))
}

func (cmds *CommandSet) UndeleteVM(args []string) {

	name := ParseVirtualMachineName(args[0])
	cmds.EnsureAuth()

	vm, err := cmds.bigv.GetVirtualMachine(name)
	if err != nil {
		exit(err)
	}

	if !vm.Deleted {
		exit(nil, fmt.Sprintf("Virtual machine %s was already undeleted", vm.Hostname))
	}

	err = cmds.bigv.UndeleteVirtualMachine(name)

	if err != nil {
		exit(err)
	}
	fmt.Printf("Successfully restored virtual machine %s\r\n", vm.Hostname)

}
