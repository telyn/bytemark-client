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
	fmt.Println()
	fmt.Println("Displays information about the given virtual machine, group, or account.")
	fmt.Println("If the --verbose flag is given to bigv show group or bigv show account, full details are given for each VM.")
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

	err = cmds.bigv.DeleteVirtualMachine(name)

	if err != nil {
		exit(err)
	}

	fmt.Println(FormatVirtualMachine(vm))

}
