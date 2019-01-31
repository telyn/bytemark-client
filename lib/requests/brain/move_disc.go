package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// MoveDisc moves the specified disc from its current server to a given server
func MoveDisc(client lib.Client, vm lib.VirtualMachineName, discLabelOrID string, newVMName lib.VirtualMachineName) (err error) {
	err = client.EnsureVirtualMachineName(&vm)
	if err != nil {
		return err
	}
	err = client.EnsureVirtualMachineName(&newVMName)
	if err != nil {
		return err
	}

	r, err := client.BuildRequest("PUT", lib.BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s", string(vm.Account), vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	newVM, err := client.GetVirtualMachine(newVMName)
	if err != nil {
		return err
	}

	disc := brain.Disc{
		VirtualMachineID: newVM.ID,
	}

	_, _, err = r.MarshalAndRun(disc, nil)
	return err
}
