package lib

import (
	"bytes"
	"encoding/json"
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

func (c *bytemarkClient) DeleteSnapshot(vm VirtualMachineName, discLabelOrID string, snapshotLabelOrID string) (err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("DELETE", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/snapshots/%s?purge=true", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID, snapshotLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) GetSnapshots(vm VirtualMachineName, discLabelOrID string) (snapshots brain.Snapshots, err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/snapshots", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &snapshots)
	return
}

func (c *bytemarkClient) RestoreSnapshot(vm VirtualMachineName, discLabelOrID string, snapshotLabelOrID string) (snapshot brain.Snapshot, err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/snapshots", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}
	restore := map[string]bool{
		"restore": true,
	}
	restoreJSON, err := json.Marshal(restore)
	if err != nil {
		return
	}

	_, _, err = r.Run(bytes.NewBuffer(restoreJSON), nil)
	return
}
