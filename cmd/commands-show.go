package cmd

import (
	"fmt"
	"strings"
)

func (dispatch *Dispatcher) ShowVM(args []string) {
	name := ParseVirtualMachineName(args[0])

	vm, err := dispatch.BigV.GetVirtualMachine(name)

	if err != nil {
		panic(err)
	}

	totalDiscSize := 0

	for _, disc := range vm.Discs {
		totalDiscSize += disc.Size
	}

	totalDiscSize = totalDiscSize / 1024

	fmt.Printf("= VM %s, %d cores, %d GiB RAM, %d GiB on %d discs =\r\n", vm.Name, vm.Cores, vm.Memory, totalDiscSize, len(vm.Discs))
	fmt.Printf("Hostname:		    %s\r\n", vm.Hostname)
	fmt.Printf("Hardware profile:	    %s\r\n", vm.HardwareProfile)
	fmt.Printf("Host machine:	    %s\r\n", vm.Head)
	for _, disc := range vm.Discs {
		fmt.Printf("Disc %s: %d GiB, %s grade\r\n", disc.Label, disc.Size/1024, disc.StorageGrade)
	}

	for _, nic := range vm.NetworkInterfaces {
		fmt.Printf("Network interface %s: %s\r\n", nic.Label, strings.Join(nic.IPs, ", "))
		//do something with extra IPs
	}

	// now I gotta decide how to pretty-print a vm.

	// --------------------------  VM main: 1 core, 2 GiB RAM, 635.0 GiB on 2 discs  -
	// Hostname:         main.personal.telyn.uk0.bigv.io
	// Hardware profile: virtio2013
	// CD-ROM:           (none)
	// VM host:          head05
	// Disc vda:                                            35.0 GiB,    sata grade
	// Disc vdb:                                           600.0 GiB, archive grade
	// Net :                    213.138.111.200, 2001:41c8:51:6c8:fcff:ff:fe00:3fe0
	// Extra IPs via 213.138.111.200:                               213.138.112.169

}

func (dispatch *Dispatcher) ShowAccount(args []string) {
	name := ParseAccountName(args[0])

	acc, err := dispatch.BigV.GetAccount(name)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Account %d: %s", acc.Id, acc.Name)

}
