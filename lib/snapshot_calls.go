package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// CreateSnapshot creates a snapshot of the given disc, returning the snapshot if it was successful.
func (c *bytemarkClient) CreateSnapshot(vm VirtualMachineName, discLabelOrID string) (snapshot brain.Snapshot, err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/snapshots", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &snapshot)
	return
}

